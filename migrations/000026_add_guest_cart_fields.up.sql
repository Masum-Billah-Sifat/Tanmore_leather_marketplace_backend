-- +goose Up
-- +goose StatementBegin

ALTER TABLE cart_items
    ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE cart_items
    ADD COLUMN guest_user_id UUID;

ALTER TABLE cart_items
    ADD COLUMN is_deprecated BOOLEAN NOT NULL DEFAULT false;

-- Enforce: Either user_id or guest_user_id must be present, but not both
ALTER TABLE cart_items
    ADD CONSTRAINT cart_user_or_guest_ck
    CHECK (
        (user_id IS NOT NULL AND guest_user_id IS NULL)
        OR (user_id IS NULL AND guest_user_id IS NOT NULL)
    );

-- Enforce: is_deprecated = true only allowed on guest_user rows
ALTER TABLE cart_items
    ADD CONSTRAINT cart_deprecated_guest_ck
    CHECK (
        is_deprecated = false
        OR (is_deprecated = true AND guest_user_id IS NOT NULL AND user_id IS NULL)
    );

-- +goose StatementEnd
