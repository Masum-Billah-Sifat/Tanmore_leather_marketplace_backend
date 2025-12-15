CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    device_fingerprint TEXT NOT NULL,

    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_device_fingerprint ON user_sessions(device_fingerprint);
CREATE INDEX idx_sessions_is_revoked ON user_sessions(is_revoked);