package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
)

type TransactionRepository interface {
	GetByCustomerIDAndDateRange(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error)
	Create(transaction *models.Transaction) error
	GetByID(id uuid.UUID) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uuid.UUID) error
	GetTotalAmountsByCustomersInPastYear() (map[uuid.UUID]float64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (r *transactionRepository) GetByCustomerIDAndDateRange(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("customer_id = ? AND time BETWEEN ? AND ?", customerID, from, to).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) GetTotalAmountsByCustomersInPastYear() (map[uuid.UUID]float64, error) {
	var results []struct {
		CustomerID  uuid.UUID
		TotalAmount float64
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	err := r.db.Model(&models.Transaction{}).
		Select("customer_id, SUM(amount) as total_amount").
		Where("time >= ?", oneYearAgo).
		Group("customer_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	totalAmounts := make(map[uuid.UUID]float64)
	for _, result := range results {
		totalAmounts[result.CustomerID] = result.TotalAmount
	}

	return totalAmounts, nil
}
