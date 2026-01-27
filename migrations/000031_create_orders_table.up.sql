CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- 🔗 Relationships
    user_id UUID NOT NULL REFERENCES users(id),
    checkout_session_id UUID NOT NULL REFERENCES checkout_sessions(id),

    -- 🔖 Order Identity
    order_code TEXT UNIQUE NOT NULL, -- e.g., TNMR-274192

    -- 💰 Financials
    subtotal NUMERIC NOT NULL,
    shipping_fee NUMERIC DEFAULT 0,
    total_amount NUMERIC NOT NULL,
    currency TEXT DEFAULT 'BDT',

    -- 💳 Payment
    payment_method TEXT NOT NULL DEFAULT 'cod'
        CHECK (payment_method IN ('cod', 'prepaid')),
    payment_status TEXT CHECK (payment_status IN (
        'pending_payment', 'payment_failed', 'payment_success', 'pending_refund'
    )),

    latest_payment_id UUID REFERENCES payments(id),

    -- 📦 Fulfillment
    status TEXT CHECK (status IN (
        'processing', 'shipped', 'delivered', 'cancelled', 'refunded', 'returned'
    )),

    -- 🎁 Discount
    platform_discount_type TEXT,
    platform_discount_value DECIMAL(10,2),
    platform_discount_amount_applied DECIMAL(10,2),
    is_platform_discount_applied BOOLEAN NOT NULL DEFAULT FALSE,

    -- 🕒 Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_checkout_session_id ON orders(checkout_session_id);
CREATE INDEX idx_orders_latest_payment_id ON orders(latest_payment_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_order_code ON orders(order_code);
