package postgres

import (
	"context"
	"database/sql"
)

type OrderItem struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"order_id"`
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	VariantLabel string  `json:"variant_label"`
	Note         string  `json:"note"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
}

type OrderItemRepository struct {
	db *sql.DB
}

func NewOrderItemRepository(db *sql.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) ListByOrderID(ctx context.Context, orderID int64) ([]OrderItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]OrderItem, error) {
		rows, err := tx.QueryContext(ctx,
			`SELECT oi.id, oi.order_id, oi.product_id, COALESCE(p.name,''), COALESCE(oi.variant_label,''), COALESCE(oi.note,''),
				oi.quantity, oi.price, oi.subtotal
			 FROM order_items oi
			 LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id = oi.tenant_id
			 WHERE oi.order_id=$1 AND oi.tenant_id=$2
			 ORDER BY oi.id`, orderID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var list []OrderItem
		for rows.Next() {
			var oi OrderItem
			if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductName, &oi.VariantLabel, &oi.Note, &oi.Quantity, &oi.Price, &oi.Subtotal); err != nil {
				return nil, err
			}
			list = append(list, oi)
		}
		return list, rows.Err()
	})
}
