-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS addresses CASCADE;
DROP TABLE IF EXISTS users CASCADE;
