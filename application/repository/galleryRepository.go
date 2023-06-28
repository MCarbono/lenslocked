package repository

import "lenslocked/domain/entity"

type GalleryRepository interface {
	Create(gallery *entity.Gallery) error
}
