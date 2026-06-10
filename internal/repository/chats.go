package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) UpsertChat(c tgbotapi.Chat) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO chats (telegram_id, title, chat_type)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (telegram_id) DO UPDATE SET
			 title = EXCLUDED.title,
			 updated_at = NOW()`,
		c.ID, c.Title, c.Type,
	)
	return err
}
