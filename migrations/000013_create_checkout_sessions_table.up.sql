CREATE TABLE checkout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    subtotal DECIMAL(10,2) NOT NULL,
    total_weight_grams INTEGER NOT NULL,
    delivery_charge DECIMAL(10,2) NOT NULL,
    total_payable DECIMAL(10,2) NOT NULL,

    shipping_address_id UUID REFERENCES shipping_addresses(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_checkout_sessions_user_id ON checkout_sessions(user_id);
CREATE INDEX idx_checkout_sessions_shipping_address_id ON checkout_sessions(shipping_address_id);
