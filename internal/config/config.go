package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP       `yaml:"http"`
	Postgresql `yaml:"postgresql"`
	JWT        `yaml:"jwt"`
}

type HTTP struct {
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
}

type Postgresql struct {
	User       string `yaml:"user" env:"PG_USER" env-default:"postgres"`
	Password   string `yaml:"password" env:"PG_PASSWORD" env-default:"5432"`
	Host       string `yaml:"host" env:"PG_HOST" env-default:"127.0.0.1"`
	Port       string `yaml:"port" env:"PG_PORT" env-default:"5432"`
	Database   string `yaml:"database" env:"PG_DATABASE" env-default:"postgres"`
	SSLMode    string `yaml:"ssl_mode" env:"PG_SSL" env-default:"disable"`
	AutoCreate bool   `yaml:"auto_create" env:"PG_AUTO_CREATE" env-default:"true"`
}

type JWT struct {
	SignKey  string        `yaml:"sign_key" env:"JWT_KEY"`
	TokenTTL time.Duration `yaml:"token_ttl" env:"JWT_TTL"`
}

func New(log *slog.Logger) (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		log.Error("config path is empty")
		os.Exit(1)
	}

	if _, err := os.Stat(path); err != nil {
		log.Error("failed to open config file", logger.Error(err))
		os.Exit(1)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "sets path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
