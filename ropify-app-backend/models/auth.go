package models

import (
	"context"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type AuthCredentials struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"password" validate:"required"`
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *AuthCredentials) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
	RegisterOAuthUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

type AuthService interface {
	Login(ctx context.Context, loginData *AuthCredentials) (string, *User, error)
	Register(ctx context.Context, registerData *AuthCredentials) (string, *User, error)
}

// Check if a password matches a hash
func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Checks if an email is valid
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
