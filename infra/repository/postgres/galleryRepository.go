package repository

import (
	"database/sql"
	"lenslocked/domain/entity"
)

type GalleryRepositoryPostgres struct {
	DB *sql.DB
}

func NewGalleryRepositoryPostgres(db *sql.DB) *GalleryRepositoryPostgres {
	return &GalleryRepositoryPostgres{
		DB: db,
	}
}

func (p *GalleryRepositoryPostgres) Create(gallery *entity.Gallery) error {
	_, err := p.DB.Exec(`INSERT INTO users (id, user_id, title) VALUES ($1, $2, $3) `, gallery.ID, gallery.UserID, gallery.Title)
	return err
}
