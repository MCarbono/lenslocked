package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"lenslocked/domain/entity"
	"os"
	"path/filepath"
)

type ImageRepositoryDisk struct {
	path       string
	extensions []string
}

func NewImageRepositoryDisk(path string, extensions []string) *ImageRepositoryDisk {
	return &ImageRepositoryDisk{
		path:       path,
		extensions: extensions,
	}
}

func (r *ImageRepositoryDisk) Find(galleryID string) ([]*entity.Image, error) {
	return nil, nil
}

func (r *ImageRepositoryDisk) FindOne(galleryID, filename string) (*entity.Image, error) {
	imagePath := filepath.Join(r.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, fmt.Errorf("querying for image: %w", err)
	}
	return entity.NewImage(galleryID, imagePath, filename), nil
}

func (r *ImageRepositoryDisk) galleryDir(galleryID string) string {
	return filepath.Join(r.path, fmt.Sprintf("gallery-%s", galleryID))
}
