package services

import (
	"encoding/base64"

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

	totalAmounts, err := s.transactionRepo.GetTotalAmountsByCustomersInPastYear()
	if err != nil {
		return nil, err
	}

	var customerDTOs []*models.CustomerDTO
	for _, customer := range customers {
		totalAmount := totalAmounts[customer.ID] // Default to 0 if no record found

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

// GetLimitedCustomers retrieves a limited number of customers and their total transaction amounts in the past year.
func (s *customerService) GetLimitedCustomers(num int) ([]*models.CustomerDTO, error) {
	customers, err := s.customerRepo.GetLimitedCustomers(num)
	if err != nil {
		return nil, err
	}

	totalAmounts, err := s.transactionRepo.GetTotalAmountsByCustomersInPastYear()
	if err != nil {
		return nil, err
	}

	var customerDTOs []*models.CustomerDTO
	for _, customer := range customers {
		totalAmount := totalAmounts[customer.ID] // Default to 0 if no record found

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
	return s.customerRepo.CreateCustomers(customer)
}

// CreateMultiCustomers hashes passwords for multiple customers and saves them in batch.
// Returns the count of successful and failed creations.
func (s *customerService) CreateMultiCustomers(customers []*models.Customer) (int, int, error) {
	successCount := 0
	failCount := 0
	validCustomers := make([]*models.Customer, 0, len(customers))

	for _, customer := range customers {
		if len(customer.Password) < 8 {
			failCount++
			continue // Skip customers with insufficient password length
		}
		customer.ID = uuid.New()
		hashedPassword, err := s.hashPassword(customer.Password)
		if err != nil {
			failCount++
			continue
		}
		customer.Password = hashedPassword
		validCustomers = append(validCustomers, customer)
	}

	if len(validCustomers) == 0 {
		return successCount, failCount, nil
	}

	rowsAffected, err := s.customerRepo.CreateMultiCustomers(validCustomers)
	if err != nil {
		// Handle potential errors such as database connection issues
		failCount += len(validCustomers) - int(rowsAffected)
		successCount += int(rowsAffected)
		return successCount, failCount, err
	}

	successCount += int(rowsAffected)
	failCount += len(validCustomers) - int(rowsAffected)

	return successCount, failCount, nil
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
