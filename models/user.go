package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

// Create will create the provided user and store it in the database
// the password will be hashed before it is stored
// email must be unique
func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id;`, user.Email, user.PasswordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}
