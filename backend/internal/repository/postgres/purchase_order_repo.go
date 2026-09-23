package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PurchaseOrderItem struct {
	ID              int64   `json:"id"`
	PurchaseOrderID int64   `json:"purchase_order_id"`
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Quantity        float64 `json:"quantity"`
	Unit            string  `json:"unit"`
	ReceivedQty     float64 `json:"received_qty"`
	UnitPrice       float64 `json:"unit_price"`
	Subtotal        float64 `json:"subtotal"`
}

type PurchaseItem struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	Quantity    float64 `json:"quantity"`
	ProductName string  `json:"product_name"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	ReceivedQty float64 `json:"received_qty"`
}

type PurchaseOrder struct {
	ID           int64          `json:"id"`
	OrderNumber  string         `json:"order_number"`
	PONumber     string         `json:"po_number"`
	SupplierID   int64          `json:"supplier_id"`
	SupplierName string         `json:"supplier_name"`
	OrderDate    string         `json:"order_date"`
	ExpectedDate *string        `json:"expected_date"`
	Status       string         `json:"status"`
	Notes        string         `json:"notes"`
	TotalAmount  float64        `json:"total_amount"`
	CreatedBy    *int64         `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Items        []PurchaseItem `json:"items"`
}

func (po *PurchaseOrder) GetExpectedDate() string {
	if po.ExpectedDate != nil {
		return *po.ExpectedDate
	}
	return ""
}

func (po *PurchaseOrder) SetExpectedDate(s *string) {
	po.ExpectedDate = s
}

type PurchaseOrderRepository struct {
	db *sql.DB
}

func NewPurchaseOrderRepository(db *sql.DB) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{db: db}
}

func (r *PurchaseOrderRepository) ListAll(ctx context.Context) ([]PurchaseOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchaseOrder, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT po.id, po.order_number, po.supplier_id, COALESCE(s.name,''), po.order_date::text, po.expected_date::text, po.status, po.notes, po.total_amount, po.created_by, po.created_at, po.updated_at
			FROM purchase_orders po
			LEFT JOIN suppliers s ON s.id = po.supplier_id AND s.tenant_id = po.tenant_id
			WHERE po.tenant_id=$1
			ORDER BY po.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PurchaseOrder, 0)
		for rows.Next() {
			var po PurchaseOrder
			if err := rows.Scan(&po.ID, &po.OrderNumber, &po.SupplierID, &po.SupplierName, &po.OrderDate, &po.ExpectedDate, &po.Status, &po.Notes, &po.TotalAmount, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, po.ID)
			po.Items = items
			if po.Items == nil {
				po.Items = make([]PurchaseItem, 0)
			}
			list = append(list, po)
		}
		return list, nil
	})
}

func (r *PurchaseOrderRepository) GetByID(ctx context.Context, id int64) (*PurchaseOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PurchaseOrder, error) {
		po := &PurchaseOrder{}
		err := tx.QueryRowContext(ctx, `
			SELECT po.id, po.order_number, po.supplier_id, COALESCE(s.name,''), po.order_date::text, po.expected_date::text, po.status, po.notes, po.total_amount, po.created_by, po.created_at, po.updated_at
			FROM purchase_orders po
			LEFT JOIN suppliers s ON s.id = po.supplier_id AND s.tenant_id = po.tenant_id
			WHERE po.id=$1 AND po.tenant_id=$2`, id, tenantID).
			Scan(&po.ID, &po.OrderNumber, &po.SupplierID, &po.SupplierName, &po.OrderDate, &po.ExpectedDate, &po.Status, &po.Notes, &po.TotalAmount, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			po.Items = items
		}
		return po, nil
	})
}

func (r *PurchaseOrderRepository) GetItems(ctx context.Context, orderID int64) ([]PurchaseItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchaseItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, product_id, product_name, quantity, unit, COALESCE((SELECT purchase_price FROM products WHERE id=product_id AND tenant_id=$2),0), COALESCE(received_qty,0) FROM purchase_order_items WHERE purchase_order_id=$1 AND tenant_id=$2 ORDER BY id", orderID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PurchaseItem, 0)
		for rows.Next() {
			var i PurchaseItem
			if err := rows.Scan(&i.ID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit, &i.UnitPrice, &i.ReceivedQty); err != nil {
				return nil, err
			}
			i.Subtotal = i.Quantity * i.UnitPrice
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *PurchaseOrderRepository) Create(ctx context.Context, po *PurchaseOrder) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO purchase_orders (order_number, supplier_id, order_date, expected_date, notes, total_amount, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6,$7,$8) RETURNING id",
			po.OrderNumber, po.SupplierID, po.OrderDate, po.ExpectedDate, po.Notes, po.TotalAmount, po.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range po.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO purchase_order_items (purchase_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *PurchaseOrderRepository) Update(ctx context.Context, po *PurchaseOrder) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx,
			"UPDATE purchase_orders SET supplier_id=$1, order_date=COALESCE(NULLIF($2,'')::date, CURRENT_DATE), expected_date=$3, notes=$4, total_amount=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			po.SupplierID, po.OrderDate, po.ExpectedDate, po.Notes, po.TotalAmount, po.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM purchase_order_items WHERE purchase_order_id=$1 AND tenant_id=$2", po.ID, tenantID); err != nil {
			return err
		}
		for _, item := range po.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO purchase_order_items (purchase_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				po.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PurchaseOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE purchase_orders SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseOrderRepository) UpdateItemReceivedQty(ctx context.Context, itemID int64, receivedQty float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE purchase_order_items SET received_qty=$1 WHERE id=$2 AND tenant_id=$3", receivedQty, itemID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseOrderRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM purchase_orders WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseOrderRepository) GeneratePONumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchase_orders WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("PO-%05d", count+1), nil
	})
}
