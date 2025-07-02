package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	JWTSecret   string `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment: %w", err)
	}
	return cfg, nil
}
