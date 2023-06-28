package repository

import (
	"database/sql"
	"fmt"
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
	_, err := p.DB.Exec(`INSERT INTO galleries (id, user_id, title) VALUES (?, ?, ?)`, gallery.ID, gallery.UserID, gallery.Title)
	return err
}

func (p *GalleryRepositorySQLite) FindByID(ID string) (*entity.Gallery, error) {
	var gallery entity.Gallery
	row := p.DB.QueryRow(`SELECT * FROM galleries WHERE id = ?`, ID)
	if err := row.Scan(&gallery.ID, &gallery.UserID, &gallery.Title); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &gallery, nil
}
