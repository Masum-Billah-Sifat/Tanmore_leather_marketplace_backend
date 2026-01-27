-- +goose Down
-- +goose StatementBegin

ALTER TABLE checkout_sessions
DROP COLUMN payment_method;

ALTER TABLE checkout_sessions
DROP COLUMN is_platform_discount_applied;

ALTER TABLE checkout_sessions
DROP COLUMN platform_discount_amount_applied;

ALTER TABLE checkout_sessions
DROP COLUMN platform_discount_value;

ALTER TABLE checkout_sessions
DROP COLUMN platform_discount_type;

ALTER TABLE checkout_sessions
DROP COLUMN status;

-- +goose StatementEnd
