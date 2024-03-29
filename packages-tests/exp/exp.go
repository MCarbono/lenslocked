package main

import (
	// "html/template"
	"os"

	"github.com/go-mail/mail/v2"
)

type User struct {
	AString  string
	AInteger int
}

func main() {
	// t, err := template.ParseFiles("hello.gohtml")
	// if err != nil {
	// 	panic(err)
	// }
	// user := User{
	// 	AString:  "Name",
	// 	AInteger: 25,
	// }
	// err = t.Execute(os.Stdout, user)
	// if err != nil {
	// 	panic(err)
	// }
	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	msg.WriteTo(os.Stdout)
}
