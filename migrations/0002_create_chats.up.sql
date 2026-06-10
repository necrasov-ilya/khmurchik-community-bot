CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    title TEXT,
    chat_type TEXT NOT NULL DEFAULT 'group',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chats_telegram_id ON chats(telegram_id);
