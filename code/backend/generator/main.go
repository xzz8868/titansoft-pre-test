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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()

	// 添加CORS中間件
	e.Use(middleware.CORS())

	// Instantiate the service via its factory function
	customerService := services.NewCustomerService(cfg)
	// Instantiate the transaction service via its factory function
	transactionService := services.NewTransactionService(cfg)

	// Instantiate the transaction controller via its factory function
	transactionController := controllers.NewTransactionController(transactionService)
	// Instantiate the controller via its factory function, injecting the service interface
	customerController := controllers.NewCustomerController(cfg, customerService)

	// Define the route for creating transactions
	e.POST("/generate/customer", customerController.GenerateAndSendCustomerData)
	e.POST("/generate/transactions", transactionController.CreateTransactions)

	e.Logger.Fatal(e.Start(":" + cfg.GeneratorServerPort))
}
