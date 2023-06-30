package usecases

import (
	"lenslocked/application/gateway"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
	"lenslocked/tokenManager"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepository
	emailGateway   gateway.EmailProvider
	idGenerator    idGenerator.IDGenerator
	TokenManager   tokenManager.Manager
}

func NewCreateUserUseCase(
	userRepository repository.UserRepository,
	emailGateway gateway.EmailProvider,
	idGenerator idGenerator.IDGenerator,
	tokenManager tokenManager.Manager) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		emailGateway:   emailGateway,
		idGenerator:    idGenerator,
		TokenManager:   tokenManager,
	}
}

func (uc *CreateUserUseCase) Execute(input *CreateGalleryInput) (*entity.User, error) {
	return nil, nil
}

type CreateUserInput struct {
	Email    string
	Password string
}
