package testinfra

import (
	"database/sql"
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/middleware"
	"lenslocked/templates"
	"lenslocked/views"
	"net/http"

	"github.com/go-chi/chi"
)

func CreateDatabaseTest() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../lenslocked_test.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS sessions;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS password_resets;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS galleries;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE users (id TEXT PRIMARY KEY, email TEXT UNIQUE NOT NULL,password_hash TEXT NOT NULL);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE sessions (id TEXT PRIMARY KEY, user_id INT UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE, token_hash TEXT UNIQUE NOT NULL);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE password_resets (id TEXT PRIMARY KEY, user_id INT UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE, token_hash TEXT UNIQUE NOT NULL, expires_at TIMESTAMP NOT NULL);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE galleries (id TEXT PRIMARY KEY, user_id TEXT REFERENCES users (id) ON DELETE CASCADE, title TEXT);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewRouterTest(usersC controllers.Users, galleryController controllers.Galleries) http.Handler {
	umw := middleware.UserMiddleware{
		FindUserByTokenUseCase: usersC.FindUserByTokenUseCase,
	}
	r := chi.NewRouter()
	r.Use(middleware.HTMLResponse)
	r.Use(umw.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	galleryController.Templates.New = views.Must(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))
	galleryController.Templates.Edit = views.Must(views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"))

	r.Get("/users/new", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	//Gallery
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleryController.Show)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleryController.Index)
			r.Get("/new", galleryController.New)
			r.Post("/", galleryController.Create)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}", galleryController.Update)
			r.Post("/{id}/delete", galleryController.Delete)
		})
	})
	return r
}
