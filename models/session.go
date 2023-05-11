package models

import (
	"database/sql"
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
}

func (s *SessionService) Create(userID int) (*Session, error) {
	// TODO create the sesssion token

	// TODO Implement SessionService.Create
	return nil, nil
}

func (s *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User
	return nil, nil
}
