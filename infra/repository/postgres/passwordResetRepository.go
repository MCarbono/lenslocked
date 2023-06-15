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

func (p *PasswordResetPostgres) FindByID(id int) (*entity.PasswordReset, error) {
	var passwordResets entity.PasswordReset
	row := p.DB.QueryRow(`SELECT * FROM password_resets WHERE id = $1`, id)
	if err := row.Scan(&passwordResets.ID, &passwordResets.UserID, &passwordResets.TokenHash, &passwordResets.ExpiresAt); err != nil {
		return nil, fmt.Errorf("find by id password_resets: %w", err)
	}
	return &passwordResets, nil
}

func (p *PasswordResetPostgres) FindByTokenHash(tokenHash string) (*entity.PasswordReset, error) {
	var passwordResets entity.PasswordReset
	row := p.DB.QueryRow(`SELECT * FROM password_resets WHERE token_hash = $1`, tokenHash)
	if err := row.Scan(&passwordResets.ID, &passwordResets.UserID, &passwordResets.TokenHash, &passwordResets.ExpiresAt); err != nil {
		return nil, fmt.Errorf("find by token hash password_resets: %w", err)
	}
	return &passwordResets, nil
}

func (p *PasswordResetPostgres) Delete(passwordReset *entity.PasswordReset) error {
	_, err := p.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;`, passwordReset.ID)
	return err
}
