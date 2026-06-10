CREATE TABLE moderation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    admin_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    action_type TEXT NOT NULL CHECK (action_type IN ('mute', 'ban', 'kick', 'unmute')),
    reason TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_moderation_logs_user_id ON moderation_logs(user_id);
CREATE INDEX idx_moderation_logs_chat_id ON moderation_logs(chat_id);
CREATE INDEX idx_moderation_logs_action_type ON moderation_logs(action_type);
