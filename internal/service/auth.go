package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"vibecheck/internal/domain"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrWeakPassword = errors.New("password must be at least 8 characters")

type AuthService struct {
	users    domain.UserRepository
	sessions domain.SessionRepository
}

func NewAuth(users domain.UserRepository, sessions domain.SessionRepository) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) Register(email, password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u, err := domain.NewUser(email, string(hash))
	if err != nil {
		return err
	}
	_, err = s.users.CreateUser(u)
	return err
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.users.GetUserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	if _, err = s.sessions.CreateSession(user.ID, token, time.Now().Add(30*24*time.Hour)); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessions.DeleteSession(token)
}
