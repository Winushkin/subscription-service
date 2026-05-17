//Package config хранит в себе конфигурацию приложения
package config

import (
	"subscription-service/internal/storage"

	"github.com/ilyakaznacheev/cleanenv"
)


// AppConfig содержит конфигурационные переменные приложения.
type AppConfig struct {
	PostgresConfig storage.Config `env-prefix:"POSTGRES_"`
	ServerHost     string         `env:"SERVER_HOST"`
	ServerPort     string         `env:"SERVER_PORT"`
}

// NewAppConfig считывает переменные окружения и возвращает структуру AppConfig.
func NewAppConfig() (*AppConfig, error) {
	var cfg AppConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
