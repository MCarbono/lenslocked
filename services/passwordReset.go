package services

import (
	"database/sql"
	"fmt"
	"lenslocked/application/gateway"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Duration          time.Duration
	TokenManager      tokenManager.Manager
	UserRepository    repository.UserRepository
	PasswordReset     repository.PasswordResetRepository
	EmailGateway      gateway.EmailProvider
	SessionRepository repository.SessionRepository
	idGenerator.IDGenerator
}

// We are going to consume a token and return the session associated with it, or return an error if the token wasn't valid for any reason.
//TODO: Unit of Work
func (service *PasswordResetService) Consume(token, password string) (*entity.Session, error) {
	tokenHash := service.TokenManager.Hash(token)
	pwReset, err := service.PasswordReset.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if pwReset.IsExpired() {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = service.PasswordReset.Delete(pwReset)
	if err != nil {
		return nil, fmt.Errorf("delete password reset: %w", err)
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	err = service.UserRepository.UpdatePasswordHash(pwReset.UserID, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	// Create will create a new session for the user provided. The session token
	// will be returned as the Token field on the Session type, but only the hashed
	// session token is stored in the database.
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < tokenManager.MIN_BYTES_PER_TOKEN {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, tokenHash, err = service.TokenManager.NewToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	insertedSession, err := service.SessionRepository.Upsert(entity.NewSession(service.Generate(), pwReset.UserID, token, tokenHash))
	if err != nil {
		return nil, fmt.Errorf("upsert session: %w", err)
	}
	return insertedSession, err
}
