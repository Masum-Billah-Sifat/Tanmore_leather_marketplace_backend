-- +goose Up
-- +goose StatementBegin

ALTER TABLE checkout_sessions
ADD COLUMN status TEXT NOT NULL DEFAULT 'draft'
CHECK (status IN ('draft', 'awaiting_shipping_info', 'ready_to_order', 'abandoned', 'order_created'));

ALTER TABLE checkout_sessions
ADD COLUMN platform_discount_type TEXT;

ALTER TABLE checkout_sessions
ADD COLUMN platform_discount_value DECIMAL(10,2);

ALTER TABLE checkout_sessions
ADD COLUMN platform_discount_amount_applied DECIMAL(10,2);

ALTER TABLE checkout_sessions
ADD COLUMN is_platform_discount_applied BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE checkout_sessions
ADD COLUMN payment_method TEXT NOT NULL DEFAULT 'cod'
CHECK (payment_method IN ('cod', 'prepaid'));

-- +goose StatementEnd
