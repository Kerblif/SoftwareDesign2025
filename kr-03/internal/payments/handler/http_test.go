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

	"cbd/internal/payments/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockPaymentService is a mock implementation of the service.PaymentService interface
type MockPaymentService struct {
	CreateAccountFunc       func(ctx context.Context, userID int64) (*models.Account, error)
	DepositFunc             func(ctx context.Context, userID int64, amount float64) error
	GetBalanceFunc          func(ctx context.Context, userID int64) (*models.BalanceResponse, error)
	ProcessPaymentFunc      func(ctx context.Context, payment *models.PaymentMessage) error
	ProcessInboxMessageFunc func(ctx context.Context, message *models.InboxMessage) error
}

// Implement PaymentService interface
func (m *MockPaymentService) CreateAccount(ctx context.Context, userID int64) (*models.Account, error) {
	return m.CreateAccountFunc(ctx, userID)
}

func (m *MockPaymentService) Deposit(ctx context.Context, userID int64, amount float64) error {
	return m.DepositFunc(ctx, userID, amount)
}

func (m *MockPaymentService) GetBalance(ctx context.Context, userID int64) (*models.BalanceResponse, error) {
	return m.GetBalanceFunc(ctx, userID)
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, payment *models.PaymentMessage) error {
	return m.ProcessPaymentFunc(ctx, payment)
}

func (m *MockPaymentService) ProcessInboxMessage(ctx context.Context, message *models.InboxMessage) error {
	return m.ProcessInboxMessageFunc(ctx, message)
}

// setupRouter creates a new gin router and handler for testing
func setupRouter(mockService *MockPaymentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(mockService)
	handler.RegisterRoutes(router)
	return router
}

// TestNewHandler tests the NewHandler function
func TestNewHandler(t *testing.T) {
	mockService := &MockPaymentService{}
	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.paymentSvc)
}

// TestRegisterRoutes tests the RegisterRoutes function
func TestRegisterRoutes(t *testing.T) {
	mockService := &MockPaymentService{}
	router := gin.New()
	handler := NewHandler(mockService)
	handler.RegisterRoutes(router)

	// This is a simple test to ensure routes are registered
	// A more comprehensive test would check if the routes actually work
	assert.NotNil(t, router)
}

// TestCreateAccount tests the CreateAccount handler
func TestCreateAccount(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, userID int64) (*models.Account, error)
		expectedStatus int
		expectedBody   bool // true if we expect a valid response body
	}{
		{
			name: "Success",
			requestBody: models.CreateAccountRequest{
				UserID: 1,
			},
			mockFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
				return &models.Account{
					ID:        123,
					UserID:    userID,
					Balance:   0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
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
			requestBody: models.CreateAccountRequest{
				UserID: 1,
			},
			mockFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockPaymentService{
				CreateAccountFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/account", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedBody {
				var response models.Account
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
			}
		})
	}
}

// TestDeposit tests the Deposit handler
func TestDeposit(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, userID int64, amount float64) error
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: models.DepositRequest{
				UserID: 1,
				Amount: 100.0,
			},
			mockFunc: func(ctx context.Context, userID int64, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusOK,
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
		},
		{
			name: "Negative amount",
			requestBody: models.DepositRequest{
				UserID: 1,
				Amount: -100.0,
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Zero amount",
			requestBody: models.DepositRequest{
				UserID: 1,
				Amount: 0,
			},
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service error",
			requestBody: models.DepositRequest{
				UserID: 1,
				Amount: 100.0,
			},
			mockFunc: func(ctx context.Context, userID int64, amount float64) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockPaymentService{
				DepositFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/deposit", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedStatus == http.StatusOK {
				var response struct {
					Status string `json:"status"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response.Status)
			}
		})
	}
}

// TestGetBalance tests the GetBalance handler
func TestGetBalance(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, userID int64) (*models.BalanceResponse, error)
		expectedStatus int
		expectedBody   bool // true if we expect a valid response body
	}{
		{
			name:   "Success",
			userID: "1",
			mockFunc: func(ctx context.Context, userID int64) (*models.BalanceResponse, error) {
				return &models.BalanceResponse{
					UserID:  userID,
					Balance: 100.0,
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
			mockFunc: func(ctx context.Context, userID int64) (*models.BalanceResponse, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockPaymentService{
				GetBalanceFunc: tc.mockFunc,
			}
			router := setupRouter(mockService)

			// Create request
			url := "/balance"
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
				var response models.BalanceResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)

				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				assert.Equal(t, userID, response.UserID)
				assert.Equal(t, 100.0, response.Balance)
			}
		})
	}
}

// TestHealthCheck tests the HealthCheck handler
func TestHealthCheck(t *testing.T) {
	// Setup mock service
	mockService := &MockPaymentService{}
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
