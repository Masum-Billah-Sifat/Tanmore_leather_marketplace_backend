-- name: InsertUserModeSwitchLog :exec
INSERT INTO user_mode_history (
  id, user_id, from_mode, to_mode, switched_at
) VALUES (
  $1, $2, $3, $4, $5
);