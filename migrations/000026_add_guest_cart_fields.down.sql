-- +goose Down
-- +goose StatementBegin

ALTER TABLE cart_items
    DROP CONSTRAINT IF EXISTS cart_deprecated_guest_ck;

ALTER TABLE cart_items
    DROP CONSTRAINT IF EXISTS cart_user_or_guest_ck;

ALTER TABLE cart_items
    DROP COLUMN IF EXISTS is_deprecated;

ALTER TABLE cart_items
    DROP COLUMN IF EXISTS guest_user_id;

ALTER TABLE cart_items
    ALTER COLUMN user_id SET NOT NULL;

-- +goose StatementEnd
