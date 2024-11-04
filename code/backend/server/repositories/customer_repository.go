package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	GetAllCustomers() ([]*models.Customer, error)
	GetLimitedCustomers(num int) ([]*models.Customer, error)
	CreateCustomers(customer *models.Customer) error
	CreateMultiCustomers(customers []*models.Customer) (int64, error)
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	UpdatePassword(customer *models.Customer) error
	DeleteCustomer(id uuid.UUID) error
	ResetAllCustomerData() error
}

// customerRepository implements CustomerRepository using Gorm
type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customerRepository instance
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

// GetAllCustomers retrieves all customers, omitting the Password field
func (r *customerRepository) GetAllCustomers() ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := r.db.Omit("Password").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// GetLimitedCustomers retrieves a limited number of customers, omitting the Password field
func (r *customerRepository) GetLimitedCustomers(num int) ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := r.db.Omit("Password").Limit(num).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// CreateCustomers inserts a single customer into the database
func (r *customerRepository) CreateCustomers(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

// CreateMultiCustomers performs batch insert for multiple customers, ignoring duplicates
func (r *customerRepository) CreateMultiCustomers(customers []*models.Customer) (int64, error) {
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&customers)
	return result.RowsAffected, result.Error
}

// GetCustomerByID retrieves a customer by ID, omitting the Password field
func (r *customerRepository) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Omit("Password").First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// UpdateCustomer updates the Name, Email, and Gender fields of a customer
func (r *customerRepository) UpdateCustomer(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Name", "Email", "Gender").Updates(customer).Error
}

// UpdatePassword updates the Password field of a customer
func (r *customerRepository) UpdatePassword(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Password").Updates(customer).Error
}

// DeleteCustomer deletes a customer by ID
func (r *customerRepository) DeleteCustomer(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}

// ResetAllCustomerData deletes all customer records and associated data
func (r *customerRepository) ResetAllCustomerData() error {
	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Customer{}).Error
	if err != nil {
		return err
	}
	// Related data is deleted due to foreign key constraints, if any
	return nil
}
