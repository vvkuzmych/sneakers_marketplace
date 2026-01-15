-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(255) NOT NULL,
    color VARCHAR(100),
    description TEXT,
    category VARCHAR(100),
    release_year INT,
    retail_price DECIMAL(10, 2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for products
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_model ON products(model);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_products_search ON products USING GIN (to_tsvector('english', name || ' ' || brand || ' ' || model || ' ' || COALESCE(color, '')));

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    display_order INT NOT NULL DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on product_id
CREATE INDEX idx_product_images_product_id ON product_images(product_id, display_order);
CREATE INDEX idx_product_images_primary ON product_images(product_id, is_primary) WHERE is_primary = TRUE;

-- Create sizes table (available sizes for products)
CREATE TABLE IF NOT EXISTS sizes (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size VARCHAR(10) NOT NULL,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved INT NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, size),
    CHECK (reserved <= quantity)
);

-- Create indexes for sizes
CREATE INDEX idx_sizes_product_id ON sizes(product_id);
CREATE INDEX idx_sizes_available ON sizes(product_id, size) WHERE (quantity - reserved) > 0;

-- Create inventory_transactions table (for audit trail)
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id BIGSERIAL PRIMARY KEY,
    size_id BIGINT NOT NULL REFERENCES sizes(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('addition', 'sale', 'removal', 'reservation', 'release')),
    quantity_change INT NOT NULL,
    quantity_before INT NOT NULL,
    quantity_after INT NOT NULL,
    reference_id VARCHAR(100), -- order_id or other reference
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on size_id and created_at for audit queries
CREATE INDEX idx_inventory_transactions_size_id ON inventory_transactions(size_id, created_at DESC);
CREATE INDEX idx_inventory_transactions_reference ON inventory_transactions(reference_id);

-- Create triggers for updated_at
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sizes_updated_at BEFORE UPDATE ON sizes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Note: Inventory transactions are logged manually in the application code
-- This provides better control and audit trail

-- Comments for documentation
COMMENT ON TABLE products IS 'Sneaker products catalog';
COMMENT ON TABLE product_images IS 'Product images (multiple per product)';
COMMENT ON TABLE sizes IS 'Product sizes and inventory';
COMMENT ON TABLE inventory_transactions IS 'Audit trail for all inventory changes';
COMMENT ON COLUMN sizes.reserved IS 'Quantity reserved for pending orders (not available for sale)';
