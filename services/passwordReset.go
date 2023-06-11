package services

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
	"lenslocked/domain/repository"
	"lenslocked/gateway"
	"lenslocked/rand"
	"lenslocked/token"
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

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration       time.Duration
	TokenManager   token.Manager
	UserRepository repository.UserRepository
	PasswordReset  repository.PasswordResetRepository
	EmailGateway   gateway.EmailProvider
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
	passwordReset := entity.NewPasswordReset(user.ID, token, tokenHash, duration)
	id, err := service.PasswordReset.Create(passwordReset)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	passwordReset.ID = id
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

// We are going to consume a token and return the user associated with it, or return an error if the token wasn't valid for any reason.
func (service *PasswordResetService) Consume(token, password string) (*entity.User, error) {
	tokenHash := service.TokenManager.Hash(token)
	// var user entity.User
	// var pwReset entity.PasswordReset
	// row := service.DB.QueryRow(`
	// 	SELECT password_resets.id,
	// 		password_resets.expires_at,
	// 		users.id,
	// 		users.email,
	// 		users.password_hash
	// 	FROM password_resets
	// 		JOIN users ON users.id = password_resets.user_id
	// 	WHERE password_resets.token_hash = $1;`, tokenHash)
	// err := row.Scan(
	// 	&pwReset.ID, &pwReset.ExpiresAt,
	// 	&user.ID, &user.Email, &user.PasswordHash)
	pwReset, err := service.PasswordReset.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	user, err := service.UserRepository.FindByID(pwReset.UserID)
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
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = service.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1`, user.ID, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	return user, nil
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
