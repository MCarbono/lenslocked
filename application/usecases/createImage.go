package usecases

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CreateImageUsecase struct {
}

func NewCreateImageUsecase() *CreateImageUsecase {
	return &CreateImageUsecase{}
}

func (uc *CreateImageUsecase) Execute(input *CreateImageInput) error {
	galleryDir := filepath.Join("images", fmt.Sprintf("gallery-%s", input.GalleryID))
	err := os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%s images directory: %w", galleryDir, err)
	}
	imagePath := filepath.Join(galleryDir, input.Filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, input.Contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

type CreateImageInput struct {
	GalleryID string
	Filename  string
	Contents  io.Reader
}
