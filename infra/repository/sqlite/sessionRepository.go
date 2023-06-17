package repository

import (
	"database/sql"
	"fmt"
	"lenslocked/domain/entity"
)

type SessionRepositorySQLite struct {
	DB *sql.DB
}

func NewSessionRepositorySQLite(db *sql.DB) *SessionRepositorySQLite {
	return &SessionRepositorySQLite{
		DB: db,
	}
}

func (s *SessionRepositorySQLite) Upsert(session *entity.Session) (*entity.Session, error) {
	row, err := s.DB.Exec(`UPDATE sessions SET token_hash = ? WHERE user_id = ?`, session.TokenHash, session.UserID)
	if err != nil {
		return nil, err
	}
	updatedRow, err := row.RowsAffected()
	if err != nil {
		return nil, err
	}
	if updatedRow == 0 {
		//if no session exists, we will get ErrNoRows. That means we need to
		//create a session object for that user
		_, err = s.DB.Exec(`INSERT INTO sessions (id, user_id, token_hash) VALUES (?, ?, ?)`, session.ID, session.UserID, session.TokenHash)
	}
	//If the err was not sql.ErrNoRows, we need to check to see if it was any
	//other error. If it was sql.ErrNoRows it will be overwritten inside the if
	//block, and we still need to check for any errors.
	if err != nil {
		return nil, fmt.Errorf("upsert: %w", err)
	}
	return session, nil
}

func (s *SessionRepositorySQLite) FindByTokenHash(token string) (*entity.Session, error) {
	var session entity.Session
	row := s.DB.QueryRow(`SELECT * FROM sessions WHERE token_hash = ?`, token)
	err := row.Scan(&session.ID, &session.UserID, &session.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("session: %w", err)
	}
	return &session, nil
}

func (s *SessionRepositorySQLite) Delete(token string) error {
	if _, err := s.DB.Exec(`DELETE FROM sessions WHERE token_hash = $1`, token); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
