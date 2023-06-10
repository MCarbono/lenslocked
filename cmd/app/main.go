package app

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked/controllers"
	"lenslocked/gateway"
	"lenslocked/migrations"
	"lenslocked/models"
	"lenslocked/token"

	repository "lenslocked/repository/postgres"
)

func Start() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	db, err := Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")
	userService := &models.UserService{
		UserRepository: repository.NewUserRepositoryPostgres(db),
		EmailGateway: gateway.NewEmailMailtrapGateway(gateway.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
		}),
		DB: db,
	}
	sessionService := &models.SessionService{
		DB:                db,
		SessionRepository: repository.NewSessionRepositoryPostgres(db),
		UserRepository:    repository.NewUserRepositoryPostgres(db),
		TokenManager:      token.ManagerImpl{},
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
	}
	fmt.Printf("Starting the server on port %v\n", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, controllers.NewRouter(usersC, cfg.CSRF.Key, cfg.CSRF.Secure)))
}
