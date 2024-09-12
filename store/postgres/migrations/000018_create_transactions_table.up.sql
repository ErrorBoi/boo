CREATE TABLE IF NOT EXISTS txs(
                                  id UUID NOT NULL DEFAULT gen_random_uuid(),
                                  tg_id bigint not null,
                                  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                  amount bigint not null,
                                  description text not null
);