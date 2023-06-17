package repository

import "lenslocked/domain/entity"

type PasswordResetRepository interface {
	Create(passwordReset *entity.PasswordReset) error
	FindByID(id string) (*entity.PasswordReset, error)
	FindByTokenHash(tokenHash string) (*entity.PasswordReset, error)
	Delete(passwordReset *entity.PasswordReset) error
}
