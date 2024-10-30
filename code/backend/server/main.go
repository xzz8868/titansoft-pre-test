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
	// 加載配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("載入配置失敗：%v", err)
	}

	// 資料庫連線
	dsn := cfg.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("資料庫連線失敗：%v", err)
	}

	// 自動遷移
	if err := db.AutoMigrate(&models.Customer{}, &models.Transaction{}); err != nil {
		log.Fatalf("資料庫遷移失敗：%v", err)
	}

	// 初始化 Repository
	customerRepo := repositories.NewCustomerRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// 初始化 Service
	customerService := services.NewCustomerService(customerRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// 初始化 Controller
	customerController := controllers.NewCustomerController(customerService)
	transactionController := controllers.NewTransactionController(transactionService)

	// 初始化 Echo
	e := echo.New()

	// 添加CORS中間件
	e.Use(middleware.CORS())

	// 設定路由
	e.POST("/customers", customerController.CreateCustomer)
	e.GET("/customers", customerController.GetAllCustomers)
	e.GET("/customers/:id", customerController.GetCustomerByID)
	e.GET("/customers/:id/transactions", transactionController.GetTransactionsByCustomer)
	e.PUT("/customers/:id", customerController.UpdateCustomer)
	e.DELETE("/customers/:id", customerController.DeleteCustomer)

	e.POST("/transactions", transactionController.CreateTransaction)
	e.GET("/transactions/:id", transactionController.GetTransactionByID)
	e.PUT("/transactions/:id", transactionController.UpdateTransaction)
	e.DELETE("/transactions/:id", transactionController.DeleteTransaction)

	// 啟動伺服器
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
