package entity

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

func NewUser(ID, email, passwordHash string) *User {
	return &User{
		ID:           ID,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
