BEGIN;

CREATE TABLE IF NOT EXISTS outbox
(
    event_id        UUID PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_type      TEXT        NOT NULL,
    occurred_at     TIMESTAMPTZ NOT NULL,
    value           BYTEA       NOT NULL,
    trace_context   BYTEA       NOT NULL
);

-- индекс для быстрого поиска по типу события и времени
CREATE INDEX IF NOT EXISTS idx_outbox_created_at
    ON outbox (created_at);

COMMIT;