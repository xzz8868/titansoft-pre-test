package services

import (
	"github.com/google/uuid"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/repositories"
)

type TransactionService interface {
	GetTransactionsByCustomer(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error)
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uuid.UUID) error
}

type transactionService struct {
	repo repositories.TransactionRepository
}

func (s *transactionService) GetTransactionsByCustomer(customerID uuid.UUID, from string, to string) ([]*models.Transaction, error) {
	return s.repo.GetByCustomerIDAndDateRange(customerID, from, to)
}

func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
	return s.repo.Create(transaction)
}

func (s *transactionService) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetByID(id)
}

func (s *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	return s.repo.Update(transaction)
}

func (s *transactionService) DeleteTransaction(id uuid.UUID) error {
	return s.repo.Delete(id)
}
