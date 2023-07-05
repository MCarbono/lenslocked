package usecases

import (
	"errors"
	"fmt"
	"io/fs"
	"lenslocked/domain/entity"
	"os"
	"path/filepath"
)

type FindImageUseCase struct {
	imagesDir string
}

func NewFindImageUseCase() *FindImageUseCase {
	return &FindImageUseCase{}
}

func (uc *FindImageUseCase) Execute(galleryID, filename string) (*entity.Image, error) {
	imagePath := filepath.Join(uc.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, fmt.Errorf("querying  for image: %w", err)
	}
	return &entity.Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (uc *FindImageUseCase) galleryDir(ID string) string {
	imagesDir := uc.imagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%s", ID))
}
