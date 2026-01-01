CREATE TABLE product_review_replies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    review_id UUID NOT NULL UNIQUE REFERENCES product_reviews(id) ON DELETE CASCADE,
    seller_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    reply_text TEXT NOT NULL,
    reply_image_url TEXT,

    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_replies_review_id ON product_review_replies(review_id);
CREATE INDEX idx_review_replies_seller_user_id ON product_review_replies(seller_user_id);