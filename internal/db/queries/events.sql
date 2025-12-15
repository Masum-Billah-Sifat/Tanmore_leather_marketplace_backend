-- name: InsertEvent :exec
INSERT INTO events (
    id,
    userid,
    event_type,
    event_payload,
    dispatched_at,
    created_at
) VALUES (
    $1,  $2,  $3,  $4,
    $5,  $6
);