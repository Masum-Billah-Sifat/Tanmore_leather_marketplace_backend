CREATE TABLE user_mode_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    from_mode TEXT NOT NULL CHECK (from_mode IN ('customer', 'seller')),
    to_mode   TEXT NOT NULL CHECK (to_mode   IN ('customer', 'seller')),

    switched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Query mode switches by user (very common for audit trails)
CREATE INDEX idx_user_mode_history_user_id ON user_mode_history(user_id);

-- Query by switch timestamp (admin timelines)
CREATE INDEX idx_user_mode_history_switched_at ON user_mode_history(switched_at);