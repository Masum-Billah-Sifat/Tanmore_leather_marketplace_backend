ALTER TABLE checkout_items
ALTER COLUMN discount_type DROP NOT NULL;

ALTER TABLE checkout_items
ALTER COLUMN discount_value DROP NOT NULL;
