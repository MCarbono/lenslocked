package repository

import "lenslocked/domain/entity"

type ImageRepository interface {
	Find(galleryID string) ([]*entity.Image, error)
	FindOne(galleryID, filename string) (*entity.Image, error)
}
