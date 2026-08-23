-- +goose Up
CREATE INDEX idx_match_event_match_id_event_type ON match_event(match_id, event_type);

-- +goose Down
DROP INDEX IF EXISTS idx_match_event_match_id_event_type;
