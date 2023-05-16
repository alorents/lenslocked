package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/alorents/lenslocked/rand"
)

const (
	DefaultResetDuration = time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a password reset is created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset token is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// Verify we have a valid email address and get the user ID
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Build the PasswordReset
	bytesPerToken := service.BytesPerToken
	if bytesPerToken == 0 {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	// Insert the PasswordReset into the database
	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET token_hash = $2,
				expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	// Get the User and PasswordReset from the database
	var (
		user    User
		pwReset PasswordReset
	)

	row := service.DB.QueryRow(`
		SELECT users.id, password_resets.id, password_resets.expires_at
		FROM password_resets join users on password_resets.user_id = users.id
		WHERE token_hash = $1`, service.hash(token))
	err := row.Scan(&user.ID, &pwReset.ID, &pwReset.ExpiresAt)
	if err != nil {
		switch {
		case err == sql.ErrNoRows: // TODO return a custom error to the user
			return nil, fmt.Errorf("consume: %w", err)
		default:
			return nil, fmt.Errorf("consume: %w", err)
		}
	}

	// Check if the PasswordReset has expired
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("consume: password reset token has expired %v", token)
	}

	// Delete the PasswordReset from the database
	_, err = service.DB.Exec(`
		DELETE FROM password_resets WHERE id = $1`, pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
