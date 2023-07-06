package usecases

import (
	"lenslocked/application/repository"
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
	return uc.imageRepository.DeleteOne(galleryID, filename)
}
