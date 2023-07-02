package usecases

import (
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
)

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour

	//DefaultSender is the default email address to send emails from.
	DefaultSender = "support@lenslocked.com"
)

type ForgotPasswordUseCase struct {
	userRepository          repository.UserRepository
	passwordResetRepository repository.PasswordResetRepository
	emailGateway            gateway.EmailProvider
	idGenerator             idGenerator.IDGenerator
	tokenManager            tokenManager.Manager
	duration                time.Duration
	bytesPerToken           int
}

func NewForgotPasswordUseCase(
	userRepository repository.UserRepository,
	passwordResetRepository repository.PasswordResetRepository,
	emailGateway gateway.EmailProvider,
	idGenerator idGenerator.IDGenerator,
	tokenManager tokenManager.Manager) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		userRepository:          userRepository,
		passwordResetRepository: passwordResetRepository,
		emailGateway:            emailGateway,
		idGenerator:             idGenerator,
		tokenManager:            tokenManager,
	}
}

func (uc *ForgotPasswordUseCase) Execute(input *ForgotPasswordInput) (*entity.PasswordReset, error) {
	//Build the passwordReset
	bytesPerToken := uc.bytesPerToken
	if bytesPerToken == 0 {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	tokenHash := uc.tokenManager.Hash(token)
	duration := uc.duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	input.Email = strings.ToLower(input.Email)
	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		//TODO: Consider returning a specific erroe when the user does not exist.
		return nil, fmt.Errorf("create: %w", err)
	}
	passwordReset := entity.NewPasswordReset(uc.idGenerator.Generate(), user.ID, token, tokenHash, duration)
	err = uc.passwordResetRepository.Create(passwordReset)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	vals := url.Values{
		"token": {passwordReset.Token},
	}
	err = uc.sendEmail(user.Email, input.ResetPasswordURL+vals.Encode())
	if err != nil {
		//deletar o password_reset caso de erro
		return nil, fmt.Errorf("forgot password email: %w", err)
	}
	return passwordReset, nil
}

func (uc *ForgotPasswordUseCase) sendEmail(to, resetURL string) error {
	email := entity.NewEmail(
		DefaultSender,
		to,
		"Reset your password",
		"To reset your password, please visit the following link: "+resetURL,
		`<p>To reset your password, please visit the following link: <a href="`+resetURL+`">`+resetURL+`</a></p>`,
	)
	if err := uc.emailGateway.Send(email); err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}

type ForgotPasswordInput struct {
	Email            string
	ResetPasswordURL string
}
