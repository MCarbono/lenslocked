package usecases

import (
	"errors"
	"fmt"
	"lenslocked/application/repository"
)

type UpdateGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
	userRepository    repository.UserRepository
}

func NewUpdateGalleryUseCase(galleryRepository repository.GalleryRepository, userRepository repository.UserRepository) *UpdateGalleryUseCase {
	return &UpdateGalleryUseCase{
		galleryRepository: galleryRepository,
		userRepository:    userRepository,
	}
}

func (uc *UpdateGalleryUseCase) Execute(input *UpdateGalleryInput) error {
	gallery, err := uc.galleryRepository.FindByID(input.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	user, err := uc.userRepository.FindByID(gallery.UserID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	if !gallery.IsOwnedBy(user.ID) {
		return errors.New("you are not authorized to edit this gallery")
	}
	gallery.Update(input.Title)
	if err := uc.galleryRepository.Update(gallery); err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

type UpdateGalleryInput struct {
	ID    string
	Title string
}
