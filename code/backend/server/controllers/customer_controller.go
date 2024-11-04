package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

// CustomerController handles HTTP requests related to customers
type CustomerController struct {
	service services.CustomerService
}

// NewCustomerController initializes a new CustomerController
func NewCustomerController(service services.CustomerService) *CustomerController {
	return &CustomerController{service}
}

// GetAllCustomers retrieves all customers
func (c *CustomerController) GetAllCustomers(ctx echo.Context) error {
	customers, err := c.service.GetAllCustomers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}

// GetLimitedCustomers retrieves a limited number of customers based on 'num' parameter
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

// CreateCustomer adds a new customer to the database
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

// CreateMultiCustomers adds multiple customers at once
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

// GetCustomerByID retrieves a customer by their unique ID
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

// UpdateCustomer updates customer details by ID
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

// UpdateCustomerPassword updates only the password of a customer by ID
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

// DeleteCustomer removes a customer by their unique ID
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

// ResetAllCustomerData resets all customer data in the system
func (c *CustomerController) ResetAllCustomerData(ctx echo.Context) error {
	if err := c.service.ResetAllCustomerData(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
