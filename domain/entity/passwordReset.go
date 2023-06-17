package entity

import "time"

type PasswordReset struct {
	ID     string
	UserID string
	//Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

func NewPasswordReset(ID, userID, token, tokenHash string, duration time.Duration) *PasswordReset {
	return &PasswordReset{
		ID:        ID,
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (pw *PasswordReset) IsExpired() bool {
	return time.Now().After(pw.ExpiresAt)
}
