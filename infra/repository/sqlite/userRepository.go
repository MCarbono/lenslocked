package repository

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
)

type UserRepositorySQLite struct {
	DB *sql.DB
}

func NewUserRepositorySQLite(db *sql.DB) *UserRepositorySQLite {
	return &UserRepositorySQLite{
		DB: db,
	}
}

func (p *UserRepositorySQLite) Create(user *entity.User) error {
	_, err := p.DB.Exec(`INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)`, user.ID, user.Email, user.PasswordHash)
	return err
}

func (p *UserRepositorySQLite) FindByEmail(email string) (*entity.User, error) {
	row := p.DB.QueryRow("SELECT * from users WHERE email = ?", email)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *UserRepositorySQLite) FindByID(ID string) (*entity.User, error) {
	var user entity.User
	row := p.DB.QueryRow(`SELECT * FROM users WHERE id = ?`, ID)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (p *UserRepositorySQLite) FindByTokenHash(token string) (*entity.User, error) {
	query := `
		SELECT user.id, user.email, user.password_hash from users AS user
		JOIN sessions AS session ON user.id = session.user_id
		WHERE session.token_hash = ?
	`
	row := p.DB.QueryRow(query, token)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *UserRepositorySQLite) UpdatePasswordHash(id string, passwordHash string) error {
	_, err := p.DB.Exec(`
		UPDATE users
		SET password_hash = ?
		WHERE id = ?`, passwordHash, id)
	return err
}
