package domain

import "time"

type Users struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
