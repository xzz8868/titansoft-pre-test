package main

import (
	"log"

	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/controllers"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Echo instance for handling HTTP requests
	e := echo.New()

	// Add CORS middleware for cross-origin requests
	e.Use(middleware.CORS())

	// Instantiate services
	customerService := services.NewCustomerService(cfg)
	transactionService := services.NewTransactionService(cfg)

	// Instantiate controllers with dependencies injected
	transactionController := controllers.NewTransactionController(transactionService)
	customerController := controllers.NewCustomerController(customerService)

	// Define routes for API endpoints
	e.POST("/generate/customer", customerController.GenerateAndSendCustomerData)
	e.POST("/generate/transactions", transactionController.CreateTransactions)

	// Start the server on the configured port
	e.Logger.Fatal(e.Start(":" + cfg.GeneratorServerPort))
}
