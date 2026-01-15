-- Drop triggers
DROP TRIGGER IF EXISTS log_order_status_change_trigger ON orders;
DROP TRIGGER IF EXISTS set_order_number_trigger ON orders;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop functions
DROP FUNCTION IF EXISTS log_order_status_change();
DROP FUNCTION IF EXISTS set_order_number();
DROP FUNCTION IF EXISTS generate_order_number();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS orders;
