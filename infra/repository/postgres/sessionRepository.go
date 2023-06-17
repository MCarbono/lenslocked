package repository

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
)

type SessionRepositoryPostgres struct {
	DB *sql.DB
}

func NewSessionRepositoryPostgres(db *sql.DB) *SessionRepositoryPostgres {
	return &SessionRepositoryPostgres{
		DB: db,
	}
}

func (s *SessionRepositoryPostgres) Upsert(session *entity.Session) (*entity.Session, error) {
	_, err := s.DB.Exec(`
		INSERT INTO sessions (id, user_id, token_hash)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) 
		DO UPDATE SET token_hash = $3;`, session.ID, session.UserID, session.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return session, nil
}

func (s *SessionRepositoryPostgres) FindByTokenHash(token string) (*entity.Session, error) {
	var session entity.Session
	row := s.DB.QueryRow(`SELECT * FROM sessions WHERE token_hash = $1`, token)
	err := row.Scan(&session.ID, &session.UserID, &session.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("session: %w", err)
	}
	return &session, nil
}

func (s *SessionRepositoryPostgres) Delete(token string) error {
	if _, err := s.DB.Exec(`DELETE FROM sessions WHERE token_hash = $1`, token); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
