ALTER TABLE chats
    ADD COLUMN greeting_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN greeting_time TEXT NOT NULL DEFAULT '10:00',
    ADD COLUMN greeting_timezone TEXT NOT NULL DEFAULT 'Europe/Minsk',
    ADD COLUMN greeting_message TEXT NOT NULL DEFAULT 'Я робот, и Паша заставил меня каждое утро желать вам охуенного дня!

Обнял, покружил, на место поставил! 😘',
    ADD COLUMN last_greeting_sent_date DATE;

CREATE TABLE user_marks (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    mark_type TEXT NOT NULL CHECK (mark_type IN ('balabol')),
    assigned_by BIGINT NOT NULL,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id, mark_type)
);

CREATE INDEX idx_user_marks_chat_user ON user_marks(chat_id, user_id);

CREATE TABLE reports (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    reporter_id BIGINT NOT NULL,
    reported_user_id BIGINT NOT NULL,
    reported_message_id INTEGER NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_chat_status ON reports(chat_id, status);
CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_reported_user_id ON reports(reported_user_id);
