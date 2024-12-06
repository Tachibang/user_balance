package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		PG     `yaml:"postgres"`
	}

	Server struct {
		Host string `env-required:"true" yaml:"host" env:"SERVER_HOST"`
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`

		LogLevel string `env:"LOG_LEVEL" yaml:"log_level"`
	}

	PG struct {
		DBHost     string `env-required:"true" yaml:"host" env:"DB_HOST"`
		DBPort     string `env-required:"true" yaml:"port" env:"DB_PORT"`
		DBUser     string `env-required:"true" yaml:"user" env:"DB_USER"`
		DBPassword string `env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		DBName     string `env-required:"true" yaml:"dbname" env:"DB_NAME"`
	}
)

func NewConfig(dotenvPath string) (*Config, error) {
	err := godotenv.Load(dotenvPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке .env файла: %w", err)
	}

	cfg := &Config{}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении переменных окружения: %w", err)
	}

	return cfg, nil
}
