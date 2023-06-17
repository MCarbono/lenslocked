package entity

type Session struct {
	ID     string
	UserID string
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

func NewSession(ID, userID, token, tokenHash string) *Session {
	return &Session{
		ID:        ID,
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}
}
