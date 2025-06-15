package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"cbd/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestNewWebSocketHandler tests the NewWebSocketHandler function
func TestNewWebSocketHandler(t *testing.T) {
	handler := NewWebSocketHandler()

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.connections)
	assert.NotNil(t, handler.upgrader)
}

// TestRegisterRoutes tests the RegisterRoutes function
func TestRegisterWebSocketRoutes(t *testing.T) {
	handler := NewWebSocketHandler()
	router := gin.New()
	handler.RegisterRoutes(router)

	// This is a simple test to ensure routes are registered
	assert.NotNil(t, router)
}

// TestHandleWebSocket tests the HandleWebSocket function
func TestHandleWebSocket(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewWebSocketHandler()
	handler.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Replace http with ws in the URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/orders/123"

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Verify connection is stored in the handler
	time.Sleep(100 * time.Millisecond) // Give some time for the connection to be stored

	handler.connectionsMutex.RLock()
	conns := handler.connections[123]
	handler.connectionsMutex.RUnlock()

	assert.Equal(t, 1, len(conns))
}

// TestHandleWebSocketInvalidID tests the HandleWebSocket function with an invalid order ID
func TestHandleWebSocketInvalidID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewWebSocketHandler()
	handler.RegisterRoutes(router)

	// Create a request with an invalid order ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ws/orders/invalid", nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid order id", response.Error)
}

// TestRemoveConnection tests the removeConnection function
func TestRemoveConnection(t *testing.T) {
	// Setup
	handler := NewWebSocketHandler()

	// Create a mock connection
	conn := &websocket.Conn{}

	// Add the connection to the map
	orderID := int64(123)
	handler.connectionsMutex.Lock()
	handler.connections[orderID] = append(handler.connections[orderID], conn)
	handler.connectionsMutex.Unlock()

	// Verify connection is stored
	handler.connectionsMutex.RLock()
	conns := handler.connections[orderID]
	handler.connectionsMutex.RUnlock()
	assert.Equal(t, 1, len(conns))

	// Remove the connection
	handler.removeConnection(orderID, conn)

	// Verify connection is removed
	handler.connectionsMutex.RLock()
	conns = handler.connections[orderID]
	handler.connectionsMutex.RUnlock()
	assert.Equal(t, 0, len(conns))

	// Verify the order ID is removed from the map
	handler.connectionsMutex.RLock()
	_, exists := handler.connections[orderID]
	handler.connectionsMutex.RUnlock()
	assert.False(t, exists)
}

// TestNotifyOrderStatusChange tests the NotifyOrderStatusChange function
func TestNotifyOrderStatusChange(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewWebSocketHandler()
	handler.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Replace http with ws in the URL
	orderID := int64(123)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/orders/" + strconv.FormatInt(orderID, 10)

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Give some time for the connection to be stored
	time.Sleep(100 * time.Millisecond)

	// Create an order
	order := &models.Order{
		ID:          orderID,
		UserID:      1,
		Amount:      100.0,
		Description: "Test order",
		Status:      models.OrderStatusFinished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Notify about order status change
	handler.NotifyOrderStatusChange(context.Background(), order)

	// Read the notification from the WebSocket
	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Verify the notification
	var receivedOrder models.Order
	err = json.Unmarshal(message, &receivedOrder)
	assert.NoError(t, err)
	assert.Equal(t, order.ID, receivedOrder.ID)
	assert.Equal(t, order.Status, receivedOrder.Status)
}

// TestNotifyOrderStatusChangeNoConnections tests the NotifyOrderStatusChange function with no connections
func TestNotifyOrderStatusChangeNoConnections(t *testing.T) {
	// Setup
	handler := NewWebSocketHandler()

	// Create an order
	order := &models.Order{
		ID:          123,
		UserID:      1,
		Amount:      100.0,
		Description: "Test order",
		Status:      models.OrderStatusFinished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Notify about order status change (should not panic)
	handler.NotifyOrderStatusChange(context.Background(), order)

	// No assertions needed, just making sure it doesn't panic
}
