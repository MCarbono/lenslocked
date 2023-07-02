package app

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked/application/usecases"
	"lenslocked/idGenerator"
	"lenslocked/infra/controllers"
	"lenslocked/infra/database"
	"lenslocked/infra/database/migrations"
	"lenslocked/infra/gateway"
	"lenslocked/infra/http/router"
	"lenslocked/services"
	"lenslocked/templates"
	"lenslocked/tokenManager"
	"lenslocked/views"

	repository "lenslocked/infra/repository/postgres"
)

func Start() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")
	idGenerator := idGenerator.New()
	userRepository := repository.NewUserRepositoryPostgres(db)
	sessionRepository := repository.NewSessionRepositoryPostgres(db)
	tokenManager := tokenManager.New()
	sessionService := &services.SessionService{
		DB:                db,
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		TokenManager:      tokenManager,
		IDGenerator:       idGenerator,
	}
	pwResetService := &services.PasswordResetService{
		DB:             db,
		UserRepository: userRepository,
		PasswordReset:  repository.NewPasswordResetPostgres(db),
		TokenManager:   tokenManager,
		EmailGateway: gateway.NewEmailMailtrapGateway(gateway.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
		}),
		SessionRepository: sessionRepository,
		IDGenerator:       idGenerator,
	}
	usersC := controllers.Users{
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		CreateUserUseCase:    usecases.NewCreateUserUseCase(userRepository, idGenerator),
		CreateSessionUseCase: usecases.NewCreateSessionUseCase(sessionRepository, tokenManager, idGenerator),
		SignInUseCase:        usecases.NewSignInUseCase(sessionRepository, userRepository, tokenManager, idGenerator),
		SignOutUseCase:       usecases.NewSignOutUseCase(sessionRepository, tokenManager),
		Templates: struct {
			New            controllers.Template
			SignIn         controllers.Template
			ForgotPassword controllers.Template
			CheckYourEmail controllers.Template
			ResetPassword  controllers.Template
		}{
			New:            views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml")),
			SignIn:         views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml")),
			ForgotPassword: views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml")),
			CheckYourEmail: views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml")),
			ResetPassword:  views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml")),
		},
	}

	galleryRepository := repository.NewGalleryRepositoryPostgres(db)
	createGalleryUseCase := usecases.NewCreateGalleryUseCase(galleryRepository, idGenerator)
	updateGalleryUseCase := usecases.NewUpdateGalleryUseCase(galleryRepository)
	findGalleryUseCase := usecases.NewFindGalleryUseCase(galleryRepository)
	findGalleriesUseCase := usecases.NewFindGalleriesUseCase(galleryRepository)
	deleteGalleryUseCase := usecases.NewDeleteGalleryUseCase(galleryRepository)

	galleryController := controllers.Galleries{
		CreateGalleryUseCase: createGalleryUseCase,
		UpdateGalleryUseCase: updateGalleryUseCase,
		FindGalleryUseCase:   findGalleryUseCase,
		FindGalleriesUseCase: findGalleriesUseCase,
		DeleteGalleryUseCase: deleteGalleryUseCase,
		Templates: struct {
			Show  controllers.Template
			New   controllers.Template
			Edit  controllers.Template
			Index controllers.Template
		}{
			Show:  views.Must(views.ParseFS(templates.FS, "galleries/show.gohtml", "tailwind.gohtml")),
			New:   views.Must(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml")),
			Edit:  views.Must(views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml")),
			Index: views.Must(views.ParseFS(templates.FS, "galleries/index.gohtml", "tailwind.gohtml")),
		},
	}

	fmt.Printf("Starting the server on port %v\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router.New(usersC, galleryController, cfg.CSRF.Key, cfg.CSRF.Secure)))
}
