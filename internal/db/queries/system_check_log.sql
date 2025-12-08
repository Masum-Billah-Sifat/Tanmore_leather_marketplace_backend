-- name: InsertSystemCheckLog :one
INSERT INTO system_check_log (test_label)
VALUES ($1)
RETURNING id, created_at;
