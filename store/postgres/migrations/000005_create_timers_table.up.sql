CREATE TABLE IF NOT EXISTS timers(
    id SERIAL PRIMARY KEY,
    tg_id BIGINT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    timer_type TEXT, --(daily or periodical)
    repeat_type TEXT NOT NULL DEFAULT 'once missed' , --(once missed or always)
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    description TEXT,
    period_seconds BIGINT, -- (only for periodical timers)
    trigger_time TEXT, -- (only for daily timers)
    last_trigger TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    status TEXT NOT NULL DEFAULT 'inactive', --(active or inactive)
    last_ack TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),--(if last_ack is after last_trigger, the timer is considered acknowledged and can be sent to "once missed" subs)
    next_trigger TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    link TEXT
);