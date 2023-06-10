package token

type Manager interface {
	New(bytesPerToken int) (token, tokenHash string, err error)
	Hash(token string) string
}
