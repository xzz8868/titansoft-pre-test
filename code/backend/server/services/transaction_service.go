package services

import (
	"sort"

	"github.com/google/uuid"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/repositories"
)

type TransactionService interface {
	GetTransactionsByCustomerID(id uuid.UUID) ([]*models.TransactionDTO, error)
	GetDateRangeTransactionsByCustomerID(customerID uuid.UUID, from string, to string) ([]*models.TransactionDTO, error)
	CreateMultiTransactions(transactions []*models.TransactionDTO) error
	// CreateTransaction(transaction *models.Transaction) error
	// UpdateTransaction(transaction *models.Transaction) error
	// DeleteTransaction(id uuid.UUID) error
}

type transactionService struct {
	repo repositories.TransactionRepository
}

// Constructor for creating a new TransactionService instance
func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

// Retrieves all transactions for a given customer, sorts them by time, and maps to DTOs
func (cs *transactionService) GetTransactionsByCustomerID(customerID uuid.UUID) ([]*models.TransactionDTO, error) {
	transactions, err := cs.repo.GetTransactionsByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	cs.sortTransactionsByTime(transactions)
	return cs.mapTransactionsToDTOs(transactions), nil
}

// Retrieves transactions within a date range for a customer, sorts them by time, and maps to DTOs
func (cs *transactionService) GetDateRangeTransactionsByCustomerID(customerID uuid.UUID, from string, to string) ([]*models.TransactionDTO, error) {
	transactions, err := cs.repo.GetDateRangeTransactionsByCustomerID(customerID, from, to)
	if err != nil {
		return nil, err
	}

	cs.sortTransactionsByTime(transactions)
	return cs.mapTransactionsToDTOs(transactions), nil
}

// Helper function to sort transactions by time in ascending order
func (cs *transactionService) sortTransactionsByTime(transactions []*models.Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Time.Before(transactions[j].Time)
	})
}

// Helper function to map transaction models to TransactionDTOs and assign sequences
func (cs *transactionService) mapTransactionsToDTOs(transactions []*models.Transaction) []*models.TransactionDTO {
	transactionDTOs := make([]*models.TransactionDTO, len(transactions))
	for i, txn := range transactions {
		transactionDTOs[i] = &models.TransactionDTO{
			ID:         txn.ID,
			CustomerID: txn.CustomerID,
			Amount:     txn.Amount,
			Time:       txn.Time,
			Sequence:   i + 1, // Sequence starts from 1 and increments
		}
	}
	return transactionDTOs
}


// Creates multiple transactions by mapping DTOs to ORM models and saving them in the repository
func (cs *transactionService) CreateMultiTransactions(transactions []*models.TransactionDTO) error {
	var transactionORMs []*models.Transaction
	for _, dto := range transactions {
		
		// Map TransactionDTO to Transaction ORM model
		transactionORM := &models.Transaction{
			ID:         uuid.New(),
			CustomerID: dto.CustomerID,
			Amount:     dto.Amount,
			Time:       dto.Time,
		}
		
		transactionORMs = append(transactionORMs, transactionORM)
	}
	
	// Call the Repository layer to save transactions
	return cs.repo.CreateMultiTransactions(transactionORMs)
}

// Creates a new transaction record in the repository
// func (cs *transactionService) CreateTransaction(transaction *models.Transaction) error {
// 	return cs.repo.Create(transaction)
// }

// Updates an existing transaction record in the repository
// func (cs *transactionService) UpdateTransaction(transaction *models.Transaction) error {
// 	return cs.repo.Update(transaction)
// }

// Deletes a transaction by its ID in the repository
// func (cs *transactionService) DeleteTransaction(id uuid.UUID) error {
// 	return cs.repo.Delete(id)
// }
