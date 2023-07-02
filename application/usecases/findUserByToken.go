package usecases

import (
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/tokenManager"
)

type FindUserByTokenUseCase struct {
	userRepository repository.UserRepository
	tokenManager   tokenManager.Manager
}

func NewFindUserByTokenUseCase(userRepository repository.UserRepository, tokenManager tokenManager.Manager) *FindUserByTokenUseCase {
	return &FindUserByTokenUseCase{
		userRepository: userRepository,
		tokenManager:   tokenManager,
	}
}

func (uc *FindUserByTokenUseCase) Execute(token string) (*entity.User, error) {
	tokenHash := uc.tokenManager.Hash(token)
	user, err := uc.userRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
