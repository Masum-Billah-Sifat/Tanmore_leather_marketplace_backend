-- Rollback: remove the seller_store_name column
ALTER TABLE checkout_items
DROP COLUMN seller_store_name;
