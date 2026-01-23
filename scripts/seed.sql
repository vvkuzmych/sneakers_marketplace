-- Seed data for sneakers_marketplace database
-- Run with: psql -U postgres -d sneakers_marketplace -f scripts/seed.sql

\echo '🌱 Seeding database...'

-- Clear existing data (in correct order to respect foreign keys)
\echo 'Clearing existing data...'
TRUNCATE TABLE transaction_fees, matches, bids, asks, sizes, products, users CASCADE;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE sizes_id_seq RESTART WITH 1;
ALTER SEQUENCE bids_id_seq RESTART WITH 1;
ALTER SEQUENCE asks_id_seq RESTART WITH 1;
ALTER SEQUENCE matches_id_seq RESTART WITH 1;
ALTER SEQUENCE transaction_fees_id_seq RESTART WITH 1;

-- Insert Users
\echo 'Inserting users...'
INSERT INTO users (email, password_hash, first_name, last_name, phone, is_active, created_at, updated_at)
VALUES
    ('john@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J7Z.EMoYiH6s7mS9C9b8rQZ8WVrXYi', 'John', 'Doe', '+1234567890', true, NOW(), NOW()),
    ('jane@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J7Z.EMoYiH6s7mS9C9b8rQZ8WVrXYi', 'Jane', 'Smith', '+1234567891', true, NOW(), NOW()),
    ('bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J7Z.EMoYiH6s7mS9C9b8rQZ8WVrXYi', 'Bob', 'Wilson', '+1234567892', true, NOW(), NOW()),
    ('alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J7Z.EMoYiH6s7mS9C9b8rQZ8WVrXYi', 'Alice', 'Johnson', '+1234567893', true, NOW(), NOW()),
    ('test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J7Z.EMoYiH6s7mS9C9b8rQZ8WVrXYi', 'Test', 'User', '+1234567894', true, NOW(), NOW());
-- Password for all users: password123

\echo 'Inserted 5 users'

-- Insert Products
\echo 'Inserting products...'
INSERT INTO products (name, brand, sku, description, retail_price, category, image_url, is_active, created_at, updated_at)
VALUES
    ('Air Jordan 1 Retro High OG', 'Nike', 'AJ1-001', 'Classic Air Jordan 1 in Chicago colorway', 170.00, 'sneakers', 'https://example.com/aj1.jpg', true, NOW(), NOW()),
    ('Air Jordan 4 Retro', 'Nike', 'AJ4-001', 'Air Jordan 4 in Bred colorway', 200.00, 'sneakers', 'https://example.com/aj4.jpg', true, NOW(), NOW()),
    ('Yeezy Boost 350 V2', 'Adidas', 'YZY-350', 'Yeezy Boost 350 V2 Zebra', 220.00, 'sneakers', 'https://example.com/yeezy.jpg', true, NOW(), NOW()),
    ('Nike Dunk Low', 'Nike', 'DUNK-001', 'Nike Dunk Low Panda', 110.00, 'sneakers', 'https://example.com/dunk.jpg', true, NOW(), NOW()),
    ('Air Max 1', 'Nike', 'AM1-001', 'Air Max 1 OG Red', 140.00, 'sneakers', 'https://example.com/am1.jpg', true, NOW(), NOW()),
    ('Air Force 1', 'Nike', 'AF1-001', 'Air Force 1 Low White', 90.00, 'sneakers', 'https://example.com/af1.jpg', true, NOW(), NOW()),
    ('New Balance 550', 'New Balance', 'NB550-001', 'New Balance 550 White Green', 120.00, 'sneakers', 'https://example.com/nb550.jpg', true, NOW(), NOW()),
    ('Travis Scott x Air Jordan 1', 'Nike', 'TS-AJ1', 'Travis Scott collaboration', 1500.00, 'sneakers', 'https://example.com/ts-aj1.jpg', true, NOW(), NOW());

\echo 'Inserted 8 products'

-- Insert Sizes for each product
\echo 'Inserting sizes...'
DO $$
DECLARE
    product_record RECORD;
    size_value NUMERIC;
BEGIN
    FOR product_record IN SELECT id FROM products LOOP
        FOR size_value IN 7..13 LOOP
            INSERT INTO sizes (product_id, size, quantity, created_at, updated_at)
            VALUES (product_record.id, size_value, 100, NOW(), NOW());
        END LOOP;
    END LOOP;
END $$;

\echo 'Inserted sizes for all products (7-13)'

-- Insert sample BIDs
\echo 'Inserting sample bids...'
INSERT INTO bids (user_id, product_id, size_id, price, quantity, status, created_at, updated_at)
VALUES
    (1, 1, 4, 250.00, 1, 'active', NOW(), NOW()),
    (2, 1, 5, 240.00, 1, 'active', NOW(), NOW()),
    (3, 2, 4, 300.00, 1, 'active', NOW(), NOW()),
    (1, 3, 6, 350.00, 1, 'active', NOW(), NOW()),
    (4, 4, 5, 150.00, 1, 'active', NOW(), NOW());

\echo 'Inserted 5 sample bids'

-- Insert sample ASKs
\echo 'Inserting sample asks...'
INSERT INTO asks (user_id, product_id, size_id, price, quantity, status, created_at, updated_at)
VALUES
    (2, 1, 4, 280.00, 1, 'active', NOW(), NOW()),
    (3, 1, 5, 270.00, 1, 'active', NOW(), NOW()),
    (4, 2, 4, 320.00, 1, 'active', NOW(), NOW()),
    (5, 3, 6, 380.00, 1, 'active', NOW(), NOW()),
    (2, 4, 5, 170.00, 1, 'active', NOW(), NOW());

\echo 'Inserted 5 sample asks'

-- Summary
\echo ''
\echo '✅ Database seeded successfully!'
\echo ''
\echo '📊 Summary:'
SELECT 
    (SELECT COUNT(*) FROM users) AS users,
    (SELECT COUNT(*) FROM products) AS products,
    (SELECT COUNT(*) FROM sizes) AS sizes,
    (SELECT COUNT(*) FROM bids) AS bids,
    (SELECT COUNT(*) FROM asks) AS asks;

\echo ''
\echo '👤 Test Users (all passwords: password123):'
\echo '   - john@example.com'
\echo '   - jane@example.com'
\echo '   - bob@example.com'
\echo '   - alice@example.com'
\echo '   - test@example.com'
\echo ''
