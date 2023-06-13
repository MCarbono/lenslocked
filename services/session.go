package services

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
	"lenslocked/domain/repository"
	"lenslocked/token"
)

type SessionService struct {
	DB                *sql.DB
	SessionRepository repository.SessionRepository
	UserRepository    repository.UserRepository
	TokenManager      token.Manager
	BytesPerToken     int
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*entity.Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, tokenHash, err := ss.TokenManager.New(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := entity.Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}
	insertedSession, err := ss.SessionRepository.Upsert(&session)
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
