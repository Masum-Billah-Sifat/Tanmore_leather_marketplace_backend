-- ------------------------------------------------------------
-- Make cart_items.required_quantity nullable
-- Reason:
--   Allows soft-deleted cart items (is_active = false)
--   to have no quantity value.
--   Required for remove-cart and clear-cart flows.

ALTER TABLE cart_items
ALTER COLUMN required_quantity DROP NOT NULL;
