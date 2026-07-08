CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  sku         VARCHAR(50) UNIQUE NOT NULL,
  name        VARCHAR(255) NOT NULL,
  description TEXT,
  quantity    INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0), -- constraint: no negative stock
  price       NUMERIC(10,2) NOT NULL CHECK (price >= 0),
  low_stock_threshold INTEGER NOT NULL DEFAULT 10,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Index for search
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_name ON products(name);