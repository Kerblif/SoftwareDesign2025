package main

import (
	"bytes"
	orderModels "cbd/internal/orders/models"
	paymentsModels "cbd/internal/payments/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strings"
)

type Status struct {
	Status string `json:"status"`
}

// postOrder would post an order to the API
func postOrder(accountID int64, amount float64, description string) (*orderModels.Order, error) {
	var newOrder = orderModels.CreateOrderRequest{
		UserID:      accountID,
		Amount:      amount,
		Description: description,
	}

	body, err := json.Marshal(newOrder)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/orders", apiBaseURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order orderModels.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// fetchOrders would fetch orders from the API
func fetchOrders(accountID uint) ([]orderModels.Order, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/orders?user_id=%d", apiBaseURL, accountID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orders struct {
		Orders []orderModels.Order `json:"orders"`
	}
	err = json.Unmarshal(body, &orders)
	if err != nil {
		return nil, err
	}

	return orders.Orders, nil
}

// getOrder would fetch an order from the API
func getOrder(orderID uint) (*orderModels.Order, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/orders/%d", apiBaseURL, orderID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order orderModels.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// createNewAccount would create a new account for the user
func createNewAccount(accountID int64) error {
	var newAccount = paymentsModels.CreateAccountRequest{
		UserID: accountID,
	}

	body, err := json.Marshal(newAccount)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/payments/accounts", apiBaseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to create account: %s", resp.Status)
	}

	return nil
}

// addDeposit would add a deposit to the account
func addDeposit(accountID int64, amount float64) error {
	var deposit = paymentsModels.DepositRequest{
		UserID: accountID,
		Amount: amount,
	}

	body, err := json.Marshal(deposit)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/payments/deposit", apiBaseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to add deposit: %s", resp.Status)
	}

	return nil
}

// getBalance would get the balance of the account
func getBalance(accountID int64) (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/payments/balance?user_id=%d", apiBaseURL, accountID), nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var balance paymentsModels.BalanceResponse
	err = json.Unmarshal(body, &balance)

	return balance.Balance, err
}

func getPaymentStatus() (*Status, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/health/payments", apiBaseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status Status
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// subscribeToOrderUpdates would subscribe to order updates via websocket
func subscribeToOrderUpdates(orderID int64) (chan *orderModels.PaymentMessage, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws/orders/%d", strings.TrimPrefix(apiBaseURL, "http://"), orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	msgChan := make(chan *orderModels.PaymentMessage)

	go func() {
		defer conn.Close()
		defer close(msgChan)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading from websocket: %v", err)
				return
			}

			var msg orderModels.PaymentMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			msgChan <- &msg
		}
	}()

	return msgChan, nil
}
