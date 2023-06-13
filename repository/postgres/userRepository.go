package repository

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
)

type UserRepositoryPostgres struct {
	DB *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		DB: db,
	}
}

func (p *UserRepositoryPostgres) Create(email, password string) (int, error) {
	row := p.DB.QueryRow(`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`, email, password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (p *UserRepositoryPostgres) FindByEmail(email string) (*entity.User, error) {
	row := p.DB.QueryRow("SELECT * from users WHERE email = $1", email)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *UserRepositoryPostgres) FindByID(ID int) (*entity.User, error) {
	var user entity.User
	row := p.DB.QueryRow(`SELECT * FROM users WHERE id = $1`, ID)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (p *UserRepositoryPostgres) FindByTokenHash(token string) (*entity.User, error) {
	var user entity.User
	row := p.DB.QueryRow(
		`SELECT users.id, users.email, users.password_hash FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1;
		`, token)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (p *UserRepositoryPostgres) UpdatePasswordHash(id int, passwordHash string) error {
	_, err := p.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1`, id, passwordHash)
	return err
}
