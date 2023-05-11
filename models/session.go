package models

import (
	"database/sql"
	"fmt"

	"github.com/alorents/lenslocked/rand"
)

const (
	// MinBytesPerToken is the minimum number of bytes to use for a session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session
	// Only the TokenHash will be saved in the DB
	// So at all other times Token will be empty
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
}

func (s *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := s.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO hash the token
	session := Session{
		UserID: userID,
		Token:  token,
		// TODO set the TokenHash
	}
	// TODO Implement SessionService.Create
	return &session, nil
}

func (s *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User
	return nil, nil
}
