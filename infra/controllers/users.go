package controllers

import (
	"fmt"
	"lenslocked/application/usecases"
	"lenslocked/context"
	"lenslocked/errors"
	"lenslocked/infra/http/cookie"
	repository "lenslocked/infra/repository/postgres"
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
	PasswordResetService   *services.PasswordResetService
	CreateUserUseCase      *usecases.CreateUserUseCase
	CreateSessionUseCase   *usecases.CreateSessionUseCase
	SignInUseCase          *usecases.SignInUseCase
	SignOutUseCase         *usecases.SignOutUseCase
	FindUserByTokenUseCase *usecases.FindUserByTokenUseCase
	ForgotPasswordUseCase  *usecases.ForgotPasswordUseCase
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
	input := &usecases.CreateUserInput{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	user, err := u.CreateUserUseCase.Execute(input)

	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account.")
		}
		u.Templates.New.Execute(w, r, input, err)
		return
	}

	session, err := u.CreateSessionUseCase.Execute(user.ID)
	if err != nil {
		fmt.Println(err)
		//TODO: long term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie.SetCookie(w, cookie.CookieSession, session.Token)
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
	signInInput := &usecases.SignInInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	session, err := u.SignInUseCase.Execute(signInInput)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	cookie.SetCookie(w, cookie.CookieSession, session.Token)
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
	token, err := cookie.ReadCookie(r, cookie.CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if err = u.SignOutUseCase.Execute(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	cookie.DeleteCookie(w, cookie.CookieSession)
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
	input := &usecases.ForgotPasswordInput{
		Email:            r.FormValue("email"),
		ResetPasswordURL: "http://localhost:3000/reset-pw?",
	}
	_, err := u.ForgotPasswordUseCase.Execute(input)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, input)
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

	session, err := u.PasswordResetService.Consume(data.Token, data.Password)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between server errors and invalid token errors.
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	//Sign the user in now that they have reset their password.
	//Any errors from this point onward should redirect to the sign in page.
	// session, err := u.SessionService.Create(user.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	cookie.SetCookie(w, cookie.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
