package services

import (
	"database/sql"
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

func (ss *SessionService) User(token string) (*entity.User, error) {
	tokenHash := ss.TokenManager.Hash(token)
	user, err := ss.UserRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
