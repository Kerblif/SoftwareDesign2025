package handler

import (
	"net/http"
	"strconv"

	"cbd/internal/payments/models"
	"cbd/internal/payments/service"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the payments service
type Handler struct {
	paymentSvc service.PaymentService
}

// NewHandler creates a new HTTP handler
func NewHandler(paymentSvc service.PaymentService) *Handler {
	return &Handler{
		paymentSvc: paymentSvc,
	}
}

// RegisterRoutes registers the HTTP routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/account", h.CreateAccount)
	router.POST("/deposit", h.Deposit)
	router.GET("/balance", h.GetBalance)
	router.GET("/health", h.HealthCheck)
}

// CreateAccount handles the request to create a new account
func (h *Handler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.paymentSvc.CreateAccount(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// Deposit handles the request to deposit money to an account
func (h *Handler) Deposit(c *gin.Context) {
	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	err := h.paymentSvc.Deposit(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetBalance handles the request to get the balance of an account
func (h *Handler) GetBalance(c *gin.Context) {
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

	balance, err := h.paymentSvc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// HealthCheck handles the health check request
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
