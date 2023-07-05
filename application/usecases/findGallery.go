package usecases

import (
	"fmt"
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
)

type FindGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
	imageRepository   repository.ImageRepository
}

func NewFindGalleryUseCase(galleryRepository repository.GalleryRepository, imageRepository repository.ImageRepository) *FindGalleryUseCase {
	return &FindGalleryUseCase{
		galleryRepository: galleryRepository,
		imageRepository:   imageRepository,
	}
}

func (uc *FindGalleryUseCase) Execute(ID string) (*entity.Gallery, error) {
	gallery, err := uc.galleryRepository.FindByID(ID)
	if err != nil {
		return nil, fmt.Errorf("query gallery by ID %w", err)
	}
	images, err := uc.imageRepository.Find(gallery.ID)
	if err != nil {
		return nil, err
	}
	gallery.AddImages(images)
	return gallery, nil
}
