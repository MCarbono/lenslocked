package entity

import "time"

type PasswordReset struct {
	ID     int
	UserID int
	//Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

func NewPasswordReset(userID int, token, tokenHash string, duration time.Duration) *PasswordReset {
	return &PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}
}
