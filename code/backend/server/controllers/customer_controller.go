package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

type CustomerController struct {
	service services.CustomerService
}

func NewCustomerController(service services.CustomerService) *CustomerController {
	return &CustomerController{service}
}

func (c *CustomerController) GetAllCustomers(ctx echo.Context) error {
	customers, err := c.service.GetAllCustomers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerController) GetLimitedCustomers(ctx echo.Context) error {
	num, err := strconv.Atoi(ctx.Param("num"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	customers, err := c.service.GetLimitedCustomers(num)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerController) CreateCustomer(ctx echo.Context) error {
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if len(customer.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, "Password at least 8 character")
	}
	customer.ID = uuid.New()
	if err := c.service.CreateCustomer(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, customer)
}

func (c *CustomerController) CreateMultiCustomers(ctx echo.Context) error {
	var customers []*models.Customer
	if err := ctx.Bind(&customers); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	successCount, failCount, err := c.service.CreateMultiCustomers(customers)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	result := map[string]int{
		"successCount": successCount,
		"failCount":    failCount,
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *CustomerController) GetCustomerByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	customer, err := c.service.GetCustomerByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) UpdateCustomer(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	customer.ID = id
	if err := c.service.UpdateCustomer(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) UpdateCustomerPassword(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if len(customer.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, "Password at least 8 character")
	}
	customer.ID = id
	if err := c.service.UpdateCustomerPassword(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) DeleteCustomer(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if err := c.service.DeleteCustomer(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *CustomerController) ResetAllCustomerData(ctx echo.Context) error {
	if err := c.service.ResetAllCustomerData(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
