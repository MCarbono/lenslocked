package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"
)

type CreateSessionUseCase struct {
	tokenManager      tokenManager.Manager
	sessionRepository repository.SessionRepository
	idGenerator       idGenerator.IDGenerator
	bytesPerToken     int
}

func NewCreateSessionUseCase(sessionRepository repository.SessionRepository, tokenManager tokenManager.Manager, idGenerator idGenerator.IDGenerator) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepository: sessionRepository,
		tokenManager:      tokenManager,
		idGenerator:       idGenerator,
	}
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (uc *CreateSessionUseCase) Execute(userID string) (*entity.Session, error) {
	bytesPerToken := uc.bytesPerToken
	if bytesPerToken < tokenManager.MIN_BYTES_PER_TOKEN {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, tokenHash, err := uc.tokenManager.NewToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := entity.NewSession(uc.idGenerator.Generate(), userID, token, tokenHash)
	insertedSession, err := uc.sessionRepository.Upsert(session)
	if err != nil {
		return nil, err
	}
	return insertedSession, nil
}
