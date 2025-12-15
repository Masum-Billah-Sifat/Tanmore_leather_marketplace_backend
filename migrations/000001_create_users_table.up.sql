CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    google_id TEXT UNIQUE NOT NULL,
    primary_email TEXT UNIQUE NOT NULL,

    display_name TEXT,
    profile_image_url TEXT,

    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    is_muted BOOLEAN NOT NULL DEFAULT FALSE,

    current_mode TEXT NOT NULL CHECK (current_mode IN ('customer', 'seller')),

    is_seller_profile_approved BOOLEAN NOT NULL DEFAULT FALSE,
    is_seller_profile_created BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Fast login / lookup
CREATE UNIQUE INDEX idx_users_google_id ON users(google_id);
CREATE UNIQUE INDEX idx_users_primary_email ON users(primary_email);

-- Mode filtering (customer vs seller)
CREATE INDEX idx_users_current_mode ON users(current_mode);

-- Moderation and admin filtering
CREATE INDEX idx_users_is_archived ON users(is_archived);
CREATE INDEX idx_users_is_banned ON users(is_banned);

-- Seller flow performance
CREATE INDEX idx_users_is_seller_profile_approved ON users(is_seller_profile_approved);
