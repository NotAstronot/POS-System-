package postgres

import (
	"database/sql"
)

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderItemRepository struct {
	db *sql.DB
}

func NewOrderItemRepository(db *sql.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) ListByOrderID(orderID int64) ([]OrderItem, error) {
	rows, err := r.db.Query("SELECT id, order_id, product_id, quantity, price, subtotal FROM order_items WHERE order_id=$1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []OrderItem
	for rows.Next() {
		var oi OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.Price, &oi.Subtotal); err != nil {
			return nil, err
		}
		list = append(list, oi)
	}
	return list, nil
}
