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
	CreateCustomer(customer *models.Customer) error
	CreateMultiCustomers(customers []*models.Customer) (int64, error)
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	UpdatePassword(customer *models.Customer) error
	ResetAllCustomerData() error
	// DeleteCustomer(id uuid.UUID) error
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
func (cr *customerRepository) GetAllCustomers() ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := cr.db.Omit("Password").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// GetLimitedCustomers retrieves a limited number of customers, omitting the Password field
func (cr *customerRepository) GetLimitedCustomers(num int) ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := cr.db.Omit("Password").Limit(num).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// CreateCustomer inserts a single customer into the database
func (cr *customerRepository) CreateCustomer(customer *models.Customer) error {
	return cr.db.Create(customer).Error
}

// CreateMultiCustomers performs batch insert for multiple customers, ignoring duplicates
func (cr *customerRepository) CreateMultiCustomers(customers []*models.Customer) (int64, error) {
	batchSize := 100
	result := cr.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&customers, batchSize)
	return result.RowsAffected, result.Error
}

// GetCustomerByID retrieves a customer by ID, omitting the Password field
func (cr *customerRepository) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := cr.db.Omit("Password").First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// UpdateCustomer updates the Name, Email, and Gender fields of a customer
func (cr *customerRepository) UpdateCustomer(customer *models.Customer) error {
	return cr.db.Model(&customer).Select("Name", "Email", "Gender").Updates(customer).Error
}

// UpdatePassword updates the Password field of a customer
func (cr *customerRepository) UpdatePassword(customer *models.Customer) error {
	return cr.db.Model(&customer).Select("Password").Updates(customer).Error
}

// ResetAllCustomerData deletes all customer records and associated data
func (cr *customerRepository) ResetAllCustomerData() error {
	err := cr.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Customer{}).Error
	if err != nil {
		return err
	}
	// Related data is deleted due to foreign key constraints, if any
	return nil
}

// DeleteCustomer deletes a customer by ID
// func (cr *customerRepository) DeleteCustomer(id uuid.UUID) error {
// 	return cr.db.Delete(&models.Customer{}, "id = ?", id).Error
// }
