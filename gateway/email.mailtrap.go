package gateway

import (
	"fmt"
	"lenslocked/domain/entity"

	"github.com/go-mail/mail/v2"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailMailtrapGateway struct {
	dialer *mail.Dialer
}

func NewEmailMailtrapGateway(config SMTPConfig) *EmailMailtrapGateway {
	e := EmailMailtrapGateway{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &e
}

func (e *EmailMailtrapGateway) Send(email *entity.Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	msg.SetHeader("From", email.From)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
	msg.AddAlternative("text/html", email.HTML)
	if err := e.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}
