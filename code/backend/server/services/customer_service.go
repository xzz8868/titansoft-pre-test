package services

import (
	"github.com/google/uuid"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
	"github.com/xzz8868/titansoft-pre-test/code/backend/server/repositories"
)

type CustomerService interface {
	GetAllCustomers() ([]*models.Customer, error)
	CreateCustomer(customer *models.Customer) error
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id uuid.UUID) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func (s *customerService) GetAllCustomers() ([]*models.Customer, error) {
	return s.repo.GetAll()
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo}
}

func (s *customerService) CreateCustomer(customer *models.Customer) error {
	return s.repo.Create(customer)
}

func (s *customerService) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *customerService) UpdateCustomer(customer *models.Customer) error {
	return s.repo.Update(customer)
}

func (s *customerService) DeleteCustomer(id uuid.UUID) error {
	return s.repo.Delete(id)
}
