package domain

import (
	"errors"
	"time"
)

var ErrEmailTaken = errors.New("email already registered")
var ErrDuplicateEntry = errors.New("entry already exists for this date")

type UserRepository interface {
	CreateUser(u *User) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type SessionRepository interface {
	CreateSession(userID, token string, expiresAt time.Time) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	DeleteSession(token string) error
}

type EntryRepository interface {
	CreateEntry(e *Entry) (*Entry, error)
}
