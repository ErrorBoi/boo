CREATE TABLE IF NOT EXISTS tasks(
    id BIGSERIAL PRIMARY KEY,
    reward BIGINT NOT NULL,
    reward_type TEXT NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    description TEXT NOT NULL DEFAULT 'complete the task',
    status TEXT NOT NULL DEFAULT 'pending',
    chat_id BIGINT NOT NULL
);