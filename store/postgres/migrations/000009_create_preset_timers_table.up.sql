CREATE TABLE IF NOT EXISTS preset_timers(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    timer_type TEXT, --(daily or periodical)
    description TEXT,
    period_seconds BIGINT, -- (only for periodical timers)
    trigger_time TEXT, -- (only for daily timers)
    link TEXT
);