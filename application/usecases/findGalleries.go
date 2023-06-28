package usecases

import (
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
)

type FindGalleriesUseCase struct {
	galleryRepository repository.GalleryRepository
}

func NewFindGalleriesUseCase(galleryRepository repository.GalleryRepository) *FindGalleriesUseCase {
	return &FindGalleriesUseCase{
		galleryRepository: galleryRepository,
	}
}

func (uc *FindGalleriesUseCase) Execute(UserID string) ([]*entity.Gallery, error) {
	galleries, err := uc.galleryRepository.FindAll(UserID)
	if err != nil {
		return nil, err
	}
	return galleries, nil
}
