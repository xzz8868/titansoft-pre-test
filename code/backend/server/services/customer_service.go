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
	CreateCustomer(customer *models.Customer) error
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	UpdateCustomerPassword(customer *models.Customer) error
	DeleteCustomer(id uuid.UUID) error
}

type customerService struct {
	repo            repositories.CustomerRepository
	transactionRepo repositories.TransactionRepository
	salt            string
}

func NewCustomerService(repo repositories.CustomerRepository, transactionRepo repositories.TransactionRepository, salt string) CustomerService {
	return &customerService{
		repo:            repo,
		transactionRepo: transactionRepo,
		salt:            salt,
	}
}

func (s *customerService) GetAllCustomers() ([]*models.CustomerDTO, error) {
	customers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	totalAmounts, err := s.transactionRepo.GetTotalAmountsByCustomersInPastYear()
	if err != nil {
		return nil, err
	}

	var customerDTOs []*models.CustomerDTO
	for _, customer := range customers {
		totalAmount := totalAmounts[customer.ID] // 如果没有记录，默认为0

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

func (s *customerService) CreateCustomer(customer *models.Customer) error {
	hashedPassword, err := s.hashPassword(customer.Password)
	if err != nil {
		return err
	}
	customer.Password = hashedPassword
	return s.repo.Create(customer)
}

func (s *customerService) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *customerService) UpdateCustomer(customer *models.Customer) error {
	return s.repo.Update(customer)
}

func (s *customerService) UpdateCustomerPassword(customer *models.Customer) error {
	if customer.Password != "" {
		hashedPassword, err := s.hashPassword(customer.Password)
		if err != nil {
			return err
		}
		customer.Password = hashedPassword
	}
	return s.repo.UpdatePassword(customer)
}

func (s *customerService) DeleteCustomer(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *customerService) hashPassword(password string) (string, error) {
	salt := []byte(s.salt)
	dk, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}
