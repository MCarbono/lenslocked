package repository

import (
	"database/sql"
	"fmt"
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

func (p *PasswordResetSQLite) Create(passwordReset *entity.PasswordReset) error {
	_, err := p.DB.Exec(`INSERT INTO password_resets (id, user_id, token_hash, expires_at)
		values (?, ?, ?, ?) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = excluded.token_hash, expires_at = excluded.expires_at`, passwordReset.ID, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	return err
}

func (p *PasswordResetSQLite) FindByID(id string) (*entity.PasswordReset, error) {
	var passwordResets entity.PasswordReset
	row := p.DB.QueryRow(`SELECT * FROM password_resets WHERE id = ?`, id)
	if err := row.Scan(&passwordResets.ID, &passwordResets.UserID, &passwordResets.TokenHash, &passwordResets.ExpiresAt); err != nil {
		return nil, fmt.Errorf("password_resets: %w", err)
	}
	return &passwordResets, nil
}

func (p *PasswordResetSQLite) FindByTokenHash(tokenHash string) (*entity.PasswordReset, error) {
	var passwordResets entity.PasswordReset
	row := p.DB.QueryRow(`SELECT * FROM password_resets WHERE token_hash = ?`, tokenHash)
	if err := row.Scan(&passwordResets.ID, &passwordResets.UserID, &passwordResets.TokenHash, &passwordResets.ExpiresAt); err != nil {
		return nil, fmt.Errorf("password_resets: %w", err)
	}
	return &passwordResets, nil
}

func (p *PasswordResetSQLite) Delete(passwordReset *entity.PasswordReset) error {
	_, err := p.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = ?;`, passwordReset.ID)
	return err
}
