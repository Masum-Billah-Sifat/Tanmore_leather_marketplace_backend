DROP INDEX IF EXISTS idx_orders_order_code;
DROP INDEX IF EXISTS idx_orders_payment_status;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_latest_payment_id;
DROP INDEX IF EXISTS idx_orders_checkout_session_id;
DROP INDEX IF EXISTS idx_orders_user_id;

DROP TABLE IF EXISTS orders;
