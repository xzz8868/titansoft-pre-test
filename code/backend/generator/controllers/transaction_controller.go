package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/services"
)

// TransactionController defines the interface for transaction-related handlers
type TransactionController interface {
	CreateTransactions(ctx echo.Context) error
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
func (tc *transactionController) CreateTransactions(ctx echo.Context) error {
	// Parse "transactions_num" query parameter to integer
	numTransactionsStr := ctx.QueryParam("transactions_num")
	numTransactions, err := strconv.Atoi(numTransactionsStr)
	if err != nil || numTransactions <= 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of transactions"})
	}

	if numTransactions > 5000 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Transactions num over 10000"})
	}

	// Parse "customers_num" query parameter to integer
	numCustomersStr := ctx.QueryParam("customers_num")
	numCustomers, err := strconv.Atoi(numCustomersStr)
	if err != nil || numCustomers <= 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of customers"})
	}

	// Generate and send transactions using the service layer
	if err := tc.transactionService.GenerateAndSendTransactions(numTransactions, numCustomers); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return success response
	return ctx.JSON(http.StatusOK, map[string]string{"status": "Transactions generated and sent successfully"})
}
