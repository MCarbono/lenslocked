package repository

import "lenslocked/domain/entity"

type SessionRepository interface {
	Upsert(session *entity.Session) (*entity.Session, error)
	FindByTokenHash(token string) (*entity.Session, error)
	Delete(token string) error
}
