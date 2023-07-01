package services

import (
	"database/sql"
	"fmt"
	"lenslocked/application/gateway"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"

	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	//DefaultSender is the default email address to send emails from.
	DefaultSender = "support@lenslocked.com"
)

type UserService struct {
	DB             *sql.DB
	UserRepository repository.UserRepository
	EmailGateway   gateway.EmailProvider
	idGenerator.IDGenerator
}

type CreateUserInput struct {
	Email    string
	Password string
}

func (us *UserService) Authenticate(email, password string) (*entity.User, error) {
	email = strings.ToLower(email)
	user, err := us.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return user, nil
}

//TODO: remove this method
func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
