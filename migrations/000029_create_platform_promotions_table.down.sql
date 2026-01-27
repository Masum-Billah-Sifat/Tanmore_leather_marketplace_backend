-- ------------------------------------------------------------
-- 📁 File: db/migration/000059_create_platform_promotions.down.sql
-- 🔁 Rollback: Drop platform_promotions table and indexes
-- ------------------------------------------------------------

DROP INDEX IF EXISTS uniq_active_priority;
DROP INDEX IF EXISTS idx_platform_promotions_active_time;
DROP INDEX IF EXISTS idx_platform_promotions_cart_range;
DROP INDEX IF EXISTS idx_platform_promotions_priority;
DROP INDEX IF EXISTS idx_platform_promotions_admin_filters;

DROP TABLE IF EXISTS platform_promotions;
