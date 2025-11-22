package migrations

import (
	"fmt"
	"pr_reviewer_service/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// Migration - миграции
type Migration struct {
	cfg    *config.Config
	logger *zap.Logger
}

// New - конструктор миграций
func New(cfg *config.Config, logger *zap.Logger) *Migration {
	return &Migration{
		cfg:    cfg,
		logger: logger,
	}
}

// RunMigrations выполняет все миграции из папки
func (m *Migration) RunMigrations() error {
	dsn := m.cfg.GetDSN()

	// Парсим DSN в pgx ConnConfig
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Создаём sql.DB
	db := stdlib.OpenDB(*connConfig)
	defer db.Close()

	maxRetries := m.cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 5
	}

	for i := 0; i < maxRetries; i++ {
		if err := goose.Up(db, m.cfg.PackageWithMigrations); err != nil {
			m.logger.Error("Ошибка выполнения миграций",
				zap.Error(err),
				zap.Int("attempt", i+1),
				zap.Int("max_attempts", maxRetries),
			)

			time.Sleep(m.cfg.RetryDelay)
			continue
		}

		m.logger.Info("Все миграции успешно применены")
		return nil
	}

	return fmt.Errorf("не удалось применить миграции после %d попыток", maxRetries)
}
