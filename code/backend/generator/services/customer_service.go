package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/config"
	"github.com/xzz8868/titansoft-pre-test/code/backend/generator/models"
)

const str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// CustomerService defines the interface for customer-related operations
type CustomerService interface {
	GenerateCustomerData(num int) ([]models.CustomerDTO, error)
	CreateMultiCustomersAPICall(customers []models.CustomerDTO) (int, int, error)
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

// GenerateCustomerData generates a list of random customer data
func (cs *customerService) GenerateCustomerData(num int) ([]models.CustomerDTO, error) {
	log.Printf("Generating data for %d customers", num)
	var customers []models.CustomerDTO
	for i := 0; i < num; i++ {
		name := cs.generateRandomName()
		customer := models.CustomerDTO{
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

// CreateMultiCustomersAPICall sends a batch of customer data to the backend API
func (cs *customerService) CreateMultiCustomersAPICall(customers []models.CustomerDTO) (int, int, error) {
    successCount := 0
    failCount := 0

    // Serialize customer data to JSON
    customersJSON, err := json.Marshal(customers)
    if err != nil {
        log.Printf("JSON marshalling error: %v", err)
        return successCount, failCount, err
    }

    // Construct the API request
    url := fmt.Sprintf("%s/customers/multi", cs.cfg.BackendServerEndpoint)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(customersJSON))
    if err != nil {
        log.Printf("Request creation error: %v", err)
        return successCount, failCount, err
    }
    req.Header.Set("Content-Type", "application/json")

    // Execute the HTTP request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("HTTP request error: %v", err)
        return successCount, failCount, err
    }
    defer resp.Body.Close()

    // Handle response based on status code
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


// generateRandomName generates a random 8-character string as a name
func (cs *customerService) generateRandomName() string {
	return generateRandomString(8)
}


// generateRandomPassword generates a random 16-character password
func (cs *customerService) generateRandomPassword() string {
	return generateRandomString(16)
}

// generateRandomString generates random string
func generateRandomString(n int) string {
	strLen := len(str)
	result := make([]byte, n)
	bytes := []byte(str)
	for i := 0; i < n; i++ {
		result[i] = bytes[rand.Intn(strLen)]
	}
	return string(result)
}

// generateRandomEmail creates an email address by combining a name with a random domain
func (cs *customerService) generateRandomEmail(name string) string {
	email := name + emailDomains[rand.Intn(len(emailDomains))]
	log.Printf("Generated email: %s", email)
	return email
}

// randomGender randomly selects a gender from predefined options
func (cs *customerService) randomGender() models.Gender {
	genders := []models.Gender{models.Male, models.Female, models.Other}
	return genders[rand.Intn(len(genders))]
}
