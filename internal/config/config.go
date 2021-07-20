package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT" envDefault:":8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"15s"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"3306"`
	Username string `env:"DB_USER" envDefault:"user"`
	Password string `env:"DB_PASS" envDefault:"password"`
	DbName   string `env:"DB_NAME" envDefault:"go-db"`
}

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
}

func Get() *Config {
	serverConf := &ServerConfig{}

	if err := env.Parse(serverConf); err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatalf("Server Configuration error...")
	}

	dbConf := &DBConfig{}

	if err := env.Parse(dbConf); err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatalf("Database Configuration error...")
	}

	return &Config{*serverConf, *dbConf}
}
