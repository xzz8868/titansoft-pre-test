package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

// TransactionController handles transaction-related requests.
type TransactionController struct {
	service services.TransactionService
}

// NewTransactionController initializes a new TransactionController.
func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{service}
}

// GetTransactionsByCustomerID retrieves all transactions for a specified customer.
func (t *TransactionController) GetTransactionsByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Customer ID"})
	}
	transactions, err := t.service.GetTransactionsByCustomerID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, transactions)
}

// GetDateRangeTransactionsByCustomerID retrieves transactions within a date range for a specified customer.
func (t *TransactionController) GetDateRangeTransactionsByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Customer ID"})
	}

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	transactions, err := t.service.GetDateRangeTransactionsByCustomerID(customerID, from, to)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, transactions)
}

// CreateTransaction creates a new transaction with a generated ID.
func (t *TransactionController) CreateTransaction(ctx echo.Context) error {
	transaction := new(models.Transaction)
	if err := ctx.Bind(transaction); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	transaction.ID = uuid.New() // Generate a new UUID for the transaction
	if err := t.service.CreateTransaction(transaction); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, transaction)
}

// CreateMultiTransactions creates multiple transactions from the provided DTOs.
func (t *TransactionController) CreateMultiTransactions(ctx echo.Context) error {
	var transactions []*models.TransactionDTO
	if err := ctx.Bind(&transactions); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if len(transactions) > 10000 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "transactions length over 10000"})
	}

	if err := t.service.CreateMultiTransactions(transactions); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result := map[string]string{
		"result": "success",
	}

	return ctx.JSON(http.StatusCreated, result)
}

// UpdateTransaction updates an existing transaction specified by ID.
func (t *TransactionController) UpdateTransaction(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	transaction := new(models.Transaction)
	if err := ctx.Bind(transaction); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	transaction.ID = id // Set the transaction ID for updating
	if err := t.service.UpdateTransaction(transaction); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a transaction specified by ID.
func (t *TransactionController) DeleteTransaction(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	if err := t.service.DeleteTransaction(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent) // Return 204 No Content on successful deletion
}
