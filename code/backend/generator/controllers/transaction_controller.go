package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/services"
)

// TransactionController defines the interface for transaction-related handlers
type TransactionController interface {
	CreateTransactions(c echo.Context) error
}

// transactionController is the concrete implementation of TransactionController
type transactionController struct {
	transactionService services.TransactionService
}

// NewTransactionController is the factory function that returns a TransactionController interface
func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

// CreateTransactions handles the creation of transactions
func (tc *transactionController) CreateTransactions(c echo.Context) error {
	// Parse "transactions_num" query parameter to integer
	numTransactionsStr := c.QueryParam("transactions_num")
	numTransactions, err := strconv.Atoi(numTransactionsStr)
	if err != nil || numTransactions <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of transactions"})
	}

	if numTransactions > 10000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Transactions num over 10000"})
	}

	// Parse "customers_num" query parameter to integer
	numCustomersStr := c.QueryParam("customers_num")
	numCustomers, err := strconv.Atoi(numCustomersStr)
	if err != nil || numCustomers <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of customers"})
	}

	// Generate and send transactions using the service layer
	ctx := c.Request().Context()
	if err := tc.transactionService.GenerateAndSendTransactions(ctx, numTransactions, numCustomers); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{"status": "Transactions generated and sent successfully"})
}
