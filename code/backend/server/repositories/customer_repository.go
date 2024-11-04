package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
)

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

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) GetAllCustomers() ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := r.db.Omit("Password").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) GetLimitedCustomers(num int) ([]*models.Customer, error) {
	var customers []*models.Customer
	if err := r.db.Omit("Password").Limit(num).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) CreateCustomers(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) CreateMultiCustomers(customers []*models.Customer) (int64, error) {
	// Perform batch insert while ignoring duplicate entries
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&customers)
	return result.RowsAffected, result.Error
}

func (r *customerRepository) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Omit("Password").First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) UpdateCustomer(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Name", "Email", "Gender").Updates(customer).Error
}

func (r *customerRepository) UpdatePassword(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Password").Updates(customer).Error
}

func (r *customerRepository) DeleteCustomer(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) ResetAllCustomerData() error {
	// Delete all customers and associated data
	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Customer{}).Error
	if err != nil {
		return err
	}
	// If there are associated models, they will be deleted due to foreign key constraints
	return nil
}
