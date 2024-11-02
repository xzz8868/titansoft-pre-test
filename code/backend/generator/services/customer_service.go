package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/models"
)

// CustomerService defines the interface for customer-related operations
type CustomerService interface {
	GenerateCustomerData(num int) ([]models.Customer, error)
	CreateMultiCustomersAPICall(ctx context.Context, customers []models.Customer) (int, int, error)
}

// customerService is the concrete implementation of CustomerService
type customerService struct {
	cfg *config.Config
}

// NewCustomerService is the factory function that returns a CustomerService interface
func NewCustomerService(cfg *config.Config) CustomerService {
	return &customerService{cfg: cfg}
}

var emailDomains = []string{"@gmail.com", "@yahoo.com.tw", "@outlook.com", "@icloud.com", "@hotmail.com",
	"@aol.com", "@mail.com", "@yandex.com", "@protonmail.com", "@gmx.com"}

// GenerateCustomerData generates random customer data
func (cs *customerService) GenerateCustomerData(num int) ([]models.Customer, error) {
	log.Printf("Generating data for %d customers", num)
	var customers []models.Customer
	for i := 0; i < num; i++ {
		name := cs.generateRandomName()
		customer := models.Customer{
			Name:     name,
			Password: cs.generateRandomPassword(),
			Email:    cs.generateRandomEmail(name),
			Gender:   cs.randomGender(),
		}
		customers = append(customers, customer)
	}
	log.Println("Customer data generation completed")
	return customers, nil
}

// CreateMultiCustomersAPICall sends multiple customer data to the backend API
func (cs *customerService) CreateMultiCustomersAPICall(ctx context.Context, customers []models.Customer) (int, int, error) {
	client := &http.Client{}
	successCount := 0
	failCount := 0

	customersJson, err := json.Marshal(customers)
	if err != nil {
		log.Printf("JSON marshalling error: %v", err)
		return successCount, failCount, err
	}

	url := fmt.Sprintf("%s/customers/multi", cs.cfg.BackendServerEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(customersJson))
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return successCount, failCount, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return successCount, failCount, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]int
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("JSON decoding error: %v", err)
			return successCount, failCount, err
		}
		successCount = result["successCount"]
		failCount = result["failCount"]
	} else {
		log.Printf("HTTP response status error: %d", resp.StatusCode)
		return successCount, failCount, fmt.Errorf("failed to create customers with status code: %d", resp.StatusCode)
	}

	log.Printf("API calls completed with %d successes and %d failures", successCount, failCount)
	return successCount, failCount, nil
}

func (cs *customerService) generateRandomName() string {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Error generating random name: %v", err)
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func (cs *customerService) generateRandomPassword() string {
	bytes := make([]byte, 8) // Generates a random 16 character password
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Error generating random password: %v", err)
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func (cs *customerService) generateRandomEmail(name string) string {
	email := name + emailDomains[rand.Intn(len(emailDomains))]
	log.Printf("Generated email: %s", email)
	return email
}

func (cs *customerService) randomGender() models.Gender {
	genders := []models.Gender{models.Male, models.Female, models.Other}
	return genders[rand.Intn(len(genders))]
}
