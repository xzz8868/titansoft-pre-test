package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{service}
}

func (t *TransactionController) GetAllTransactionByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid Customer ID")
	}
	transactions, err := t.service.GetAllTransactionByCustomerID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) GetDateRangeTransactionsByCustomerID(ctx echo.Context) error {
	customerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid Customer ID")
	}

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	transactions, err := t.service.GetDateRangeTransactionsByCustomerID(customerID, from, to)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) CreateTransaction(ctx echo.Context) error {
	transaction := new(models.Transaction)
	if err := ctx.Bind(transaction); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	transaction.ID = uuid.New()
	if err := t.service.CreateTransaction(transaction); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, transaction)
}

func (t *TransactionController) CreateMultiTransactions(ctx echo.Context) error {
	var transactions []*models.TransactionDTO
	if err := ctx.Bind(&transactions); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := t.service.CreateMultiTransactions(transactions); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	result := map[string]string{
		"result": "success",
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (t *TransactionController) UpdateTransaction(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	transaction := new(models.Transaction)
	if err := ctx.Bind(transaction); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	transaction.ID = id
	if err := t.service.UpdateTransaction(transaction); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, transaction)
}

func (t *TransactionController) DeleteTransaction(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if err := t.service.DeleteTransaction(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}
