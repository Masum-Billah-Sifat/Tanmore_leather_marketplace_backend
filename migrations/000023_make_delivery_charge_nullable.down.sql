-- +goose Down
-- +goose StatementBegin
ALTER TABLE checkout_sessions
    ALTER COLUMN delivery_charge SET NOT NULL;
-- +goose StatementEnd
