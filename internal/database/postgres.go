package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/evart2006/khmurchik-community-bot/internal/config"
)

func NewPool(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime != "" {
		maxLife, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			return nil, fmt.Errorf("invalid conn_max_lifetime %q: %w", cfg.ConnMaxLifetime, err)
		}
		poolCfg.MaxConnLifetime = maxLife
	}

	return pgxpool.NewWithConfig(context.Background(), poolCfg)
}
