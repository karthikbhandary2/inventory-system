CREATE TABLE audit_logs (
  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  entity_type VARCHAR(50) NOT NULL,   -- 'product', 'stock_transaction'
  entity_id   UUID NOT NULL,
  action      VARCHAR(50) NOT NULL,   -- 'create', 'update', 'delete'
  old_values  JSONB,
  new_values  JSONB,
  performed_by VARCHAR(100) NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);