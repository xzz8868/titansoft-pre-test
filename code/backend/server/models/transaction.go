package models

import (
    "github.com/google/uuid"
    "time"
)

type Transaction struct {
    ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    CustomerID uuid.UUID `gorm:"type:char(36);not null" json:"customer_id"`
    Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
    Time       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"time"`
    Sequence   int       `gorm:"not null" json:"sequence"`
    CreatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
}
