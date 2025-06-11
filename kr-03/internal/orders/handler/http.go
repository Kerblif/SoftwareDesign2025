package handler

import (
	"net/http"
	"strconv"

	"cbd/internal/orders/models"
	"cbd/internal/orders/service"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the orders service
type Handler struct {
	orderSvc service.OrderService
}

// NewHandler creates a new HTTP handler
func NewHandler(orderSvc service.OrderService) *Handler {
	return &Handler{
		orderSvc: orderSvc,
	}
}

// RegisterRoutes registers the HTTP routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/orders", h.CreateOrder)
	router.GET("/orders", h.GetOrders)
	router.GET("/orders/:id", h.GetOrder)
	router.GET("/health", h.HealthCheck)
}

// CreateOrder handles the request to create a new order
func (h *Handler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderSvc.CreateOrder(c.Request.Context(), req.UserID, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response model
	response := &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Amount:      order.Amount,
		Description: order.Description,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrders handles the request to get orders for a user
func (h *Handler) GetOrders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	orders, err := h.orderSvc.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response models
	var responses []*models.OrderResponse
	for _, order := range orders {
		responses = append(responses, &models.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			Amount:      order.Amount,
			Description: order.Description,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetOrder handles the request to get a specific order
func (h *Handler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderSvc.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response model
	response := &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Amount:      order.Amount,
		Description: order.Description,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles the health check request
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
