package controllers

import (
	"fmt"
	"lenslocked/context"
	"lenslocked/services"
	"net/http"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *services.UserService
	SessionService       *services.SessionService
	PasswordResetService *services.PasswordResetService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
		return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	user, err := u.UserService.Create(email, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		//TODO: long term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
	// if user == nil {
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	// token, err := readCookie(r, CookieSession)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusSeeOther)
	// 	return
	// }
	// user, err := u.SessionService.User(token)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }

}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if err = u.SessionService.Delete(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	_, err := u.PasswordResetService.Create(data.Email, "http://localhost:3000/reset-pw?")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token, data.Password)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between server errors and invalid token errors.
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	//Sign the user in now that they have reset their passord.
	//Any errors from this point onward should redirect to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
