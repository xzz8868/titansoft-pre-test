package models

import (
	"github.com/google/uuid"
)

type CustomerDTO struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	Email                  string    `json:"email"`
	Gender                 Gender    `json:"gender"`
	TotalTransactionAmount float64   `json:"total_transaction_amount"`
}
