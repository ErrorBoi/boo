CREATE TABLE IF NOT EXISTS balance(
    tg_id BIGINT UNIQUE NOT NULL,
    amount DECIMAL NOT NULL DEFAULT 0
);

INSERT INTO balance(tg_id, amount) SELECT tg_id, 0 FROM users WHERE tg_id NOT IN (SELECT tg_id FROM balance);