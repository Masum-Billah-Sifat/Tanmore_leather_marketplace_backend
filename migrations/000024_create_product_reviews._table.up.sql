CREATE TABLE product_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    reviewer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    review_text TEXT NOT NULL,
    review_image_url TEXT,

    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_reviews_product_id ON product_reviews(product_id);
CREATE INDEX idx_product_reviews_reviewer_user_id ON product_reviews(reviewer_user_id);
CREATE INDEX idx_product_reviews_is_archived ON product_reviews(is_archived);
CREATE INDEX idx_product_reviews_is_banned ON product_reviews(is_banned);
