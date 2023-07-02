package usecases

import (
	"lenslocked/application/repository"
	"lenslocked/tokenManager"
)

type SignOutUseCase struct {
	sessionRepository repository.SessionRepository
	tokenManager      tokenManager.Manager
}

func NewSignOutUseCase(sessionRepository repository.SessionRepository, tokenManager tokenManager.Manager) *SignOutUseCase {
	return &SignOutUseCase{
		sessionRepository: sessionRepository,
		tokenManager:      tokenManager,
	}
}

func (uc *SignOutUseCase) Execute(token string) error {
	return uc.sessionRepository.Delete(uc.tokenManager.Hash(token))
}
