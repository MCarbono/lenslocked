package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"lenslocked/idGenerator"
)

type CreateGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
	idGenerator       idGenerator.IDGenerator
}

func NewCreateGalleryUseCase(galleryRepository repository.GalleryRepository, idGenerator idGenerator.IDGenerator) *CreateGalleryUseCase {
	return &CreateGalleryUseCase{
		galleryRepository: galleryRepository,
		idGenerator:       idGenerator,
	}
}

func (uc *CreateGalleryUseCase) Execute(input *CreateGalleryInput) (*entity.Gallery, error) {
	gallery := entity.NewGallery(uc.idGenerator.Generate(), input.UserID, input.Title)
	if err := uc.galleryRepository.Create(gallery); err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return gallery, nil
}

type CreateGalleryInput struct {
	Title  string
	UserID string
}
