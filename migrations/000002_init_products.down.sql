-- Drop triggers
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_sizes_updated_at ON sizes;
DROP TRIGGER IF EXISTS log_sizes_inventory_change ON sizes;

-- Drop function
DROP FUNCTION IF EXISTS log_inventory_change();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS inventory_transactions CASCADE;
DROP TABLE IF EXISTS sizes CASCADE;
DROP TABLE IF EXISTS product_images CASCADE;
DROP TABLE IF EXISTS products CASCADE;
