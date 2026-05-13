package domain

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func NewUser(email, passwordHash string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(email) {
		return nil, errors.New("invalid email address")
	}
	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}
	return &User{Email: email, PasswordHash: passwordHash}, nil
}

func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 {
		return false
	}
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")
	return dot > 0 && dot < len(domain)-1
}
