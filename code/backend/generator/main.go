package main

import (
	"log"

	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/controllers"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cloud.google.com/go/profiler"
)

func main() {
	profilerCfg := profiler.Config{
		Service:        "pre-test-generator",
		ServiceVersion: "1.0.0",
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",

		// For OpenCensus users:
		// To see Profiler agent spans in APM backend,
		// set EnableOCTelemetry to true
		// EnableOCTelemetry: true,
	}

	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profilerCfg); err != nil {
		log.Fatalf("Failed init profiler：%v", err)
	}

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
