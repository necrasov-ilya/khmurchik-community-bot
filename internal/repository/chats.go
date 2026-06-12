package repository

import (
	"context"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	db       *pgxpool.Pool
	defaults config.DefaultsConfig
}

type ChatSettings struct {
	TelegramID           int64
	Title                string
	ChatType             string
	GreetingEnabled      bool
	GreetingTime         string
	GreetingTimezone     string
	GreetingMessage      string
	LastGreetingSentDate *time.Time
}

func NewChatRepository(db *pgxpool.Pool, defaults config.DefaultsConfig) *ChatRepository {
	return &ChatRepository{db: db, defaults: defaults}
}

func (r *ChatRepository) UpsertChat(c tgbotapi.Chat) error {
	if c.Type != "group" && c.Type != "supergroup" {
		return nil
	}

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO chats (telegram_id, title, chat_type, greeting_enabled, greeting_time, greeting_timezone, greeting_message)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (telegram_id) DO UPDATE SET
			 title = EXCLUDED.title,
			 chat_type = EXCLUDED.chat_type,
			 updated_at = NOW()`,
		c.ID, c.Title, c.Type, r.defaults.GreetingEnabled, r.defaults.GreetingTime, r.defaults.Timezone, r.defaults.GreetingMessage,
	)
	return err
}

func (r *ChatRepository) ListGreetingEnabled(ctx context.Context) ([]ChatSettings, error) {
	rows, err := r.db.Query(ctx,
		`SELECT telegram_id, title, chat_type, greeting_enabled, greeting_time, greeting_timezone, greeting_message, last_greeting_sent_date
		 FROM chats
		 WHERE greeting_enabled = TRUE AND chat_type IN ('group', 'supergroup')`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []ChatSettings
	for rows.Next() {
		var chat ChatSettings
		if err := rows.Scan(&chat.TelegramID, &chat.Title, &chat.ChatType, &chat.GreetingEnabled, &chat.GreetingTime, &chat.GreetingTimezone, &chat.GreetingMessage, &chat.LastGreetingSentDate); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}

func (r *ChatRepository) GetSettings(ctx context.Context, chatID int64) (ChatSettings, error) {
	var chat ChatSettings
	err := r.db.QueryRow(ctx,
		`SELECT telegram_id, title, chat_type, greeting_enabled, greeting_time, greeting_timezone, greeting_message, last_greeting_sent_date
		 FROM chats
		 WHERE telegram_id = $1`,
		chatID,
	).Scan(&chat.TelegramID, &chat.Title, &chat.ChatType, &chat.GreetingEnabled, &chat.GreetingTime, &chat.GreetingTimezone, &chat.GreetingMessage, &chat.LastGreetingSentDate)
	return chat, err
}

func (r *ChatRepository) SetGreetingEnabled(ctx context.Context, chatID int64, enabled bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE chats SET greeting_enabled = $2, updated_at = NOW() WHERE telegram_id = $1`,
		chatID, enabled,
	)
	return err
}

func (r *ChatRepository) SetGreetingTime(ctx context.Context, chatID int64, greetingTime, timezone string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE chats SET greeting_time = $2, greeting_timezone = $3, updated_at = NOW() WHERE telegram_id = $1`,
		chatID, greetingTime, timezone,
	)
	return err
}

func (r *ChatRepository) SetGreetingMessage(ctx context.Context, chatID int64, message string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE chats SET greeting_message = $2, updated_at = NOW() WHERE telegram_id = $1`,
		chatID, message,
	)
	return err
}

func (r *ChatRepository) MarkGreetingSent(ctx context.Context, chatID int64, sentDate time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE chats SET last_greeting_sent_date = $2, updated_at = NOW() WHERE telegram_id = $1`,
		chatID, sentDate,
	)
	return err
}
