package repository

import (
	"database/sql"
	"fmt"
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
	_, err := p.DB.Exec(`INSERT INTO galleries (id, user_id, title) VALUES ($1, $2, $3) `, gallery.ID, gallery.UserID, gallery.Title)
	return err
}

func (p *GalleryRepositoryPostgres) FindByID(ID string) (*entity.Gallery, error) {
	var gallery entity.Gallery
	row := p.DB.QueryRow(`SELECT * FROM galleries WHERE id = $1`, ID)
	if err := row.Scan(&gallery.ID, &gallery.UserID, &gallery.Title); err != nil {
		return nil, fmt.Errorf("gallery: %w", err)
	}
	return &gallery, nil
}
