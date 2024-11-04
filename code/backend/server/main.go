package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/controllers"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/repositories"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	dsn := cfg.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto-migrate database models
	if err := db.AutoMigrate(&models.Customer{}, &models.Transaction{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Initialize repositories
	customerRepo := repositories.NewCustomerRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	customerService := services.NewCustomerService(customerRepo, transactionRepo, cfg.Salt)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize controllers
	customerController := controllers.NewCustomerController(customerService)
	transactionController := controllers.NewTransactionController(transactionService)

	// Initialize Echo instance
	e := echo.New()

	// Add CORS middleware
	e.Use(middleware.CORS())

	// Set up routes
	// Routes for FrontEnd
	e.GET("/customers", customerController.GetAllCustomers)
	e.GET("/customers/:id", customerController.GetCustomerByID)
	e.POST("/customers", customerController.CreateCustomer)
	e.PUT("/customers/:id", customerController.UpdateCustomer)
	e.PUT("/customers/password/:id", customerController.UpdateCustomerPassword)

	e.GET("/customers/:id/transactions", transactionController.GetAllTransactionByCustomerID)
	e.GET("/customers/:id/transactions/date", transactionController.GetDateRangeTransactionsByCustomerID)

	e.DELETE("/customers/reset", customerController.ResetAllCustomerData)

	// Routes for Generator
	e.GET("/customers/limit/:num", customerController.GetLimitedCustomers)
	e.POST("/customers/multi", customerController.CreateMultiCustomers)

	e.POST("/transactions/multi", transactionController.CreateMultiTransactions)

	// Disabled routes
	// e.DELETE("/customers/:id", customerController.DeleteCustomer)
	// e.POST("/transactions", transactionController.CreateTransaction)
	// e.PUT("/transactions/:id", transactionController.UpdateTransaction)
	// e.DELETE("/transactions/:id", transactionController.DeleteTransaction)

	// Start the server
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
