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
	"lenslocked/token"

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
	userService := &services.UserService{
		UserRepository: repository.NewUserRepositoryPostgres(db),
		DB:             db,
		IDGenerator:    idGenerator.New(),
	}
	sessionService := &services.SessionService{
		DB:                db,
		SessionRepository: repository.NewSessionRepositoryPostgres(db),
		UserRepository:    repository.NewUserRepositoryPostgres(db),
		TokenManager:      token.ManagerImpl{},
		IDGenerator:       idGenerator.New(),
	}
	pwResetService := &services.PasswordResetService{
		DB:             db,
		UserRepository: repository.NewUserRepositoryPostgres(db),
		PasswordReset:  repository.NewPasswordResetPostgres(db),
		TokenManager:   token.ManagerImpl{},
		EmailGateway: gateway.NewEmailMailtrapGateway(gateway.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
		}),
		SessionRepository: repository.NewSessionRepositoryPostgres(db),
		IDGenerator:       idGenerator.New(),
	}
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
	}

	galleryRepository := repository.NewGalleryRepositoryPostgres(db)
	createGalleryUseCase := usecases.NewCreateGalleryUseCase(galleryRepository, idGenerator.New())
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
	}

	fmt.Printf("Starting the server on port %v\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router.New(usersC, galleryController, cfg.CSRF.Key, cfg.CSRF.Secure)))
}
