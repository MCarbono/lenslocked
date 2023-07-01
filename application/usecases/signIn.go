package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type SignInUseCase struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	tokenManager      tokenManager.Manager
	idGenerator       idGenerator.IDGenerator
	bytesPerToken     int
}

func NewSignInUseCase(
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
	tokenManager tokenManager.Manager,
	idGenerator idGenerator.IDGenerator,
) *SignInUseCase {
	return &SignInUseCase{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		tokenManager:      tokenManager,
		idGenerator:       idGenerator,
	}
}

// Execute will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (uc *SignInUseCase) Execute(input *SignInInput) (*entity.Session, error) {
	input.Email = strings.ToLower(input.Email)
	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	bytesPerToken := uc.bytesPerToken
	if bytesPerToken < tokenManager.MIN_BYTES_PER_TOKEN {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, tokenHash, err := uc.tokenManager.NewToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := entity.NewSession(uc.idGenerator.Generate(), user.ID, token, tokenHash)
	_, err = uc.sessionRepository.Upsert(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

type SignInInput struct {
	Email    string
	Password string
}
