ALTER TABLE payments
ADD CONSTRAINT fk_payments_order_id
FOREIGN KEY (order_id)
REFERENCES orders(id)
ON DELETE CASCADE;
