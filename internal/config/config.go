package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config описывает структуру конфига
type Config struct {
	DatabaseDSN         string `envconfig:"DATABASE_DSN"`
	SentryDSN           string `envconfig:"SENTRY_DSN"`
	AccessTokenKey      string `envconfig:"ACCESS_TOKEN_KEY"`
	RefreshTokenKey     string `envconfig:"REFRESH_TOKEN_KEY"`
	NotificationBaseURL string `envconfig:"NOTIFICATION_BASE_URL"`
}

// InitConfig возвращает конфиг
func InitConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)

	return &cfg, err
}

// MustInitConfig возвращает конфиг или паникует при ошибке
func MustInitConfig() *Config {
	cfg, err := InitConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
