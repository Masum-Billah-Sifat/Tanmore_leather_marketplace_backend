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

-- name: FetchUndispatchedEvents :many
SELECT
  id,
  userid,
  event_type,
  event_payload,
  dispatched_at,
  created_at
FROM events
WHERE dispatched_at IS NULL
ORDER BY created_at ASC
LIMIT $1;  -- Use positional arg for LIMIT


-- name: MarkEventDispatched :exec
UPDATE events
SET dispatched_at = $1
WHERE id = $2
  AND dispatched_at IS NULL;
