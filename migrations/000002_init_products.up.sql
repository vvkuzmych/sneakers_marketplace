-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(255) NOT NULL,
    colorway VARCHAR(100),
    sku VARCHAR(100) UNIQUE,
    retail_price DECIMAL(10, 2),
    release_date DATE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for products
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_model ON products(model);
CREATE INDEX idx_products_release_date ON products(release_date DESC);
CREATE INDEX idx_products_is_active ON products(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_products_search ON products USING GIN (to_tsvector('english', brand || ' ' || model || ' ' || COALESCE(colorway, '')));

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on product_id
CREATE INDEX idx_product_images_product_id ON product_images(product_id, position);

-- Create sizes table (available sizes for products)
CREATE TABLE IF NOT EXISTS sizes (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size VARCHAR(10) NOT NULL,
    stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    reserved_quantity INT NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, size)
);

-- Create indexes for sizes
CREATE INDEX idx_sizes_product_id ON sizes(product_id);
CREATE INDEX idx_sizes_stock ON sizes(product_id, size) WHERE stock_quantity > 0;

-- Create inventory_transactions table (for audit trail)
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id BIGSERIAL PRIMARY KEY,
    size_id BIGINT NOT NULL REFERENCES sizes(id) ON DELETE CASCADE,
    quantity_change INT NOT NULL,
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('purchase', 'sale', 'adjustment', 'return', 'reservation', 'release')),
    reference_id BIGINT, -- order_id or other reference
    reference_type VARCHAR(50), -- 'order', 'bid', 'ask', etc.
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on size_id and created_at for audit queries
CREATE INDEX idx_inventory_transactions_size_id ON inventory_transactions(size_id, created_at DESC);
CREATE INDEX idx_inventory_transactions_reference ON inventory_transactions(reference_type, reference_id);

-- Create triggers for updated_at
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sizes_updated_at BEFORE UPDATE ON sizes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to automatically create inventory transaction on size update
CREATE OR REPLACE FUNCTION log_inventory_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.stock_quantity != NEW.stock_quantity THEN
        INSERT INTO inventory_transactions (
            size_id,
            quantity_change,
            transaction_type,
            notes
        ) VALUES (
            NEW.id,
            NEW.stock_quantity - OLD.stock_quantity,
            'adjustment',
            'Automatic inventory adjustment'
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for automatic inventory logging
CREATE TRIGGER log_sizes_inventory_change AFTER UPDATE ON sizes
    FOR EACH ROW
    WHEN (OLD.stock_quantity IS DISTINCT FROM NEW.stock_quantity)
    EXECUTE FUNCTION log_inventory_change();

-- Comments for documentation
COMMENT ON TABLE products IS 'Sneaker products catalog';
COMMENT ON TABLE product_images IS 'Product images (multiple per product)';
COMMENT ON TABLE sizes IS 'Product sizes and inventory';
COMMENT ON TABLE inventory_transactions IS 'Audit trail for all inventory changes';
COMMENT ON COLUMN sizes.reserved_quantity IS 'Quantity reserved for pending orders (not available for sale)';
