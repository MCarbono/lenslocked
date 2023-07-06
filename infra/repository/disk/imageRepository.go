package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"lenslocked/domain/entity"
	"os"
	"path/filepath"
	"strings"
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
	globPattern := filepath.Join(r.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []*entity.Image
	for _, file := range allFiles {
		if r.hasExtension(file) {
			images = append(images, entity.NewImage(galleryID, file, filepath.Base(file)))
		}
	}
	return images, nil
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

func (r *ImageRepositoryDisk) DeleteOne(galleryID, filename string) error {
	imagePath := filepath.Join(r.galleryDir(galleryID), filename)
	err := os.Remove(imagePath)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryDisk) galleryDir(galleryID string) string {
	return filepath.Join(r.path, fmt.Sprintf("gallery-%s", galleryID))
}

func (r *ImageRepositoryDisk) hasExtension(file string) bool {
	for _, ext := range r.extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
