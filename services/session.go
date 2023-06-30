package services

import (
	"database/sql"
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"
)

type SessionService struct {
	DB                *sql.DB
	SessionRepository repository.SessionRepository
	UserRepository    repository.UserRepository
	TokenManager      tokenManager.Manager
	BytesPerToken     int
	idGenerator.IDGenerator
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID string) (*entity.Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < tokenManager.MIN_BYTES_PER_TOKEN {
		bytesPerToken = tokenManager.MIN_BYTES_PER_TOKEN
	}
	token, tokenHash, err := ss.TokenManager.NewToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := entity.NewSession(ss.Generate(), userID, token, tokenHash)
	insertedSession, err := ss.SessionRepository.Upsert(session)
	if err != nil {
		return nil, err
	}
	return insertedSession, err
}

func (ss *SessionService) User(token string) (*entity.User, error) {
	tokenHash := ss.TokenManager.Hash(token)
	user, err := ss.UserRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ss *SessionService) Delete(token string) error {
	return ss.SessionRepository.Delete(ss.TokenManager.Hash(token))
}
