package usecases

import (
	"lenslocked/application/repository"
	"os"
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
	image, err := uc.imageRepository.FindOne(galleryID, filename)
	if err != nil {
		return err
	}
	err = os.Remove(image.Path)
	if err != nil {
		return err
	}
	return nil
}
