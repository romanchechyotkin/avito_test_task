package config

import (
	"flag"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// todo 127.0.0.1 => 0.0.0.0
type Config struct {
	HTTP struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
		Host string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	} `yaml:"http"`
	Postgresql struct {
		User       string `yaml:"user" env:"PG_USER" env-default:"postgres"`
		Password   string `yaml:"password" env:"PG_PASSWORD" env-default:"5432"`
		Host       string `yaml:"host" env:"PG_HOST" env-default:"127.0.0.1"`
		Port       string `yaml:"port" env:"PG_PORT" env-default:"5432"`
		Database   string `yaml:"database" env:"PG_DATABASE" env-default:"postgres"`
		SSLMode    string `yaml:"ssl_mode" env:"PG_SSL" env-default:"disable"`
		AutoCreate bool   `yaml:"auto_create" env:"PG_AUTO_CREATE" env-default:"true"`
	} `yaml:"postgresql"`
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
