package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

// Config - структура конфига
type Config struct {
	DbName                string        `env:"DB_NAME" env-default:"postgres"`
	DbUser                string        `env:"DB_USER" env-default:"postgres"`
	DbPassword            string        `env:"DB_PASSWORD" env-default:"postgres"`
	DbHost                string        `env:"DB_HOST" env-default:"localhost"`
	DbPort                int           `env:"DB_PORT" env-default:"5432"`
	MaxRetries            int           `env:"MAX_RETRIES" env-default:"5"`
	RetryDelay            time.Duration `env:"RETRY_DELAY" env-default:"3s"`
	PackageWithMigrations string        `env:"PACKAGE_WITH_MIGRATIONS" env-default:"./migrations"`
}

// New - конструктор конфига
func New(logger *zap.Logger) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		logger.Error("Error reading config", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully read config")

	return &cfg, nil
}

// GetDSN - подключение для бд
func (cfg *Config) GetDSN() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		strconv.Itoa(cfg.DbPort),
		cfg.DbName,
	)

	return dsn
}
