package token

import (
	"crypto/sha256"
	"encoding/base64"
	"lenslocked/rand"
)

type ManagerImpl struct {
	minBytesPerToken int
}

func (tm ManagerImpl) NewToken(bytesPerToken int) (token, tokenHash string, err error) {
	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", err
	}
	return token, tm.Hash(token), nil
}

func (tm ManagerImpl) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func New() ManagerImpl {
	return ManagerImpl{
		minBytesPerToken: 32,
	}
}
