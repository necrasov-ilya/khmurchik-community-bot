package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Report struct {
	ID                int64
	ChatID            int64
	ReporterID        int64
	ReportedUserID    int64
	ReportedMessageID int
	Reason            string
	Status            string
	CreatedAt         time.Time
}

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, report Report) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		`INSERT INTO reports (chat_id, reporter_id, reported_user_id, reported_message_id, reason, status)
		 VALUES ($1, $2, $3, $4, $5, 'open')
		 RETURNING id`,
		report.ChatID, report.ReporterID, report.ReportedUserID, report.ReportedMessageID, report.Reason,
	).Scan(&id)
	return id, err
}
