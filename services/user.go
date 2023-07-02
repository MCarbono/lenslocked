package services

import (
	"database/sql"
	"lenslocked/application/gateway"
	"lenslocked/application/repository"
	"lenslocked/idGenerator"
)

const (
	//DefaultSender is the default email address to send emails from.
	DefaultSender = "support@lenslocked.com"
)

type UserService struct {
	DB             *sql.DB
	UserRepository repository.UserRepository
	EmailGateway   gateway.EmailProvider
	idGenerator.IDGenerator
}
