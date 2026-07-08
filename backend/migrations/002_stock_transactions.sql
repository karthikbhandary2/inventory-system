CREATE TYPE stock_operation AS ENUM ('in', 'out', 'adjustment');

CREATE TABLE stock_transactions (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id   UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  operation    stock_operation NOT NULL,
  quantity     INTEGER NOT NULL CHECK (quantity > 0),
  notes        TEXT,
  performed_by VARCHAR(100) NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stock_txns_product ON stock_transactions(product_id);
CREATE INDEX idx_stock_txns_created ON stock_transactions(created_at DESC);