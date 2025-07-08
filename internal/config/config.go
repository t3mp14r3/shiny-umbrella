package config

import (
	"fmt"
	"log"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
    AppAddr         string  `env:"APP_ADDR"`

    DBUser          string  `env:"DB_USER"`
    DBHost          string  `env:"DB_HOST"`
    DBName          string  `env:"DB_NAME"`
    DBPass          string  `env:"DB_PASS"`
    DBPort          int     `env:"DB_PORT"`
    DBMigrations    string  `env:"DB_MIGRATIONS"`
}

func New() (*Config, error) {
    var cfg Config

    err := env.Parse(&cfg)

    if err != nil {
        log.Printf("Failed to initialize new config: %v\n", err)
        return nil, err
    }

    return &cfg, nil
}

func (c *Config) RepositoryConnString() string {
    return fmt.Sprintf("user=%s dbname=%s host=%s port=%d password=%s sslmode=disable",
        c.DBUser,
        c.DBName,
        c.DBHost,
        c.DBPort,
        c.DBPass,
    )
}
