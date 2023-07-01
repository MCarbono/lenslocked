package router

import (
	"lenslocked/infra/controllers"
	"lenslocked/infra/http/middleware"
	"lenslocked/templates"
	"lenslocked/views"
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

func New(usersC controllers.Users, galleryController controllers.Galleries, csrfKey string, csrfSecure bool) http.Handler {
	umw := middleware.UserMiddleware{
		SessionService: usersC.SessionService,
	}

	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(csrfSecure), csrf.Path("/"))
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(middleware.HTMLResponse)
	r.Use(csrfMw)
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
