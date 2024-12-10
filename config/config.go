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
		Kafka  `yaml:"kafka"`
		Cron   `yaml:"cron"`
	}

	Server struct {
		Host  string `env-required:"true" yaml:"host" env:"SERVER_HOST"`
		Port  string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		DBHost     string `env-required:"true" yaml:"host" env:"DB_HOST"`
		DBPort     string `env-required:"true" yaml:"port" env:"DB_PORT"`
		DBUser     string `env-required:"true" yaml:"user" env:"DB_USER"`
		DBPassword string `env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		DBName     string `env-required:"true" yaml:"dbname" env:"DB_NAME"`
	}

	Kafka struct {
		Brokers string `env-required:"true" yaml:"brokers" env:"KAFKA_BROKERS"`
		Topic   string `env-required:"true" yaml:"topic" env:"KAFKA_TOPIC"`
	}

	Cron struct {
		Schedule string `env-required:"true" yaml:"schedule" env:"CRON_SCHEDULE"`
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
