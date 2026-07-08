package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID                uuid.UUID `json:"id" db:"id"`
	SKU               string    `json:"sku" db:"sku" validate:"required,min=3,max=50"`
	Name              string    `json:"name" db:"name" validate:"required,min=1,max=255"`
	Description       string    `json:"description" db:"description"`
	Quantity          int       `json:"quantity" db:"quantity"`
	Price             float64   `json:"price" db:"price" validate:"required,gte=0"`
	LowStockThreshold int       `json:"low_stock_threshold" db:"low_stock_threshold"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type StockTransaction struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	Operation   string    `json:"operation" validate:"required,oneof=in out adjustment"`
	Quantity    int       `json:"quantity" validate:"required,gt=0"`
	Notes       string    `json:"notes"`
	PerformedBy string    `json:"performed_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuditLog struct {
	ID          uuid.UUID `json:"id"`
	EntityType  string    `json:"entity_type"`
	EntityID    uuid.UUID `json:"entity_id"`
	Action      string    `json:"action"`
	OldValues   any       `json:"old_values"`
	NewValues   any       `json:"new_values"`
	PerformedBy string    `json:"performed_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// Report types
type InventoryReport struct {
	TotalProducts int       `json:"total_products"`
	TotalValue    float64   `json:"total_value"`
	LowStockItems []Product `json:"low_stock_items"`
	TotalLowStock int       `json:"total_low_stock"`
}
