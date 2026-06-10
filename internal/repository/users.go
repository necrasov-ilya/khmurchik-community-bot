package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(u *tgbotapi.User) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (telegram_id, username, first_name, last_name, language_code)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (telegram_id) DO UPDATE SET
			 username = EXCLUDED.username,
			 first_name = EXCLUDED.first_name,
			 last_name = EXCLUDED.last_name,
			 language_code = EXCLUDED.language_code,
			 updated_at = NOW()`,
		u.ID, u.UserName, u.FirstName, u.LastName, u.LanguageCode,
	)
	return err
}

func (r *UserRepository) FindByUsername(username string) (int64, error) {
	var id int64
	err := r.db.QueryRow(context.Background(),
		`SELECT telegram_id FROM users WHERE username = $1 LIMIT 1`,
		username,
	).Scan(&id)
	return id, err
}
