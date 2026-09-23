package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SalesOrderItem struct {
	ID           int64   `json:"id"`
	SalesOrderID int64   `json:"sales_order_id"`
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	FulfilledQty float64 `json:"fulfilled_qty"`
	DeliveredQty float64 `json:"delivered_qty"`
	UnitPrice    float64 `json:"unit_price"`
	Subtotal     float64 `json:"subtotal"`
}

type SalesOrder struct {
	ID              int64            `json:"id"`
	OrderNumber     string           `json:"order_number"`
	CustomerID      int64            `json:"customer_id"`
	CustomerName    string           `json:"customer_name"`
	SalesCategoryID *int64           `json:"sales_category_id"`
	QuotationID     *int64           `json:"quotation_id"`
	OrderDate       string           `json:"order_date"`
	DueDate         string           `json:"due_date"`
	Status          string           `json:"status"`
	Notes           string           `json:"notes"`
	TotalAmount     float64          `json:"total_amount"`
	CreatedBy       *int64           `json:"created_by"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Items           []SalesOrderItem `json:"items"`
}

type SalesOrderRepository struct {
	db *sql.DB
}

func NewSalesOrderRepository(db *sql.DB) *SalesOrderRepository {
	return &SalesOrderRepository{db: db}
}

func (r *SalesOrderRepository) ListAll(ctx context.Context) ([]SalesOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesOrder, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT so.id, so.order_number, so.customer_id, COALESCE(c.name,''), so.order_date::text, so.due_date::text, so.status, so.notes, so.total_amount, so.created_by, so.created_at, so.updated_at
			FROM sales_orders so
			LEFT JOIN customers c ON c.id = so.customer_id AND c.tenant_id = so.tenant_id
			WHERE so.tenant_id=$1
			ORDER BY so.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesOrder, 0)
		for rows.Next() {
			var so SalesOrder
			if err := rows.Scan(&so.ID, &so.OrderNumber, &so.CustomerID, &so.CustomerName, &so.OrderDate, &so.DueDate, &so.Status, &so.Notes, &so.TotalAmount, &so.CreatedBy, &so.CreatedAt, &so.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, so.ID)
			so.Items = items
			if so.Items == nil {
				so.Items = make([]SalesOrderItem, 0)
			}
			list = append(list, so)
		}
		return list, nil
	})
}

func (r *SalesOrderRepository) GetByID(ctx context.Context, id int64) (*SalesOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesOrder, error) {
		so := &SalesOrder{}
		err := tx.QueryRowContext(ctx, `
			SELECT so.id, so.order_number, so.customer_id, COALESCE(c.name,''), so.order_date::text, so.due_date::text, so.status, so.notes, so.total_amount, so.created_by, so.created_at, so.updated_at
			FROM sales_orders so
			LEFT JOIN customers c ON c.id = so.customer_id AND c.tenant_id = so.tenant_id
			WHERE so.id=$1 AND so.tenant_id=$2`, id, tenantID).
			Scan(&so.ID, &so.OrderNumber, &so.CustomerID, &so.CustomerName, &so.OrderDate, &so.DueDate, &so.Status, &so.Notes, &so.TotalAmount, &so.CreatedBy, &so.CreatedAt, &so.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			so.Items = items
		}
		return so, nil
	})
}

func (r *SalesOrderRepository) GetItems(ctx context.Context, orderID int64) ([]SalesOrderItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesOrderItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, sales_order_id, product_id, product_name, quantity, unit, fulfilled_qty FROM sales_order_items WHERE sales_order_id=$1 AND tenant_id=$2 ORDER BY id", orderID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesOrderItem, 0)
		for rows.Next() {
			var i SalesOrderItem
			if err := rows.Scan(&i.ID, &i.SalesOrderID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit, &i.FulfilledQty); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *SalesOrderRepository) Create(ctx context.Context, so *SalesOrder) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO sales_orders (order_number, customer_id, order_date, due_date, notes, total_amount, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6,$7,$8) RETURNING id",
			so.OrderNumber, so.CustomerID, so.OrderDate, so.DueDate, so.Notes, so.TotalAmount, so.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range so.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO sales_order_items (sales_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *SalesOrderRepository) Update(ctx context.Context, so *SalesOrder) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx,
			"UPDATE sales_orders SET customer_id=$1, order_date=COALESCE(NULLIF($2,'')::date, CURRENT_DATE), due_date=$3, notes=$4, total_amount=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			so.CustomerID, so.OrderDate, so.DueDate, so.Notes, so.TotalAmount, so.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM sales_order_items WHERE sales_order_id=$1 AND tenant_id=$2", so.ID, tenantID); err != nil {
			return err
		}
		for _, item := range so.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO sales_order_items (sales_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				so.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *SalesOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_orders SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesOrderRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_orders WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_orders WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SO-%05d", count+1), nil
	})
}

func (r *SalesOrderRepository) UpdateItemDeliveredQty(ctx context.Context, itemID int64, deliveredQty float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_order_items SET fulfilled_qty=$1 WHERE id=$2 AND tenant_id=$3", deliveredQty, itemID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
