package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

// TransactionController defines the interface for transaction-related handlers
type TransactionController interface {
	GetTransactionsByCustomerID(ctx echo.Context) error
	GetDateRangeTransactionsByCustomerID(ctx echo.Context) error
	CreateMultiTransactions(ctx echo.Context) error
	// CreateTransaction(ctx echo.Context) error
	// UpdateTransaction(ctx echo.Context) error
	// DeleteTransaction(ctx echo.Context) error
}

// transactionController is the concrete implementation of TransactionController
type transactionController struct {
	transactionService services.TransactionService
}

// NewTransactionController initializes a new TransactionController.
func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

// GetTransactionsByCustomerID retrieves all transactions for a specified customer.
func (tc *transactionController) GetTransactionsByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Customer ID"})
	}
	transactions, err := tc.transactionService.GetTransactionsByCustomerID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, transactions)
}

// GetDateRangeTransactionsByCustomerID retrieves transactions within a date range for a specified customer.
func (tc *transactionController) GetDateRangeTransactionsByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Customer ID"})
	}

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	transactions, err := tc.transactionService.GetDateRangeTransactionsByCustomerID(customerID, from, to)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, transactions)
}

// CreateMultiTransactions creates multiple transactions from the provided DTOs.
func (tc *transactionController) CreateMultiTransactions(ctx echo.Context) error {
	var transactions []*models.TransactionDTO
	if err := ctx.Bind(&transactions); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if len(transactions) > 5000 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "transactions length over 5000"})
	}

	if err := tc.transactionService.CreateMultiTransactions(transactions); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result := map[string]string{
		"result": "success",
	}

	return ctx.JSON(http.StatusCreated, result)
}

// CreateTransaction creates a new transaction with a generated ID.
// func (tc *transactionController) CreateTransaction(ctx echo.Context) error {
// 	transaction := new(models.Transaction)
// 	if err := ctx.Bind(transaction); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
// 	}
// 	transaction.ID = uuid.New() // Generate a new UUID for the transaction
// 	if err := tc.transactionService.CreateTransaction(transaction); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}
// 	return ctx.JSON(http.StatusCreated, transaction)
// }

// UpdateTransaction updates an existing transaction specified by ID.
// func (tc *transactionController) UpdateTransaction(ctx echo.Context) error {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
// 	}
// 	transaction := new(models.Transaction)
// 	if err := ctx.Bind(transaction); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
// 	}
// 	transaction.ID = id // Set the transaction ID for updating
// 	if err := tc.transactionService.UpdateTransaction(transaction); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}
// 	return ctx.JSON(http.StatusOK, transaction)
// }

// DeleteTransaction deletes a transaction specified by ID.
// func (tc *transactionController) DeleteTransaction(ctx echo.Context) error {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
// 	}
// 	if err := tc.transactionService.DeleteTransaction(id); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}
// 	return ctx.NoContent(http.StatusNoContent) // Return 204 No Content on successful deletion
// }
