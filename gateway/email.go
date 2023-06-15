package gateway

import "lenslocked/domain/entity"

//go:generate mockgen -destination=../gen/mock/gateway_email_mock.go -package=mock . EmailProvider
type EmailProvider interface {
	Send(email *entity.Email) error
}
