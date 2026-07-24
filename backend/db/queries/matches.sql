-- name: CreateMatch :one
INSERT INTO match (
    public_id,
    white_user_id,
    black_user_id,
    white_display_name,
    black_display_name,
    initial_time_ms,
    increment_ms
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: InsertMatchEvent :one
INSERT INTO match_event (
    match_id,
    seq_num,
    event_type,
    payload
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
