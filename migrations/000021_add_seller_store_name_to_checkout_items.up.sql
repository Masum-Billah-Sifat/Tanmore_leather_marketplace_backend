-- Add a new column to store the seller's store name
ALTER TABLE checkout_items
ADD COLUMN seller_store_name TEXT NOT NULL DEFAULT '';
