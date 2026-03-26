package model

import (
	"errors"
	"regexp"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validation errors
var (
	ErrInvalidName  = errors.New("name must be between 2 and 100 characters")
	ErrInvalidEmail = errors.New("email must be a valid format")
	ErrInvalidAge   = errors.New("age must be between 18 and 120")
)

// Validate performs validation on the user model
func (u *User) Validate() error {
	if len(u.Name) < 2 || len(u.Name) > 100 {
		return ErrInvalidName
	}

	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}

	if u.Age < 18 || u.Age > 120 {
		return ErrInvalidAge
	}

	return nil
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

// UpdateUser represents partial user update
type UpdateUser struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}