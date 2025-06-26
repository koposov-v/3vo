package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Name                  string        `env:"SERVICE_NAME,required"`
	GracefulShutdownDelay time.Duration `env:"SERVICE_GRACEFUL_SHUTDOWN_DELAY" envDefault:"5s"`
	HTTP                  Server
	JWTSecret             string `env:"JWT_SECRET"`
}

type Server struct {
	Host            string        `env:"SERVICE_HTTP_HOST" envDefault:"0.0.0.0"`
	Port            string        `env:"SERVICE_HTTP_PORT,required"`
	ShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"20s"`
}

func Init() (*Config, error) {
	var cfg Config

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to find .env file at path: .env")
	}

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
