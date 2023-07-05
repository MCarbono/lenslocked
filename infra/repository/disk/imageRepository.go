package repository

import "lenslocked/domain/entity"

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
	return nil, nil
}
