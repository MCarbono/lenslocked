package app

import (
	"lenslocked/infra/database"
	"lenslocked/infra/gateway"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	PSQL database.PostgresConfig
	SMTP gateway.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Port string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = database.DefaultPostgresConfig()
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	csrfSecure := os.Getenv("CSRF_SECURE")
	cfg.CSRF.Secure, err = strconv.ParseBool(csrfSecure)
	if err != nil {
		return cfg, err
	}
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	return cfg, nil
}
