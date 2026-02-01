-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255) NOT NULL,
    shoe_name VARCHAR(255) NOT NULL,
    base_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create countries table
CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    country_code VARCHAR(2) NOT NULL UNIQUE,
    currency VARCHAR(3) NOT NULL
);

-- Create product_variants table
CREATE TABLE IF NOT EXISTS product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    color VARCHAR(255) NOT NULL,
    shoe_size VARCHAR(50) NOT NULL,
    unique_identifier VARCHAR(500) NOT NULL,
    UNIQUE(product_id, color, shoe_size)
);

-- Create index on unique_identifier for faster lookups
CREATE INDEX IF NOT EXISTS idx_product_variants_unique_identifier ON product_variants(unique_identifier);
CREATE INDEX IF NOT EXISTS idx_product_variants_product_id ON product_variants(product_id);

-- Create price_history table
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    price NUMERIC(10, 2) NOT NULL,
    is_in_stock BOOLEAN NOT NULL DEFAULT false,
    captured_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_price_history_variant_id ON price_history(variant_id);
CREATE INDEX IF NOT EXISTS idx_price_history_country_id ON price_history(country_id);
CREATE INDEX IF NOT EXISTS idx_price_history_captured_at ON price_history(captured_at);
CREATE INDEX IF NOT EXISTS idx_price_history_variant_country ON price_history(variant_id, country_id);
