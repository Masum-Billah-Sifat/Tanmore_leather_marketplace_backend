CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID, -- remove FK for now


    ssl_session_id TEXT,
    ssl_trx_id TEXT,
    bank_trx_id TEXT,

    amount NUMERIC NOT NULL,
    currency TEXT DEFAULT 'BDT',
    payment_method TEXT, -- e.g. card, bkash, nagad

    status TEXT CHECK (status IN (
        'pending', 'success', 'failed', 'cancelled', 'refunded'
    )) DEFAULT 'pending',

    risk_level TEXT,
    gateway_response JSONB,
    validation_response JSONB,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Optional indexes if needed later:
-- CREATE INDEX idx_payments_order_id ON payments(order_id);
-- CREATE INDEX idx_payments_status ON payments(status);
