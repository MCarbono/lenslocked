package repository

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
)

type PasswordResetPostgres struct {
	DB *sql.DB
}

func NewPasswordResetPostgres(db *sql.DB) *PasswordResetPostgres {
	return &PasswordResetPostgres{
		DB: db,
	}
}

func (p *PasswordResetPostgres) Create(passwordReset *entity.PasswordReset) (int, error) {
	row := p.DB.QueryRow(`INSERT INTO password_resets (user_id, token_hash, expires_at)
		values ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}
