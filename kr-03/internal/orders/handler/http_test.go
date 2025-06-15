package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"cbd/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockOrderService is a mock implementation of the service.OrderService interface
type MockOrderService struct {
	CreateOrderFunc          func(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error)
	GetOrderByIDFunc         func(ctx context.Context, id int64) (*models.Order, error)
	GetOrdersByUserIDFunc    func(ctx context.Context, userID int64) ([]*models.Order, error)
	ProcessPaymentResultFunc func(ctx context.Context, result *models.PaymentResult) error
}

// Implement OrderService interface
func (m *MockOrderService) CreateOrder(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error) {
	return m.CreateOrderFunc(ctx, userID, amount, description)
}

func (m *MockOrderService) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	return m.GetOrderByIDFunc(ctx, id)
}

func (m *MockOrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	return m.GetOrdersByUserIDFunc(ctx, userID)
}

func (m *MockOrderService) ProcessPaymentResult(ctx context.Context, result *models.PaymentResult) error {
	return m.ProcessPaymentResultFunc(ctx, result)
}

// setupRouter creates a new gin router and handler for testing
func setupRouter(mockService *MockOrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(mockService)
	handler.RegisterRoutes(router)
	return router
}

// TestNewHandler tests the NewHandler function
func TestNewHandler(t *testing.T) {
	mockService := &MockOrderService{}
	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.orderSvc)
}

// TestRegisterRoutes tests the RegisterRoutes function
func TestRegisterRoutes(t *testing.T) {
	mockService := &MockOrderService{}
	router := gin.New()
	handler := NewHandler(mockService)
	handler.RegisterRoutes(router)

	// This is a simple test to ensure routes are registered
	// A more comprehensive test would check if the routes actually work
	assert.NotNil(t, router)
}

// TestCreateOrder tests the CreateOrder handler
func TestCreateOrder(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error)
		expectedStatus int
		expectedBody   bool // true if we expect a valid response body
	}{
		{
			name: "Success",
			requestBody: models.CreateOrderRequest{
				UserID:      1,
				Amount:      100.0,
				Description: "Test order",
			},
			mockFunc: func(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error) {
				return &models.Order{
					ID:          123,
					UserID:      userID,
					Amount:      amount,
					Description: description,
					Status:      models.OrderStatusNew,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   true,
		},
		{
			name: "Invalid request body",
			requestBody: struct {
				InvalidField string `json:"invalid_field"`
			}{
				InvalidField: "invalid",
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name: "Service error",
			requestBody: models.CreateOrderRequest{
				UserID:      1,
				Amount:      100.0,
				Description: "Test order",
			},
			mockFunc: func(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockOrderService{
				CreateOrderFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedBody {
				var response models.OrderResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
			}
		})
	}
}

// TestGetOrders tests the GetOrders handler
func TestGetOrders(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, userID int64) ([]*models.Order, error)
		expectedStatus int
		expectedBody   bool // true if we expect a valid response body
	}{
		{
			name:   "Success",
			userID: "1",
			mockFunc: func(ctx context.Context, userID int64) ([]*models.Order, error) {
				return []*models.Order{
					{
						ID:          123,
						UserID:      userID,
						Amount:      100.0,
						Description: "Test order 1",
						Status:      models.OrderStatusNew,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          124,
						UserID:      userID,
						Amount:      200.0,
						Description: "Test order 2",
						Status:      models.OrderStatusFinished,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "Missing user_id",
			userID:         "",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:           "Invalid user_id",
			userID:         "invalid",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:   "Service error",
			userID: "1",
			mockFunc: func(ctx context.Context, userID int64) ([]*models.Order, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockOrderService{
				GetOrdersByUserIDFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			url := "/orders"
			if tc.userID != "" {
				url += "?user_id=" + tc.userID
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedBody {
				var response struct {
					Orders []*models.OrderResponse `json:"orders"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Orders)
				assert.Equal(t, 2, len(response.Orders))
			}
		})
	}
}

// TestGetOrder tests the GetOrder handler
func TestGetOrder(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		orderID        string
		mockFunc       func(ctx context.Context, id int64) (*models.Order, error)
		expectedStatus int
		expectedBody   bool // true if we expect a valid response body
	}{
		{
			name:    "Success",
			orderID: "123",
			mockFunc: func(ctx context.Context, id int64) (*models.Order, error) {
				return &models.Order{
					ID:          id,
					UserID:      1,
					Amount:      100.0,
					Description: "Test order",
					Status:      models.OrderStatusNew,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "Invalid order_id",
			orderID:        "invalid",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:    "Service error",
			orderID: "123",
			mockFunc: func(ctx context.Context, id int64) (*models.Order, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockOrderService{
				GetOrderByIDFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/orders/"+tc.orderID, nil)
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedBody {
				var response models.OrderResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)

				id, _ := strconv.ParseInt(tc.orderID, 10, 64)
				assert.Equal(t, id, response.ID)
			}
		})
	}
}

// TestHealthCheck tests the HealthCheck handler
func TestHealthCheck(t *testing.T) {
	// Setup mock service
	mockService := &MockOrderService{}
	router := setupRouter(mockService)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check response
	assert.Equal(t, http.StatusOK, resp.Code)

	var response struct {
		Status string `json:"status"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
}
