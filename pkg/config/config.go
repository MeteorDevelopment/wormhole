package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/pkg/errors"
	"log"
)

type Config struct {
	ApiBase string `env:"API_BASE"`
	Port    string `env:"PORT"`

	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbName     string `env:"DB_NAME"`

	EmailUsername string `env:"EMAIL_USERNAME"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
	EmailHost     string `env:"EMAIL_HOST"`
	EmailPort     string `env:"EMAIL_PORT"`

	JwtSigningKey     string `env:"JWT_SIGNING_KEY"`
	JwtExpirationTime int64  `env:"JWT_EXPIRATION_TIME"`
}

var config *Config

func Get() *Config {
	return config
}

func Init() {
	config = &Config{}
	err := env.Parse(config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse private config"))
	}
}
