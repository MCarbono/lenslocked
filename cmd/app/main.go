package app

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked/controllers"
	"lenslocked/gateway"
	"lenslocked/migrations"
	"lenslocked/services"
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
	userService := &services.UserService{
		UserRepository: repository.NewUserRepositoryPostgres(db),
		DB:             db,
	}
	sessionService := &services.SessionService{
		DB:                db,
		SessionRepository: repository.NewSessionRepositoryPostgres(db),
		UserRepository:    repository.NewUserRepositoryPostgres(db),
		TokenManager:      token.ManagerImpl{},
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
	}
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
	}
	fmt.Printf("Starting the server on port %v\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, controllers.NewRouter(usersC, cfg.CSRF.Key, cfg.CSRF.Secure)))
}
