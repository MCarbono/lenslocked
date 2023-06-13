package repository

import "lenslocked/domain/entity"

type UserRepository interface {
	Create(email, password string) (int, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(ID int) (*entity.User, error)
	FindByTokenHash(token string) (*entity.User, error)
	UpdatePasswordHash(id int, passwordHash string) error
}
