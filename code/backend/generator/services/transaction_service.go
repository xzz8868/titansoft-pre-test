package services

import (
	"bytes"
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
	GenerateAndSendTransactions(numTransactions int, numCustomers int) error
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
func (ts *transactionService) GenerateAndSendTransactions(numTransactions int, numCustomers int) error {
	// Step 1: Retrieve customer IDs
	customerIDs, err := ts.getCustomerIDs(numCustomers)
	if err != nil {
		return fmt.Errorf("failed to retrieve customer IDs: %w", err)
	}

	if len(customerIDs) == 0 {
		return fmt.Errorf("no customer data")
	}

	// Step 2: Generate transactions
	transactions := ts.generateTransactions(numTransactions, customerIDs)

	// Step 3: Send transactions to backend
	if err := ts.sendTransactions(transactions); err != nil {
		return fmt.Errorf("failed to send transactions: %w", err)
	}

	return nil
}

// getCustomerIDs retrieves customer IDs from the backend server
func (ts *transactionService) getCustomerIDs(numCustomers int) ([]uuid.UUID, error) {
	// Construct the API request
	url := fmt.Sprintf("%s/customers/limit/%d", ts.cfg.BackendServerEndpoint, numCustomers)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check HTTP Status Code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend server responded with status: %d", resp.StatusCode)
	}

	// Analyze response data
	var customers []models.CustomerDTO
	if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
		return nil, err
	}

	// Extract customer ID
	customerIDs := make([]uuid.UUID, len(customers))
	for i, customer := range customers {
		customerIDs[i] = customer.ID
	}

	return customerIDs, nil
}

// generateTransactions creates a list of random transactions
func (ts *transactionService) generateTransactions(numTransactions int, customerIDs []uuid.UUID) []models.TransactionDTO {
	transactions := make([]models.TransactionDTO, numTransactions)
	customerCount := len(customerIDs)
	var wg sync.WaitGroup
	wg.Add(numTransactions)

	// Generate each transaction in a separate goroutine
	for i := 0; i < numTransactions; i++ {
		go func(i int) {
			defer wg.Done()
			transactions[i] = models.TransactionDTO{
				CustomerID: customerIDs[rand.Intn(customerCount)],
				Amount:     rand.Float64() * 1000000, // Random amount up to $1000000
				Time:       ts.randomTimeWithinMonths(18),
			}
		}(i)
	}

	wg.Wait()
	return transactions
}

// sendTransactions posts the transactions to the backend server
func (ts *transactionService) sendTransactions(transactions []models.TransactionDTO) error {
	// Serialize customer data to JSON
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	// Create POST request with transaction data
	url := fmt.Sprintf("%s/transactions/multi", ts.cfg.BackendServerEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(transactionsJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle response based on status code
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("backend server responded with status: %d", resp.StatusCode)
	}

	log.Println("Transactions successfully sent to backend server")
	return nil
}

// randomTimeWithinMonths generates a random time within the past specified months
func (ts *transactionService) randomTimeWithinMonths(months int) time.Time {
	now := time.Now()
	past := now.AddDate(0, -months, 0)
	delta := now.Unix() - past.Unix()
	sec := rand.Int63n(delta) + past.Unix()
	return time.Unix(sec, 0)
}

// nopCloser is a helper to create an io.ReadCloser from bytes.Reader
type nopCloser struct {
	*bytes.Reader
}

func (nopCloser) Close() error { return nil }
