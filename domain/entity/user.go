package entity

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

func NewUser(ID int, email, passwordHash string) *User {
	return &User{
		ID:           ID,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
