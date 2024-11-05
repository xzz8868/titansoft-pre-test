package models

import "github.com/google/uuid"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

type CustomerDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Gender   Gender    `json:"gender"`
}
