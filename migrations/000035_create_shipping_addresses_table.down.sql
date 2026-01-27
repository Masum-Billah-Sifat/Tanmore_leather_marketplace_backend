-- Drop indexes first
DROP INDEX IF EXISTS idx_shipping_addresses_recipient_phone;
DROP INDEX IF EXISTS idx_shipping_addresses_location;
DROP INDEX IF EXISTS idx_shipping_addresses_checkout_session_id;

-- Drop the table
DROP TABLE IF EXISTS shipping_addresses;
