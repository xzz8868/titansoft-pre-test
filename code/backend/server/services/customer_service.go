package services

import (
	"encoding/base64"
	"fmt"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/repositories"
)

type CustomerService interface {
	GetAllCustomers() ([]*models.CustomerDTO, error)
	GetLimitedCustomers(num int) ([]*models.CustomerDTO, error)
	CreateCustomer(customer *models.Customer) error
	CreateMultiCustomers(customers []*models.Customer) (int, int, error)
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	UpdateCustomerPassword(customer *models.Customer) error
	DeleteCustomer(id uuid.UUID) error
	ResetAllCustomerData() error
}

type customerService struct {
	customerRepo    repositories.CustomerRepository
	transactionRepo repositories.TransactionRepository
	salt            string
}

// NewCustomerService creates a new instance of CustomerService with required dependencies.
func NewCustomerService(repo repositories.CustomerRepository, transactionRepo repositories.TransactionRepository, salt string) CustomerService {
	return &customerService{
		customerRepo:    repo,
		transactionRepo: transactionRepo,
		salt:            salt,
	}
}

// GetAllCustomers retrieves all customers and their total transaction amounts in the past year.
func (s *customerService) GetAllCustomers() ([]*models.CustomerDTO, error) {
	customers, err := s.customerRepo.GetAllCustomers()
	if err != nil {
		return nil, err
	}

	customerDTOs, err := s.buildCustomerDTOsWithTransactions(customers)
	if err != nil {
		return nil, err
	}

	return customerDTOs, nil
}

// GetLimitedCustomers retrieves a specified number of customers and their total transaction amounts for the past year.
func (s *customerService) GetLimitedCustomers(num int) ([]*models.CustomerDTO, error) {
	// Fetch a limited number of customers from the repository
	customers, err := s.customerRepo.GetLimitedCustomers(num)
	if err != nil {
		return nil, err
	}

	// Build and return customer DTOs enriched with transaction data
	customerDTOs, err := s.buildCustomerDTOsWithTransactions(customers)
	if err != nil {
		return nil, err
	}

	return customerDTOs, nil
}

// buildCustomerDTOsWithTransactions constructs CustomerDTOs with total transaction amounts from the past year.
func (s *customerService) buildCustomerDTOsWithTransactions(customers []*models.Customer) ([]*models.CustomerDTO, error) {
	// Retrieve total transaction amounts for each customer from the past year
	totalAmounts, err := s.transactionRepo.GetTotalAmountsByCustomersInPastYear()
	if err != nil {
		return nil, err
	}

	var customerDTOs []*models.CustomerDTO
	// Map each customer to a DTO, attaching their transaction total
	for _, customer := range customers {
		totalAmount := totalAmounts[customer.ID] // Default to zero if not found in map
		customerDTO := &models.CustomerDTO{
			ID:                     customer.ID,
			Name:                   customer.Name,
			Email:                  customer.Email,
			Gender:                 customer.Gender,
			TotalTransactionAmount: totalAmount,
		}
		customerDTOs = append(customerDTOs, customerDTO)
	}
	return customerDTOs, nil
}

// CreateCustomer hashes the customer's password and saves the customer to the repository.
func (s *customerService) CreateCustomer(customer *models.Customer) error {
	hashedPassword, err := s.hashPassword(customer.Password)
	if err != nil {
		return err
	}
	customer.Password = hashedPassword
	return s.customerRepo.CreateCustomer(customer)
}

// CreateMultiCustomers hashes passwords for multiple customers and saves them in batch.
// Returns the count of successful and failed creations.
func (s *customerService) CreateMultiCustomers(customers []*models.Customer) (int, int, error) {
	successCount := 0
	failCount := 0
	validCustomers := make([]*models.Customer, 0, len(customers))

	type result struct {
		customer *models.Customer
		err      error
	}

	results := make(chan result, len(customers))
	var wg sync.WaitGroup

	// Set max concurrency
	maxGoroutines := runtime.GOMAXPROCS(2)
	sem := make(chan struct{}, maxGoroutines)

	for _, customer := range customers {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore slot for goroutine
		go func(c *models.Customer) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore slot

			// Validate password length
			if len(c.Password) < 8 {
				results <- result{nil, fmt.Errorf("password too short for customer %v", c.Email)}
				return
			}

			// Generate UUID and hash password
			c.ID = uuid.New()
			hashedPassword, err := s.hashPassword(c.Password)
			if err != nil {
				results <- result{nil, fmt.Errorf("failed to hash password for customer %v: %w", c.Email, err)}
				return
			}
			c.Password = hashedPassword
			results <- result{c, nil}
		}(customer)
	}

	// Close results channel once all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from goroutines
	for res := range results {
		if res.err != nil {
			failCount++
		} else {
			validCustomers = append(validCustomers, res.customer)
			successCount++
		}
	}

	if len(validCustomers) == 0 {
		return successCount, failCount, nil
	}

	// Batch insert valid customers into the database
	rowsAffected, err := s.customerRepo.CreateMultiCustomers(validCustomers)
	if err != nil {
		// Assuming rowsAffected is accurate, adjust successCount and failCount accordingly
		failCount += len(validCustomers) - int(rowsAffected)
		return int(rowsAffected), failCount, fmt.Errorf("batch insert error: %w", err)
	}

	return int(rowsAffected), failCount, nil
}

// GetCustomerByID retrieves a customer by their unique ID.
func (s *customerService) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	return s.customerRepo.GetCustomerByID(id)
}

// UpdateCustomer updates the customer's information in the repository.
func (s *customerService) UpdateCustomer(customer *models.Customer) error {
	return s.customerRepo.UpdateCustomer(customer)
}

// UpdateCustomerPassword hashes the new password (if provided) and updates it in the repository.
func (s *customerService) UpdateCustomerPassword(customer *models.Customer) error {
	if customer.Password != "" {
		hashedPassword, err := s.hashPassword(customer.Password)
		if err != nil {
			return err
		}
		customer.Password = hashedPassword
	}
	return s.customerRepo.UpdatePassword(customer)
}

// DeleteCustomer removes a customer from the repository by their unique ID.
func (s *customerService) DeleteCustomer(id uuid.UUID) error {
	return s.customerRepo.DeleteCustomer(id)
}

// ResetAllCustomerData clears all customer data in the repository.
func (s *customerService) ResetAllCustomerData() error {
	return s.customerRepo.ResetAllCustomerData()
}

// hashPassword hashes a password using the scrypt algorithm and encodes it in base64.
func (s *customerService) hashPassword(password string) (string, error) {
	salt := []byte(s.salt)
	dk, err := scrypt.Key([]byte(password), salt, 1024, 8, 1, 32) // N=1024 for demonstration, higher values recommended
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}
