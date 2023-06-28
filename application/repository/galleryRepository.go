package repository

import "lenslocked/domain/entity"

type GalleryRepository interface {
	Create(gallery *entity.Gallery) error
	FindByID(ID string) (*entity.Gallery, error)
	FindAll(UserID string) ([]*entity.Gallery, error)
	Update(gallery *entity.Gallery) error
	Delete(ID string) error
}
