package models

import (
	"github.com/google/uuid"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

type Customer struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	Password string    `gorm:"type:varchar(255);not null" json:"password"`
	Email    string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Gender   Gender    `gorm:"type:enum('male','female','other');not null" json:"gender"`
}
