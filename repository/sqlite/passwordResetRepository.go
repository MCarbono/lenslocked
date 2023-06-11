package repository

import (
	"database/sql"
	"lenslocked/domain/entity"
)

type PasswordResetSQLite struct {
	DB *sql.DB
}

func NewPasswordResetSQLite(db *sql.DB) *PasswordResetSQLite {
	return &PasswordResetSQLite{
		DB: db,
	}
}

func (p *PasswordResetSQLite) Create(passwordReset *entity.PasswordReset) (int, error) {
	row, err := p.DB.Exec(`INSERT INTO password_resets (user_id, token_hash, expires_at)
		values (?, ?, ?) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = excluded.token_hash, expires_at = excluded.expires_at`, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	if err != nil {
		return 0, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}
