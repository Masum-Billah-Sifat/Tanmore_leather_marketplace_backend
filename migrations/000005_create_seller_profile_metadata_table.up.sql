CREATE TABLE seller_profile_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    sellerstorename TEXT NOT NULL,
    sellercontactno TEXT NOT NULL,
    sellerwhatsappcontactno TEXT NOT NULL,

    sellerwebsitelink TEXT,          -- optional
    sellerfacebookpagename TEXT,     -- optional
    selleremail TEXT,                -- optional

    sellerphysicallocation TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Every seller has one profile: enforce fast lookup
CREATE UNIQUE INDEX idx_seller_profile_seller_id ON seller_profile_metadata(seller_id);

-- For searching or filtering sellers
CREATE INDEX idx_seller_profile_name ON seller_profile_metadata(sellerstorename);

-- Contact lookup audit
CREATE INDEX idx_seller_profile_contact ON seller_profile_metadata(sellercontactno);
