package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseHost     string `env:"DATABASE_HOST" env-default:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" env-default:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" env-default:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" env-default:"postgres"`
	DatabaseDB       string `env:"DATABASE_DB" env-default:"news"`
}

func NewLoadConfig() (Config, error) {
	var cfg Config
	cleanenv.ReadConfig(".env", &cfg)
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
