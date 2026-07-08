package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/karthikbhandary2/inventory/internal/middleware"
	"github.com/karthikbhandary2/inventory/internal/models"
)

type ProductService struct{}

var ErrInsufficientStock = errors.New("insufficient stock")

func (s *ProductService) CreateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	tx := middleware.TxFromContext(ctx)

	row := tx.QueryRow(ctx, `
        INSERT INTO products (sku, name, description, quantity, price, low_stock_threshold)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `, p.SKU, p.Name, p.Description, p.Quantity, p.Price, p.LowStockThreshold)

	if err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}

	// Write audit log inside the same transaction
	if err := s.writeAudit(ctx, tx, "product", p.ID.String(), "create", nil, p); err != nil {
		return nil, err
	}

	return p, nil
}

// StockOperation applies a stock in/out operation atomically.
// Uses optimistic locking via FOR UPDATE to prevent race conditions.
func (s *ProductService) StockOperation(ctx context.Context, txn *models.StockTransaction) error {
	tx := middleware.TxFromContext(ctx)

	// Lock the row before reading — prevents concurrent updates
	var current models.Product
	err := tx.QueryRow(ctx,
		`SELECT id, quantity FROM products WHERE id = $1 FOR UPDATE`,
		txn.ProductID,
	).Scan(&current.ID, &current.Quantity)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("product not found: %w", pgx.ErrNoRows)
	}
	if err != nil {
		return fmt.Errorf("fetch product: %w", err)
	}

	// Calculate new quantity
	newQty := current.Quantity
	switch txn.Operation {
	case "in":
		newQty += txn.Quantity
	case "out":
		newQty -= txn.Quantity
		if newQty < 0 {
			return fmt.Errorf("%w: have %d, need %d", ErrInsufficientStock, current.Quantity, txn.Quantity)
		}
	case "adjustment":
		newQty = txn.Quantity
	}

	// Update quantity
	_, err = tx.Exec(ctx,
		`UPDATE products SET quantity = $1, updated_at = NOW() WHERE id = $2`,
		newQty, txn.ProductID,
	)
	if err != nil {
		return fmt.Errorf("update quantity: %w", err)
	}

	// Record the transaction
	_, err = tx.Exec(ctx, `
        INSERT INTO stock_transactions (product_id, operation, quantity, notes, performed_by)
        VALUES ($1, $2, $3, $4, $5)
    `, txn.ProductID, txn.Operation, txn.Quantity, txn.Notes, txn.PerformedBy)
	return err
}

func (s *ProductService) GetReport(ctx context.Context) (*models.InventoryReport, error) {
	tx := middleware.TxFromContext(ctx)

	report := &models.InventoryReport{}

	// Total value = SUM(quantity * price)
	err := tx.QueryRow(ctx, `
        SELECT COUNT(*), COALESCE(SUM(quantity * price), 0)
        FROM products
    `).Scan(&report.TotalProducts, &report.TotalValue)
	if err != nil {
		return nil, err
	}

	// Low stock items
	rows, err := tx.Query(ctx, `
        SELECT id, sku, name, description, quantity, price, low_stock_threshold, created_at, updated_at
        FROM products
        WHERE quantity <= low_stock_threshold
        ORDER BY quantity ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Quantity,
			&p.Price, &p.LowStockThreshold, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		report.LowStockItems = append(report.LowStockItems, p)
	}
	report.TotalLowStock = len(report.LowStockItems)
	return report, nil
}

func (s *ProductService) writeAudit(ctx context.Context, tx pgx.Tx, entityType, entityID, action string, old, new any) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO audit_logs (entity_type, entity_id, action, old_values, new_values, performed_by)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, entityType, entityID, action, old, new, "system")
	return err
}

func (s *ProductService) ListProducts(ctx context.Context, search string, lowStockOnly bool) ([]models.Product, error) {
	tx := middleware.TxFromContext(ctx)

	query := `SELECT id, sku, name, description, quantity, price, low_stock_threshold, created_at, updated_at
               FROM products WHERE 1=1`
	args := []any{}
	argN := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argN, argN)
		args = append(args, "%"+search+"%")
		argN++
	}
	if lowStockOnly {
		query += " AND quantity <= low_stock_threshold"
	}
	query += " ORDER BY name ASC"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Quantity,
			&p.Price, &p.LowStockThreshold, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	tx := middleware.TxFromContext(ctx)

	var p models.Product
	err := tx.QueryRow(ctx, `
        SELECT id, sku, name, description, quantity, price, low_stock_threshold, created_at, updated_at
        FROM products WHERE id = $1
    `, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Quantity,
		&p.Price, &p.LowStockThreshold, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err // bubbles up pgx.ErrNoRows for the handler to check
	}
	return &p, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	tx := middleware.TxFromContext(ctx)

	// Fetch old values first, for the audit log diff
	old, err := s.GetProduct(ctx, p.ID.String())
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, `
        UPDATE products
        SET name = $1, description = $2, price = $3, low_stock_threshold = $4, updated_at = NOW()
        WHERE id = $5
        RETURNING quantity, created_at, updated_at
    `, p.Name, p.Description, p.Price, p.LowStockThreshold, p.ID).
		Scan(&p.Quantity, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.SKU = old.SKU // SKU is immutable after creation — don't let it change here

	if err := s.writeAudit(ctx, tx, "product", p.ID.String(), "update", old, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	tx := middleware.TxFromContext(ctx)

	old, err := s.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	// ON DELETE RESTRICT on stock_transactions.product_id will fail this
	// if the product has transaction history — caught and surfaced as 409
	_, err = tx.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return s.writeAudit(ctx, tx, "product", id, "delete", old, nil)
}

func (s *ProductService) GetAuditLogs(ctx context.Context, entityID string) ([]models.AuditLog, error) {
	tx := middleware.TxFromContext(ctx)

	query := `SELECT id, entity_type, entity_id, action, old_values, new_values, performed_by, created_at
               FROM audit_logs`
	var args []any
	if entityID != "" {
		query += " WHERE entity_id = $1"
		args = append(args, entityID)
	}
	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action,
			&l.OldValues, &l.NewValues, &l.PerformedBy, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
