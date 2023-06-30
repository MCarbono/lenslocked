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
		return nil, err
	}
	return &gallery, nil
}

func (p *GalleryRepositoryPostgres) FindAll(UserID string) ([]*entity.Gallery, error) {
	rows, err := p.DB.Query("SELECT * from galleries where user_id = $1", UserID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user id %w", err)
	}
	var galleries []*entity.Gallery
	for rows.Next() {
		var gallery entity.Gallery
		err := rows.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id %w", err)
		}
		galleries = append(galleries, &gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user id %w", err)
	}
	return galleries, nil
}

func (p *GalleryRepositoryPostgres) Update(gallery *entity.Gallery) error {
	_, err := p.DB.Exec(`UPDATE galleries SET title = $1 WHERE id = $2`, gallery.Title, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery %w", err)
	}
	return nil
}

func (p *GalleryRepositoryPostgres) Delete(ID string) error {
	_, err := p.DB.Exec(`DELETE FROM galleries WHERE id = $1`, ID)
	if err != nil {
		return fmt.Errorf("delete gallery %w", err)
	}
	return nil
}
