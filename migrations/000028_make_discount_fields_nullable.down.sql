-- Restore NOT NULL constraints (may fail if data exists with NULLs)
ALTER TABLE checkout_items
ALTER COLUMN discount_type SET NOT NULL;

ALTER TABLE checkout_items
ALTER COLUMN discount_value SET NOT NULL;
