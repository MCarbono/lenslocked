package usecases

import (
	"lenslocked/application/repository"
	"path/filepath"
)

type DeleteImageUseCase struct {
	imageRepository repository.ImageRepository
}

func NewDeleteImageUseCase(imageRepository repository.ImageRepository) *DeleteImageUseCase {
	return &DeleteImageUseCase{
		imageRepository: imageRepository,
	}
}

func (uc *DeleteImageUseCase) Execute(galleryID, filename string) error {
	return uc.imageRepository.DeleteOne(galleryID, filepath.Base(filename))
}
