package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDTO struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
}
