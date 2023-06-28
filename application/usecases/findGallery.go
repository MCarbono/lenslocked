package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
)

type FindGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
}

func NewFindGalleryUseCase(galleryRepository repository.GalleryRepository) *FindGalleryUseCase {
	return &FindGalleryUseCase{
		galleryRepository: galleryRepository,
	}
}

func (uc *FindGalleryUseCase) Execute(ID string) (*entity.Gallery, error) {
	gallery, err := uc.galleryRepository.FindByID(ID)
	if err != nil {
		return nil, fmt.Errorf("query gallery by ID %w", err)
	}
	return gallery, nil
}
