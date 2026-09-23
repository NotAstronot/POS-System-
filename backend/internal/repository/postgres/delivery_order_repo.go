package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DeliveryOrderItem struct {
	ID              int64   `json:"id"`
	DeliveryOrderID int64   `json:"delivery_order_id"`
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Quantity        float64 `json:"quantity"`
	Unit            string  `json:"unit"`
	DeliveredQty    float64 `json:"delivered_qty"`
}

type DeliveryOrder struct {
	ID           int64               `json:"id"`
	OrderNumber  string              `json:"order_number"`
	DONumber     string              `json:"do_number"`
	SalesOrderID int64               `json:"sales_order_id"`
	CustomerID   int64               `json:"customer_id"`
	CustomerName string              `json:"customer_name"`
	DeliveryDate string              `json:"delivery_date"`
	Status       string              `json:"status"`
	Notes        string              `json:"notes"`
	CreatedBy    *int64              `json:"created_by"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	Items        []DeliveryOrderItem `json:"items"`
}

type DeliveryOrderRepository struct {
	db *sql.DB
}

func NewDeliveryOrderRepository(db *sql.DB) *DeliveryOrderRepository {
	return &DeliveryOrderRepository{db: db}
}

func (r *DeliveryOrderRepository) ListAll(ctx context.Context) ([]DeliveryOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]DeliveryOrder, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT do.id, do.order_number, COALESCE(do.do_number,''), do.sales_order_id, do.customer_id, COALESCE(c.name,''), do.delivery_date::text, do.status, do.notes, do.created_by, do.created_at, do.updated_at
			FROM delivery_orders do
			LEFT JOIN customers c ON c.id = do.customer_id AND c.tenant_id = do.tenant_id
			WHERE do.tenant_id=$1
			ORDER BY do.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]DeliveryOrder, 0)
		for rows.Next() {
			var d DeliveryOrder
			if err := rows.Scan(&d.ID, &d.OrderNumber, &d.DONumber, &d.SalesOrderID, &d.CustomerID, &d.CustomerName, &d.DeliveryDate, &d.Status, &d.Notes, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, d.ID)
			d.Items = items
			if d.Items == nil {
				d.Items = make([]DeliveryOrderItem, 0)
			}
			list = append(list, d)
		}
		return list, nil
	})
}

func (r *DeliveryOrderRepository) GetByID(ctx context.Context, id int64) (*DeliveryOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*DeliveryOrder, error) {
		d := &DeliveryOrder{}
		err := tx.QueryRowContext(ctx, `
			SELECT do.id, do.order_number, COALESCE(do.do_number,''), do.sales_order_id, do.customer_id, COALESCE(c.name,''), do.delivery_date::text, do.status, do.notes, do.created_by, do.created_at, do.updated_at
			FROM delivery_orders do
			LEFT JOIN customers c ON c.id = do.customer_id AND c.tenant_id = do.tenant_id
			WHERE do.id=$1 AND do.tenant_id=$2`, id, tenantID).
			Scan(&d.ID, &d.OrderNumber, &d.DONumber, &d.SalesOrderID, &d.CustomerID, &d.CustomerName, &d.DeliveryDate, &d.Status, &d.Notes, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			d.Items = items
		}
		return d, nil
	})
}

func (r *DeliveryOrderRepository) GetItems(ctx context.Context, orderID int64) ([]DeliveryOrderItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]DeliveryOrderItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, delivery_order_id, product_id, product_name, quantity, unit, delivered_qty FROM delivery_order_items WHERE delivery_order_id=$1 AND tenant_id=$2 ORDER BY id", orderID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]DeliveryOrderItem, 0)
		for rows.Next() {
			var i DeliveryOrderItem
			if err := rows.Scan(&i.ID, &i.DeliveryOrderID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit, &i.DeliveredQty); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *DeliveryOrderRepository) Create(ctx context.Context, d *DeliveryOrder) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO delivery_orders (order_number, sales_order_id, customer_id, delivery_date, notes, created_by, tenant_id) VALUES ($1,$2,$3,COALESCE(NULLIF($4,'')::date, CURRENT_DATE),$5,$6,$7) RETURNING id",
			d.OrderNumber, d.SalesOrderID, d.CustomerID, d.DeliveryDate, d.Notes, d.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range d.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO delivery_order_items (delivery_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *DeliveryOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE delivery_orders SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *DeliveryOrderRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM delivery_orders WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *DeliveryOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM delivery_orders WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("DO-%05d", count+1), nil
	})
}

func (r *DeliveryOrderRepository) GenerateDONumber(ctx context.Context) (string, error) {
	return r.GenerateOrderNumber(ctx)
}

func (r *DeliveryOrderRepository) Update(ctx context.Context, d *DeliveryOrder) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		if _, err := tx.ExecContext(ctx,
			"UPDATE delivery_orders SET do_number=$1, sales_order_id=$2, customer_id=$3, delivery_date=COALESCE(NULLIF($4,'')::date, CURRENT_DATE), status=$5, notes=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8",
			d.DONumber, d.SalesOrderID, d.CustomerID, d.DeliveryDate, d.Status, d.Notes, d.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM delivery_order_items WHERE delivery_order_id=$1 AND tenant_id=$2", d.ID, tenantID); err != nil {
			return err
		}
		for _, item := range d.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO delivery_order_items (delivery_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				d.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}
