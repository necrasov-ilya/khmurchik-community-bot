DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS user_marks;

ALTER TABLE chats
    DROP COLUMN IF EXISTS last_greeting_sent_date,
    DROP COLUMN IF EXISTS greeting_message,
    DROP COLUMN IF EXISTS greeting_timezone,
    DROP COLUMN IF EXISTS greeting_time,
    DROP COLUMN IF EXISTS greeting_enabled;
