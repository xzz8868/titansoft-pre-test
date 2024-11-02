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
	numTransactionsStr := c.QueryParam("transactions_num")
	numCustomersStr := c.QueryParam("customers_num")

	numTransactions, err := strconv.Atoi(numTransactionsStr)
	if err != nil || numTransactions <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of transactions"})
	}

	numCustomers, err := strconv.Atoi(numCustomersStr)
	if err != nil || numCustomers <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of customers"})
	}

	ctx := c.Request().Context()
	if err := tc.transactionService.GenerateAndSendTransactions(ctx, numTransactions, numCustomers); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Transactions generated and sent successfully"})
}