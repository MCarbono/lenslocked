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
		return nil, fmt.Errorf("gallery: %w", err)
	}
	return &gallery, nil
}

func (p *GalleryRepositorySQLite) FindAll(UserID string) ([]*entity.Gallery, error) {
	rows, err := p.DB.Query("SELECT * from galleries where user_id = ?", UserID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user id %w", err)
	}
	var galleries []*entity.Gallery
	for rows.Next() {
		var gallery *entity.Gallery
		err := rows.Scan(gallery.ID, gallery.UserID, gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user id %w", err)
	}
	return galleries, nil
}

func (p *GalleryRepositorySQLite) Update(gallery *entity.Gallery) error {
	_, err := p.DB.Exec(`UPDATE galleries SET title = ? WHERE id = ?`, gallery.Title, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery %w", err)
	}
	return nil
}

func (p *GalleryRepositorySQLite) Delete(ID string) error {
	_, err := p.DB.Exec(`DELETE FROM galleries WHERE id = ?`, ID)
	if err != nil {
		return fmt.Errorf("delete gallery %w", err)
	}
	return nil
}
