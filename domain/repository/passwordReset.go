package repository

import "lenslocked/domain/entity"

type PasswordResetRepository interface {
	Create(passwordReset *entity.PasswordReset) (int, error)
	FindByID(id int) (*entity.PasswordReset, error)
}
