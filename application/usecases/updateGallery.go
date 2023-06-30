package usecases

import (
	"fmt"
	"lenslocked/application/repository"
)

type UpdateGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
}

func NewUpdateGalleryUseCase(galleryRepository repository.GalleryRepository) *UpdateGalleryUseCase {
	return &UpdateGalleryUseCase{
		galleryRepository: galleryRepository,
	}
}

func (uc *UpdateGalleryUseCase) Execute(input *UpdateGalleryInput) error {
	gallery, err := uc.galleryRepository.FindByID(input.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
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
