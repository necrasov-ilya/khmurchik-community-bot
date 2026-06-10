package moderation

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ModerationLog struct {
	ID         int64
	UserID     int64
	AdminID    int64
	ChatID     int64
	ActionType string
	Reason     string
	ExpiresAt  *time.Time
	CreatedAt  time.Time
}

type Repository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewRepository(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) LogAction(ctx context.Context, log ModerationLog) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO moderation_logs (user_id, admin_id, chat_id, action_type, reason, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		log.UserID, log.AdminID, log.ChatID, log.ActionType, log.Reason, log.ExpiresAt,
	)
	return err
}
