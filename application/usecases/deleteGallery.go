package usecases

import (
	"lenslocked/application/repository"
)

type DeleteGalleryUseCase struct {
	galleryRepository repository.GalleryRepository
}

func NewDeleteGalleryUseCase(galleryRepository repository.GalleryRepository) *DeleteGalleryUseCase {
	return &DeleteGalleryUseCase{
		galleryRepository: galleryRepository,
	}
}

func (uc *DeleteGalleryUseCase) Execute(ID string) error {
	return uc.galleryRepository.Delete(ID)
}
