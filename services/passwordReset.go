package services

import (
	"database/sql"
	"fmt"
	"lenslocked/application/gateway"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/rand"
	"lenslocked/tokenManager"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour
)

const (
	//DefaultSender is the default email address to send emails from.
	DefaultSender = "support@lenslocked.com"
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

func (service *PasswordResetService) Create(email, resetPasswordURL string) (*entity.PasswordReset, error) {
	email = strings.ToLower(email)
	user, err := service.UserRepository.FindByEmail(email)
	if err != nil {
		//TODO: Consider returning a specific erroe when the user does not exist.
		return nil, fmt.Errorf("create: %w", err)
	}
	//Build the passwordReset
	bytesPerToken := service.BytesPerToken
	if bytesPerToken == 0 {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
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
	passwordReset := entity.NewPasswordReset(service.Generate(), user.ID, token, tokenHash, duration)
	err = service.PasswordReset.Create(passwordReset)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	vals := url.Values{
		"token": {passwordReset.Token},
	}
	err = service.forgotPassword(user.Email, resetPasswordURL+vals.Encode())
	if err != nil {
		//deletar o password_reset caso de erro
		return nil, fmt.Errorf("forgot password email: %w", err)
	}
	return passwordReset, nil
}

func (us *PasswordResetService) forgotPassword(to, resetURL string) error {
	email := entity.NewEmail(
		DefaultSender,
		to,
		"Reset your password",
		"To reset your password, please visit the following link: "+resetURL,
		`<p>To reset your password, please visit the following link: <a href="`+resetURL+`">`+resetURL+`</a></p>`,
	)
	if err := us.EmailGateway.Send(email); err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
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
