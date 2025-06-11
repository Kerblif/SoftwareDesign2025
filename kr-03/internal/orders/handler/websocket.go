package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"cbd/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time order status updates
type WebSocketHandler struct {
	// Map of order ID to WebSocket connections
	connections      map[int64][]*websocket.Conn
	connectionsMutex sync.RWMutex
	upgrader         websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		connections: make(map[int64][]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins
			},
		},
	}
}

// RegisterRoutes registers the WebSocket routes
func (h *WebSocketHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws/orders/:id", h.HandleWebSocket)
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get order ID from URL
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Add connection to the map
	h.connectionsMutex.Lock()
	h.connections[orderID] = append(h.connections[orderID], conn)
	h.connectionsMutex.Unlock()

	// Remove connection when it's closed
	defer func() {
		conn.Close()
		h.removeConnection(orderID, conn)
	}()

	// Keep the connection alive
	for {
		// Read message (just to detect disconnection)
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// removeConnection removes a WebSocket connection from the map
func (h *WebSocketHandler) removeConnection(orderID int64, conn *websocket.Conn) {
	h.connectionsMutex.Lock()
	defer h.connectionsMutex.Unlock()

	conns := h.connections[orderID]
	for i, c := range conns {
		if c == conn {
			// Remove the connection from the slice
			h.connections[orderID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// If there are no more connections for this order, remove the entry
	if len(h.connections[orderID]) == 0 {
		delete(h.connections, orderID)
	}
}

// NotifyOrderStatusChange notifies all connected clients about an order status change
func (h *WebSocketHandler) NotifyOrderStatusChange(ctx context.Context, order *models.Order) {
	h.connectionsMutex.RLock()
	conns := h.connections[order.ID]
	h.connectionsMutex.RUnlock()

	if len(conns) == 0 {
		return
	}

	// Create notification message
	notification := map[string]interface{}{
		"order_id": order.ID,
		"status":   order.Status,
	}

	// Marshal notification to JSON
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	// Send notification to all connected clients
	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, notificationBytes)
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
			// Connection might be closed, but we'll let the read loop handle it
		}
	}
}