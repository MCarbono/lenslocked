package entity

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

//criar funcao construtora
func NewEmail(from, to, subject, plaintext, html string) *Email {
	return &Email{
		From:      from,
		To:        to,
		Subject:   subject,
		Plaintext: plaintext,
		HTML:      html,
	}
}
