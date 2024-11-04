package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/services"

	"github.com/labstack/echo/v4"
)

// CustomerController defines the interface for customer-related operations
type CustomerController interface {
	GenerateAndSendCustomerData(c echo.Context) error
}

// customerController is the concrete implementation of CustomerController
type customerController struct {
	cfg             *config.Config
	customerService services.CustomerService
}

// NewCustomerController is the factory function that returns a CustomerController interface
func NewCustomerController(cfg *config.Config, customerService services.CustomerService) CustomerController {
	return &customerController{
		cfg:             cfg,
		customerService: customerService,
	}
}

// GenerateAndSendCustomerData handles generating customer data and sending it to the backend
func (cc *customerController) GenerateAndSendCustomerData(c echo.Context) error {
	// Parse the 'num' query parameter and validate it
	numStr := c.QueryParam("num")
	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		log.Printf("Error converting num parameter: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid number of records"})
	}
	log.Printf("Received request to generate %d customer records", num)

	var sameFailureCounter int
	var generateDuration time.Duration
	var sendDuration time.Duration

	// Loop to attempt data generation and sending until successful or retry limit reached
	for {
		generateStartTime := time.Now()
		log.Println("Starting generation of customer data")
		// Generate customer data using the service interface
		customers, err := cc.customerService.GenerateCustomerData(num)
		if err != nil {
			log.Printf("Error generating customer data: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate customer data"})
		}
		generateDuration += time.Since(generateStartTime)

		sendAPIStartTime := time.Now()
		log.Println("Starting API call to send customer data")
		// Send customer data to the backend server using the service interface
		ctx := context.Background()
		_, failedCount, err := cc.customerService.CreateMultiCustomersAPICall(ctx, customers)
		if err != nil {
			log.Printf("Error during API call to backend: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to send customer data to backend: %s", err)})
		}
		sendDuration += time.Since(sendAPIStartTime)

		log.Printf("API call complete with %d failures", failedCount)
		// Check if there is a decrease in failures or if failures are constant
		if failedCount > 0 {
			if failedCount == num {
				// Increment counter if failure count remains constant
				sameFailureCounter++
				log.Printf("Failure count remained constant at %d, retry %d", failedCount, sameFailureCounter)
				if sameFailureCounter >= 5 {
					// Stop retries after 5 consecutive constant failures
					log.Println("Persistent failures reached, stopping retries")
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error":  "Persistent failures, stopping retries",
						"failed": strconv.Itoa(failedCount),
					})
				}
			} else {
				// Reset counter if failure count decreases
				sameFailureCounter = 0
				log.Println("Failure count decreased, resetting sameFailureCounter")
			}
			// Set num to failedCount for the next iteration
			num = failedCount
		} else {
			// Break loop if there are no failures
			log.Println("No failures, breaking loop")
			break
		}
	}

	// Return success response with generation and sending durations
	return c.JSON(http.StatusOK, map[string]string{
		"status":          "Customer data generated and sent to backend server",
		"generation_time": fmt.Sprintf("%v", generateDuration),
		"send_time":       fmt.Sprintf("%v", sendDuration),
	})
}
