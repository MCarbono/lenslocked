package repository

import "lenslocked/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(ID string) (*entity.User, error)
	FindByTokenHash(token string) (*entity.User, error)
	UpdatePasswordHash(id string, passwordHash string) error
}
