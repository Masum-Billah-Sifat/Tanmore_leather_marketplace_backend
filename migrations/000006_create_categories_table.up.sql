CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,

    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,

    level INT NOT NULL DEFAULT 0,
    is_leaf BOOLEAN NOT NULL DEFAULT FALSE,

    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_leaf ON categories(is_leaf);
CREATE INDEX idx_categories_is_archived ON categories(is_archived);
CREATE INDEX idx_categories_level ON categories(level);

