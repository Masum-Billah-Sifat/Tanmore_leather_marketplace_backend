-- ------------------------------------------------------------
-- Rollback: Restore NOT NULL constraint on required_quantity
-- ------------------------------------------------------------

ALTER TABLE cart_items
ALTER COLUMN required_quantity SET NOT NULL;
