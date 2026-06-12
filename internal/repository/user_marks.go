package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const MarkBalabol = "balabol"

type UserMarkRepository struct {
	db *pgxpool.Pool
}

func NewUserMarkRepository(db *pgxpool.Pool) *UserMarkRepository {
	return &UserMarkRepository{db: db}
}

func (r *UserMarkRepository) SetMark(ctx context.Context, chatID, userID, assignedBy int64, markType, reason string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO user_marks (chat_id, user_id, mark_type, assigned_by, reason)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (chat_id, user_id, mark_type) DO UPDATE SET
			 assigned_by = EXCLUDED.assigned_by,
			 reason = EXCLUDED.reason,
			 created_at = NOW()`,
		chatID, userID, markType, assignedBy, reason,
	)
	return err
}

func (r *UserMarkRepository) RemoveMark(ctx context.Context, chatID, userID int64, markType string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM user_marks WHERE chat_id = $1 AND user_id = $2 AND mark_type = $3`,
		chatID, userID, markType,
	)
	return err
}

func (r *UserMarkRepository) HasMark(ctx context.Context, chatID, userID int64, markType string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM user_marks WHERE chat_id = $1 AND user_id = $2 AND mark_type = $3
		)`,
		chatID, userID, markType,
	).Scan(&exists)
	return exists, err
}
