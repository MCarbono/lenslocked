package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"

	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordUseCase struct {
	userRepository          repository.UserRepository
	passwordResetRepository repository.PasswordResetRepository
	sessionRepository       repository.SessionRepository
	idGenerator             idGenerator.IDGenerator
	tokenManager            tokenManager.Manager
	bytesPerToken           int
}

func NewResetPasswordUseCase(
	userRepository repository.UserRepository,
	passwordResetRepository repository.PasswordResetRepository,
	sessionRepository repository.SessionRepository,
	idGenerator idGenerator.IDGenerator,
	tokenManager tokenManager.Manager) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepository:          userRepository,
		passwordResetRepository: passwordResetRepository,
		sessionRepository:       sessionRepository,
		idGenerator:             idGenerator,
		tokenManager:            tokenManager,
	}
}

// We are going to consume a token and return the session associated with it, or return an error if the token wasn't valid for any reason.
//TODO: Unit of Work
func (uc *ResetPasswordUseCase) Execute(input *ResetPasswordInput) (*entity.Session, error) {
	tokenHash := uc.tokenManager.Hash(input.Token)
	pwReset, err := uc.passwordResetRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if pwReset.IsExpired() {
		return nil, fmt.Errorf("token expired: %v", input.Token)
	}
	err = uc.passwordResetRepository.Delete(pwReset)
	if err != nil {
		return nil, fmt.Errorf("delete password reset: %w", err)
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	err = uc.userRepository.UpdatePasswordHash(pwReset.UserID, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	// Create will create a new session for the user provided. The session token
	// will be returned as the Token field on the Session type, but only the hashed
	// session token is stored in the database.
	bytesPerToken := uc.bytesPerToken
	if bytesPerToken < tokenManager.MIN_BYTES_PER_TOKEN {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, tokenHash, err := uc.tokenManager.NewToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}
	session := entity.NewSession(uc.idGenerator.Generate(), pwReset.UserID, token, tokenHash)
	_, err = uc.sessionRepository.Upsert(session)
	if err != nil {
		return nil, fmt.Errorf("upsert session: %w", err)
	}
	return session, err
}

type ResetPasswordInput struct {
	Token    string
	Password string
}
