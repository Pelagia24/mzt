package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB  DB  `mapstructure:"db"`
	Jwt Jwt `mapstructure:"jwt"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

type Jwt struct {
	AccessKey        string        `mapstructure:"access_key"`
	RefreshKey       string        `mapstructure:"refresh_key"`
	AccessExpiresIn  time.Duration `mapstructure:"access_expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
}

func NewConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return &Config{
		DB: DB{
			Host:     os.Getenv("PGHOST"),
			Port:     os.Getenv("PGPORT"),
			User:     os.Getenv("PGUSER"),
			Name:     os.Getenv("PGDATABASE"),
			Password: os.Getenv("PGPASSWORD"),
		},
		Jwt: Jwt{
			AccessKey:        os.Getenv("JWT_ACCESS_KEY"),
			RefreshKey:       os.Getenv("JWT_REFRESH_KEY"),
			AccessExpiresIn:  time.Minute * 30,
			RefreshExpiresIn: time.Hour * 24 * 14,
		},
	}
}
