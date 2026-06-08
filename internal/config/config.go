package config

import (
	"net"
	"net/url"
	"strconv"
)

type Config struct {
	App   AppConfig      `json:"app"`
	DB    DatabaseConfig `json:"-"`
	Redis RedisConfig    `json:"-"`
	Log   LoggerConfig   `json:"log"`
}

type AppConfig struct {
	Name        string `json:"name"`
	SkinportUrl string `json:"skinport_url"`
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Database string `env:"POSTGRES_DATABASE" env-default:"kolikosoft"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
}

func (c DatabaseConfig) URL() string {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   "/" + c.Database,
	}

	query := databaseURL.Query()
	query.Set("sslmode", c.SSLMode)
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	Database int    `env:"REDIS_DATABASE" env-default:"0"`
}

func (c RedisConfig) URL() string {
	redisURL := url.URL{
		Scheme: "redis",
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   "/" + strconv.Itoa(c.Database),
	}
	if c.Password != "" {
		redisURL.User = url.UserPassword("", c.Password)
	}

	return redisURL.String()
}

type LoggerConfig struct {
	Level string `json:"level"`
}
