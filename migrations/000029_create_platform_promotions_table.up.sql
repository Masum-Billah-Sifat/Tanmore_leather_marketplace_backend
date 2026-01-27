-- ------------------------------------------------------------
-- 📁 File: db/migration/000059_create_platform_promotions.up.sql
-- 🧱 Create platform_promotions table and indexes
-- ------------------------------------------------------------

CREATE TABLE platform_promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- 🎯 Promotion Metadata
    title TEXT NOT NULL,
    description TEXT,

    -- 💸 Discount Logic
    discount_type TEXT NOT NULL CHECK (discount_type IN ('flat', 'percentage')),
    discount_value DECIMAL(10,2) NOT NULL,

    min_cart_value DECIMAL(10,2),
    max_cart_value DECIMAL(10,2),
    max_discount_cap DECIMAL(10,2),

    -- 🕒 Validity Window
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    -- 🧠 Status & Priority
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_archived BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 100 CHECK (priority >= 0),

    -- 🕒 Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    -- 🔐 Logical Constraint
    CHECK (
        min_cart_value IS NULL
        OR max_cart_value IS NULL
        OR min_cart_value <= max_cart_value
    )
);

-- ✅ Ensure no two active, non-archived promos share the same priority
CREATE UNIQUE INDEX uniq_active_priority
ON platform_promotions(priority)
WHERE is_active = true AND is_archived = false;

-- 🔍 Find applicable promos for checkout
CREATE INDEX idx_platform_promotions_active_time
ON platform_promotions (is_active, start_time, end_time);

-- 🔍 Support cart value filtering
CREATE INDEX idx_platform_promotions_cart_range
ON platform_promotions (min_cart_value, max_cart_value);

-- 🔍 Priority lookup (optional if you use ORDER BY priority)
CREATE INDEX idx_platform_promotions_priority
ON platform_promotions(priority);

-- 🔍 Support admin filtering (active + archived)
CREATE INDEX idx_platform_promotions_admin_filters
ON platform_promotions(is_active, is_archived);
