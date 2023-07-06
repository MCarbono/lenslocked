package usecases

import (
	"lenslocked/application/repository"
	"lenslocked/domain/entity"
	"path/filepath"
)

type FindImageUseCase struct {
	imageRepository repository.ImageRepository
}

func NewFindImageUseCase(imageRepository repository.ImageRepository) *FindImageUseCase {
	return &FindImageUseCase{
		imageRepository: imageRepository,
	}
}

func (uc *FindImageUseCase) Execute(galleryID, filename string) (*entity.Image, error) {
	return uc.imageRepository.FindOne(galleryID, filepath.Base(filename))
}
