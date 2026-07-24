-- +goose Up
CREATE TABLE match (
    id                 BIGSERIAL PRIMARY KEY,
    public_id          TEXT NOT NULL UNIQUE,
    white_user_id      BIGINT REFERENCES users(id) ON DELETE SET NULL,
    black_user_id      BIGINT REFERENCES users(id) ON DELETE SET NULL,
    white_display_name TEXT NOT NULL,
    black_display_name TEXT NOT NULL,
    initial_time_ms    BIGINT NOT NULL CHECK (initial_time_ms > 0),
    increment_ms       BIGINT NOT NULL CHECK (increment_ms >= 0),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_match_white_user_id ON match(white_user_id) WHERE white_user_id IS NOT NULL;
CREATE INDEX idx_match_black_user_id ON match(black_user_id) WHERE black_user_id IS NOT NULL;
CREATE INDEX idx_match_created_at ON match(created_at);

CREATE TYPE match_event_type AS ENUM (
    'game_started',
    'move',
    'game_ended'
);

CREATE TABLE match_event (
    id         BIGSERIAL PRIMARY KEY,
    match_id   BIGINT NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    seq_num    INTEGER NOT NULL CHECK (seq_num > 0),
    event_type match_event_type NOT NULL,
    payload    JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(match_id, seq_num)
);

-- +goose Down
DROP TABLE IF EXISTS match_event;
DROP TABLE IF EXISTS match;
DROP TYPE IF EXISTS match_event_type;
