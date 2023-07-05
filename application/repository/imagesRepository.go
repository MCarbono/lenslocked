package repository

import "lenslocked/domain/entity"

type ImagesRepository interface {
	Find(galleryID string) ([]*entity.Image, error)
	FindOne(galleryID, filename string) (*entity.Image, error)
}
