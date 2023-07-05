package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"lenslocked/domain/entity"
	"net/url"
	"os"
	"path/filepath"
)

type ImageRepositoryDisk struct {
	path string
}

func NewImageRepositoryDisk(path string) *ImageRepositoryDisk {
	return &ImageRepositoryDisk{
		path: path,
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
	return &entity.Image{
		Filename:        filename,
		GalleryID:       galleryID,
		Path:            imagePath,
		FilenameEscaped: url.PathEscape(filepath.Base(filename)),
	}, nil
}

func (r *ImageRepositoryDisk) galleryDir(galleryID string) string {
	return filepath.Join(r.path, fmt.Sprintf("gallery-%s", galleryID))
}
