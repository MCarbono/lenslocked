package gateway

import "lenslocked/domain/entity"

type EmailProvider interface {
	Send(email *entity.Email) error
}
