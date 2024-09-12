CREATE TABLE IF NOT EXISTS tasks_users(
    task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    tg_id BIGINT,
    created_at VARCHAR (50) NOT NULL,
    CONSTRAINT tasks_users_pkey PRIMARY KEY ( task_id, tg_id )
);