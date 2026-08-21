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

-- name: GetMatchByPublicID :one
SELECT * FROM match
WHERE public_id = $1;

-- name: GetMatchEventsByMatchID :many
SELECT * FROM match_event
WHERE match_id = $1
ORDER BY seq_num ASC;

-- name: GetUserEndedMatches :many
SELECT 
    m.public_id,
    m.white_user_id,
    m.black_user_id,
    m.white_display_name,
    m.black_display_name,
    wu.username AS white_username,
    bu.username AS black_username,
    m.initial_time_ms,
    m.increment_ms,
    m.created_at,
    e.payload AS ended_payload,
    (SELECT COUNT(*) FROM match_event me WHERE me.match_id = m.id AND me.event_type = 'move') AS move_count
FROM match m
INNER JOIN match_event e ON e.match_id = m.id AND e.event_type = 'game_ended'
LEFT JOIN users wu ON wu.id = m.white_user_id
LEFT JOIN users bu ON bu.id = m.black_user_id
WHERE m.white_user_id = @user_id OR m.black_user_id = @user_id
ORDER BY m.created_at DESC;
