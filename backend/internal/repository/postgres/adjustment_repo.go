package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StockAdjustmentItem struct {
	ID           int64   `json:"id"`
	AdjustmentID int64   `json:"adjustment_id"`
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Quantity     float64 `json:"quantity"`
	ActualQty    float64 `json:"actual_qty"`
	SystemQty    float64 `json:"system_qty"`
	Difference   float64 `json:"difference"`
	Reason       string  `json:"reason"`
	Unit         string  `json:"unit"`
}

type StockAdjustment struct {
	ID               int64                 `json:"id"`
	AdjustmentNumber string                `json:"adjustment_number"`
	AdjustmentNo     string                `json:"adjustment_no"`
	WarehouseID      int64                 `json:"warehouse_id"`
	WarehouseName    string                `json:"warehouse_name"`
	AdjustmentDate   string                `json:"adjustment_date"`
	Type             string                `json:"type"`
	Reason           string                `json:"reason"`
	Status           string                `json:"status"`
	Notes            string                `json:"notes"`
	CreatedBy        *int64                `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Items            []StockAdjustmentItem `json:"items"`
}

type StockAdjustmentRepository struct {
	db *sql.DB
}

func NewStockAdjustmentRepository(db *sql.DB) *StockAdjustmentRepository {
	return &StockAdjustmentRepository{db: db}
}

func (r *StockAdjustmentRepository) ListAll(ctx context.Context) ([]StockAdjustment, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]StockAdjustment, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT a.id, a.adjustment_no, a.warehouse_id, w.name, a.adjustment_date::text, a.type, a.reason, a.status, a.notes, a.created_by, a.created_at, a.updated_at
			FROM stock_adjustments a
			JOIN warehouses w ON w.id = a.warehouse_id
			WHERE a.tenant_id=$1
			ORDER BY a.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]StockAdjustment, 0)
		for rows.Next() {
			var a StockAdjustment
			if err := rows.Scan(&a.ID, &a.AdjustmentNo, &a.WarehouseID, &a.WarehouseName, &a.AdjustmentDate, &a.Type, &a.Reason, &a.Status, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, a.ID)
			a.Items = items
			if a.Items == nil {
				a.Items = make([]StockAdjustmentItem, 0)
			}
			list = append(list, a)
		}
		return list, nil
	})
}

func (r *StockAdjustmentRepository) GetByID(ctx context.Context, id int64) (*StockAdjustment, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*StockAdjustment, error) {
		a := &StockAdjustment{}
		err := tx.QueryRowContext(ctx, `
			SELECT a.id, a.adjustment_no, a.warehouse_id, w.name, a.adjustment_date::text, a.type, a.reason, a.status, a.notes, a.created_by, a.created_at, a.updated_at
			FROM stock_adjustments a
			JOIN warehouses w ON w.id = a.warehouse_id
			WHERE a.id=$1 AND a.tenant_id=$2`, id, tenantID).
			Scan(&a.ID, &a.AdjustmentNo, &a.WarehouseID, &a.WarehouseName, &a.AdjustmentDate, &a.Type, &a.Reason, &a.Status, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			a.Items = items
		}
		return a, nil
	})
}

func (r *StockAdjustmentRepository) GetItems(ctx context.Context, adjustmentID int64) ([]StockAdjustmentItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]StockAdjustmentItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, adjustment_id, product_id, product_name, quantity, unit FROM stock_adjustment_items WHERE adjustment_id=$1 AND tenant_id=$2 ORDER BY id", adjustmentID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]StockAdjustmentItem, 0)
		for rows.Next() {
			var i StockAdjustmentItem
			if err := rows.Scan(&i.ID, &i.AdjustmentID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *StockAdjustmentRepository) Create(ctx context.Context, a *StockAdjustment) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO stock_adjustments (adjustment_no, warehouse_id, adjustment_date, type, reason, notes, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6,$7,$8) RETURNING id",
			a.AdjustmentNo, a.WarehouseID, a.AdjustmentDate, a.Type, a.Reason, a.Notes, a.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range a.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO stock_adjustment_items (adjustment_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *StockAdjustmentRepository) Update(ctx context.Context, a *StockAdjustment) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx,
			"UPDATE stock_adjustments SET warehouse_id=$1, adjustment_date=COALESCE(NULLIF($2,'')::date, CURRENT_DATE), type=$3, reason=$4, notes=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			a.WarehouseID, a.AdjustmentDate, a.Type, a.Reason, a.Notes, a.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM stock_adjustment_items WHERE adjustment_id=$1 AND tenant_id=$2", a.ID, tenantID); err != nil {
			return err
		}
		for _, item := range a.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO stock_adjustment_items (adjustment_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				a.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StockAdjustmentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE stock_adjustments SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *StockAdjustmentRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM stock_adjustments WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *StockAdjustmentRepository) GenerateAdjustmentNumber(ctx context.Context) (string, error) {
	return r.GenerateAdjustmentNo(ctx)
}

func (r *StockAdjustmentRepository) GenerateAdjustmentNo(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_adjustments WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("ADJ-%05d", count+1), nil
	})
}
