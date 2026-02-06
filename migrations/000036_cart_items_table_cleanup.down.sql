-- 1. Re-add columns
ALTER TABLE cart_items ADD COLUMN guest_user_id uuid;
ALTER TABLE cart_items ADD COLUMN is_deprecated boolean NOT NULL DEFAULT false;

-- 2. Re-add check constraints
ALTER TABLE cart_items ADD CONSTRAINT cart_deprecated_guest_ck
  CHECK (is_deprecated = false OR (is_deprecated = true AND guest_user_id IS NOT NULL AND user_id IS NULL));

ALTER TABLE cart_items ADD CONSTRAINT cart_user_or_guest_ck
  CHECK (user_id IS NOT NULL AND guest_user_id IS NULL OR user_id IS NULL AND guest_user_id IS NOT NULL);

-- 3. Make user_id nullable again
ALTER TABLE cart_items ALTER COLUMN user_id DROP NOT NULL;

-- 4. Recreate old indexes
CREATE INDEX idx_cart_items_is_active ON cart_items (is_active);
CREATE INDEX idx_cart_items_user_id ON cart_items (user_id);
CREATE INDEX idx_cart_items_variant_id ON cart_items (variant_id);

-- 5. Drop new index
DROP INDEX IF EXISTS idx_cart_items_user_variant;
