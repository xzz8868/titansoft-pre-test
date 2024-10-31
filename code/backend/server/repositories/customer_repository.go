package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xzz8868/titansoft-pre-test/code/backend/server/models"
)

type CustomerRepository interface {
	GetAll() ([]*models.Customer, error)
	Create(customer *models.Customer) error
	GetByID(id uuid.UUID) (*models.Customer, error)
	Update(customer *models.Customer) error
	UpdatePassword(customer *models.Customer) error
	Delete(id uuid.UUID) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) GetAll() ([]*models.Customer, error) {
    var customers []*models.Customer
    if err := r.db.Omit("Password").Find(&customers).Error; err != nil {
        return nil, err
    }
    return customers, nil
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
    var customer models.Customer
    if err := r.db.Omit("Password").First(&customer, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &customer, nil
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Name", "Email", "Gender").Updates(customer).Error
}

func (r *customerRepository) UpdatePassword(customer *models.Customer) error {
	return r.db.Model(&customer).Select("Password").Updates(customer).Error
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}
