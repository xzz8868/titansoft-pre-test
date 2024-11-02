package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/models"
)

// TransactionService defines the interface for transaction-related operations
type TransactionService interface {
	GenerateAndSendTransactions(ctx context.Context, numTransactions int, numCustomers int) error
}

// transactionService is the concrete implementation of TransactionService
type transactionService struct {
	cfg *config.Config
}

// NewTransactionService is the factory function that returns a TransactionService interface
func NewTransactionService(cfg *config.Config) TransactionService {
	return &transactionService{cfg: cfg}
}

// GenerateAndSendTransactions generates transaction data and sends it to the backend server
func (ts *transactionService) GenerateAndSendTransactions(ctx context.Context, numTransactions int, numCustomers int) error {
	// Step 1: Retrieve customer IDs
	customerIDs, err := ts.getCustomerIDs(ctx, numCustomers)
	if err != nil {
		return fmt.Errorf("failed to retrieve customer IDs: %w", err)
	}

	// Step 2: Generate transactions
	transactions := ts.generateTransactions(numTransactions, customerIDs)

	// Step 3: Send transactions to backend
	if err := ts.sendTransactions(ctx, transactions); err != nil {
		return fmt.Errorf("failed to send transactions: %w", err)
	}

	return nil
}

// getCustomerIDs retrieves customer IDs from the backend server
func (ts *transactionService) getCustomerIDs(ctx context.Context, numCustomers int) ([]uuid.UUID, error) {
	url := fmt.Sprintf("%s/customers/limit/%d", ts.cfg.BackendServerEndpoint, numCustomers)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend server responded with status: %d", resp.StatusCode)
	}

	var customers []models.Customer
	if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
		return nil, err
	}

	customerIDs := make([]uuid.UUID, len(customers))
	for i, customer := range customers {
		customerIDs[i] = customer.ID
	}

	return customerIDs, nil
}

// generateTransactions creates a list of random transactions
func (ts *transactionService) generateTransactions(numTransactions int, customerIDs []uuid.UUID) []models.Transaction {
	transactions := make([]models.Transaction, numTransactions)
	customerCount := len(customerIDs)
	var wg sync.WaitGroup
	wg.Add(numTransactions)

	for i := 0; i < numTransactions; i++ {
		go func(i int) {
			defer wg.Done()
			transactions[i] = models.Transaction{
				CustomerID: customerIDs[rand.Intn(customerCount)],
				Amount:     rand.Float64() * 1000000, // Random amount up to $1000000
				Time:       randomTimeWithinMonths(18),
			}
		}(i)
	}

	wg.Wait()
	return transactions
}

// sendTransactions posts the transactions to the backend server
func (ts *transactionService) sendTransactions(ctx context.Context, transactions []models.Transaction) error {
	url := fmt.Sprintf("%s/transactions/multi", ts.cfg.BackendServerEndpoint)
	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}
	req.Body = http.NoBody
	req.Body = nopCloser{bytes.NewReader(data)}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("backend server responded with status: %d", resp.StatusCode)
	}

	log.Println("Transactions successfully sent to backend server")
	return nil
}

// randomTimeWithinMonths generates a random time within the past specified months
func randomTimeWithinMonths(months int) time.Time {
	now := time.Now()
	past := now.AddDate(0, -months, 0)
	delta := now.Unix() - past.Unix()
	sec := rand.Int63n(delta) + past.Unix()
	return time.Unix(sec, 0)
}

// Helper to create an io.ReadCloser from bytes.Reader
type nopCloser struct {
	*bytes.Reader
}

func (nopCloser) Close() error { return nil }
