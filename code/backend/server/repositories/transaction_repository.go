package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
)

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	GetTransactionsByCustomerID(id uuid.UUID) ([]*models.Transaction, error)
	GetDateRangeTransactionsByCustomerID(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error)
	CreateMultiTransactions(transactions []*models.Transaction) error
	GetTotalAmountsByCustomersInPastYear() (map[uuid.UUID]float64, error)
	// CreateTransaction(transaction *models.Transaction) error
	// UpdateTransaction(transaction *models.Transaction) error
	// DeleteTransaction(id uuid.UUID) error
}

// transactionRepository is the concrete implementation of TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository returns a new instance of transactionRepository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// GetTransactionsByCustomerID retrieves all transactions for a specific customer
func (tr *transactionRepository) GetTransactionsByCustomerID(customerID uuid.UUID) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := tr.db.Where("customer_id = ?", customerID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetDateRangeTransactionsByCustomerID retrieves transactions for a customer within a date range
func (tr *transactionRepository) GetDateRangeTransactionsByCustomerID(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := tr.db.Where("customer_id = ? AND time BETWEEN ? AND ?", customerID, from, to).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateMultiTransactions inserts multiple transaction records into the database
func (tr *transactionRepository) CreateMultiTransactions(transactions []*models.Transaction) error {
	batchSize := 100
	return tr.db.CreateInBatches(transactions, batchSize).Error
}

// GetTotalAmountsByCustomersInPastYear calculates the total transaction amounts for each customer in the past year
func (tr *transactionRepository) GetTotalAmountsByCustomersInPastYear() (map[uuid.UUID]float64, error) {
	var results []struct {
		CustomerID  uuid.UUID
		TotalAmount float64
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0).Truncate(24 * time.Hour)
	err := tr.db.Model(&models.Transaction{}).
		Select("customer_id, SUM(amount) as total_amount").
		Where("time >= ?", oneYearAgo).
		Group("customer_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map customer IDs to their respective total transaction amounts
	totalAmounts := make(map[uuid.UUID]float64)
	for _, result := range results {
		totalAmounts[result.CustomerID] = result.TotalAmount
	}

	return totalAmounts, nil
}

// CreateTransaction inserts a new transaction record into the database
// func (tr *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
// 	return tr.db.Create(transaction).Error
// }

// UpdateTransaction modifies an existing transaction record in the database
// func (tr *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
// 	return tr.db.Save(transaction).Error
// }

// DeleteTransaction removes a transaction record by its ID
// func (tr *transactionRepository) DeleteTransaction(id uuid.UUID) error {
// 	return tr.db.Delete(&models.Transaction{}, "id = ?", id).Error
// }
