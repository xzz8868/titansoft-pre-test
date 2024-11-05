package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

// CustomerController defines the interface for customer-related operations
type CustomerController interface {
	GetAllCustomers(ctx echo.Context) error
	GetLimitedCustomers(ctx echo.Context) error
	CreateCustomer(ctx echo.Context) error
	CreateMultiCustomers(ctx echo.Context) error
	GetCustomerByID(ctx echo.Context) error
	UpdateCustomer(ctx echo.Context) error
	UpdateCustomerPassword(ctx echo.Context) error
	DeleteCustomer(ctx echo.Context) error
	ResetAllCustomerData(ctx echo.Context) error
}

// CustomerController handles HTTP requests related to customers
type customerController struct {
	customerService services.CustomerService
}

// NewCustomerController initializes a new CustomerController
func NewCustomerController(customerService services.CustomerService) CustomerController {
	return &customerController{
		customerService: customerService,
	}
}

// GetAllCustomers retrieves all customers
func (cc *customerController) GetAllCustomers(ctx echo.Context) error {
	customers, err := cc.customerService.GetAllCustomers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, customers)
}

// GetLimitedCustomers retrieves a limited number of customers based on 'num' parameter
func (cc *customerController) GetLimitedCustomers(ctx echo.Context) error {
	num, err := strconv.Atoi(ctx.Param("num"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	customers, err := cc.customerService.GetLimitedCustomers(num)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, customers)
}

// CreateCustomer adds a new customer to the database
func (cc *customerController) CreateCustomer(ctx echo.Context) error {
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if len(customer.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters"})
	}
	customer.ID = uuid.New()
	if err := cc.customerService.CreateCustomer(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, customer)
}

// CreateMultiCustomers adds multiple customers at once
func (cc *customerController) CreateMultiCustomers(ctx echo.Context) error {
	var customers []*models.Customer
	if err := ctx.Bind(&customers); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if len(customers) > 2000 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Customers length over 2000"})
	}

	successCount, failCount, err := cc.customerService.CreateMultiCustomers(customers)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result := map[string]int{
		"successCount": successCount,
		"failCount":    failCount,
	}
	return ctx.JSON(http.StatusCreated, result)
}

// GetCustomerByID retrieves a customer by their unique ID
func (cc *customerController) GetCustomerByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	customer, err := cc.customerService.GetCustomerByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, customer)
}

// UpdateCustomer updates customer details by ID
func (cc *customerController) UpdateCustomer(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	customer.ID = id
	if err := cc.customerService.UpdateCustomer(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, customer)
}

// UpdateCustomerPassword updates only the password of a customer by ID
func (cc *customerController) UpdateCustomerPassword(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if len(customer.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters"})
	}
	customer.ID = id
	if err := cc.customerService.UpdateCustomerPassword(customer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, customer)
}

// DeleteCustomer removes a customer by their unique ID
func (cc *customerController) DeleteCustomer(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	if err := cc.customerService.DeleteCustomer(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}

// ResetAllCustomerData resets all customer data in the system
func (cc *customerController) ResetAllCustomerData(ctx echo.Context) error {
	if err := cc.customerService.ResetAllCustomerData(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "All customer data reset successfully"})
}
