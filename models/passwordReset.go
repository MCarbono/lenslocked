package models

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
	"lenslocked/rand"
	"lenslocked/token"
	"strings"
	"time"
)

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour
)

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration     time.Duration
	TokenManager token.Manager
}

func (service *PasswordResetService) Create(email string) (*entity.PasswordReset, error) {
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		//TODO: Consider returning a specific erroe when the user does not exist.
		return nil, fmt.Errorf("create: %w", err)
	}
	//Build the passwordReset
	bytesPerToken := service.BytesPerToken
	if bytesPerToken == 0 {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	tokenHash := service.TokenManager.Hash(token)
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	passwordReset := entity.NewPasswordReset(userID, token, tokenHash, duration)
	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		values ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	err = row.Scan(&passwordReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return passwordReset, nil
}

// We are going to consume a token and return the user associated with it, or return an error if the token wasn't valid for any reason.
func (service *PasswordResetService) Consume(token string) (*entity.User, error) {
	tokenHash := service.TokenManager.Hash(token)
	var user entity.User
	var pwReset entity.PasswordReset
	row := service.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email,
			users.password_hash
		FROM password_resets
			JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = service.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
