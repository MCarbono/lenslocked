package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepository
	idGenerator    idGenerator.IDGenerator
}

func NewCreateUserUseCase(
	userRepository repository.UserRepository,
	idGenerator idGenerator.IDGenerator) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		idGenerator:    idGenerator,
	}
}

func (uc *CreateUserUseCase) Execute(input *CreateUserInput) (*entity.User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)
	user := entity.NewUser(uc.idGenerator.Generate(), input.Email, passwordHash)
	err = uc.userRepository.Create(user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

type CreateUserInput struct {
	Email    string
	Password string
}
