package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDTO struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Sequence   int       `json:"sequence"`
	Time       time.Time `json:"time"`
}
