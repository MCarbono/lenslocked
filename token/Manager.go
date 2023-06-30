package token

type Manager interface {
	NewToken(bytesPerToken int) (token, tokenHash string, err error)
	Hash(token string) string
}
