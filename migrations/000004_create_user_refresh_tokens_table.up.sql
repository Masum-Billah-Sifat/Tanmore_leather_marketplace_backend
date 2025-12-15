CREATE TABLE user_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES user_sessions(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL,

    deprecated_reason TEXT,
    is_deprecated BOOLEAN NOT NULL DEFAULT FALSE,
    deprecated_at TIMESTAMPTZ,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Validate tokens fast
CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON user_refresh_tokens(token_hash);

-- Token cleanup & expiration
CREATE INDEX idx_refresh_tokens_expires_at ON user_refresh_tokens(expires_at);

-- Token state
CREATE INDEX idx_refresh_tokens_is_deprecated ON user_refresh_tokens(is_deprecated);

-- Common queries for session or user
CREATE INDEX idx_refresh_tokens_user_id ON user_refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_session_id ON user_refresh_tokens(session_id);


