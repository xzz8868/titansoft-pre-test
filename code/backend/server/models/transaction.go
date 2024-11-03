package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:char(36);not null;index" json:"customer_id"`
	Customer   Customer  `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Time       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"time"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
}
