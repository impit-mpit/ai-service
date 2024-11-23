package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	OpenApiUrl string `env:"OPENAPIURL" env-default:"localhost"`
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
