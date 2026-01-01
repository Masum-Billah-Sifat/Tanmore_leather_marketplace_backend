-- +goose Up
-- +goose StatementBegin
ALTER TABLE checkout_sessions
    ALTER COLUMN delivery_charge DROP NOT NULL;
-- +goose StatementEnd
