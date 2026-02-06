-- 1. Drop old indexes
DROP INDEX IF EXISTS idx_cart_items_is_active;
DROP INDEX IF EXISTS idx_cart_items_user_id;
DROP INDEX IF EXISTS idx_cart_items_variant_id;

-- 2. Drop check constraints that use the fields we're removing
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_deprecated_guest_ck;
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_user_or_guest_ck;

-- 3. Drop unwanted columns (individually)
ALTER TABLE cart_items DROP COLUMN IF EXISTS guest_user_id;
ALTER TABLE cart_items DROP COLUMN IF EXISTS is_deprecated;

-- 4. Make user_id NOT NULL
ALTER TABLE cart_items ALTER COLUMN user_id SET NOT NULL;

-- 5. Add useful composite index
CREATE INDEX idx_cart_items_user_variant ON cart_items (user_id, variant_id);
