package repository

import (
	"database/sql"
	"lenslocked/domain/entity"
)

type GalleryRepositorySQLite struct {
	DB *sql.DB
}

func NewGalleryRepositorySQLite(db *sql.DB) *GalleryRepositorySQLite {
	return &GalleryRepositorySQLite{
		DB: db,
	}
}

func (p *GalleryRepositorySQLite) Create(gallery *entity.Gallery) error {
	_, err := p.DB.Exec(`INSERT INTO users (id, user_id, title) VALUES (?, ?, ?)`, gallery.ID, gallery.UserID, gallery.Title)
	return err
}
