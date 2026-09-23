package postgres

import (
	"context"
	"database/sql"
	"time"
)

type StockMovement struct {
	ID            int64     `json:"id"`
	ProductID     int64     `json:"product_id"`
	ProductName   string    `json:"product_name"`
	WarehouseID   int64     `json:"warehouse_id"`
	WarehouseName string    `json:"warehouse_name"`
	Quantity      float64   `json:"quantity"`
	MovementType  string    `json:"movement_type"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   *int64    `json:"reference_id"`
	Note          string    `json:"note"`
	CreatedBy     *int64    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) GetProductIDByName(ctx context.Context, name string) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "SELECT id FROM products WHERE name=$1 AND tenant_id=$2 ORDER BY id LIMIT 1", name, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *StockRepository) GetProductNameByID(ctx context.Context, id int64) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var name string
		err := tx.QueryRowContext(ctx, "SELECT name FROM products WHERE id=$1 AND tenant_id=$2", id, tenantID).Scan(&name)
		if err != nil {
			return "", err
		}
		return name, nil
	})
}

func (r *StockRepository) GetDefaultWarehouseID(ctx context.Context) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "SELECT id FROM warehouses WHERE is_active=true AND tenant_id=$1 ORDER BY id LIMIT 1", tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *StockRepository) GetStock(ctx context.Context, productID, warehouseID int64) (float64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (float64, error) {
		var qty float64
		err := tx.QueryRowContext(ctx, "SELECT COALESCE(quantity,0) FROM stock WHERE product_id=$1 AND warehouse_id=$2 AND tenant_id=$3", productID, warehouseID, tenantID).Scan(&qty)
		if err != nil {
			return 0, err
		}
		return qty, nil
	})
}

func (r *StockRepository) AddStock(ctx context.Context, productID, warehouseID int64, quantity float64, movementType, refType string, refID *int64, note string, createdBy *int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		_, err := tx.ExecContext(ctx, `
			INSERT INTO stock (product_id, warehouse_id, quantity, tenant_id) VALUES ($1,$2,$3,$4)
			ON CONFLICT (product_id, warehouse_id)
			DO UPDATE SET quantity = stock.quantity + $3, updated_at = NOW()`,
			productID, warehouseID, quantity, tenantID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO stock_movements (product_id, warehouse_id, quantity, movement_type, reference_type, reference_id, note, created_by, tenant_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			productID, warehouseID, quantity, movementType, refType, refID, note, createdBy, tenantID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE products SET stock = COALESCE((SELECT SUM(quantity) FROM stock WHERE product_id=$1 AND tenant_id=$2), 0) WHERE id=$1 AND tenant_id=$2`,
			productID, tenantID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *StockRepository) SetStock(ctx context.Context, productID, warehouseID int64, newQty float64, movementType, refType string, refID *int64, note string, createdBy *int64) error {
	cur, err := r.GetStock(ctx, productID, warehouseID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	diff := newQty - cur
	return r.AddStock(ctx, productID, warehouseID, diff, movementType, refType, refID, note, createdBy)
}

func (r *StockRepository) ListMovements(ctx context.Context) ([]StockMovement, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]StockMovement, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT m.id, m.product_id, COALESCE(p.name,''), m.warehouse_id, COALESCE(w.name,''),
				m.quantity, m.movement_type, m.reference_type, m.reference_id, m.note, m.created_by, m.created_at
			FROM stock_movements m
			LEFT JOIN products p ON p.id = m.product_id AND p.tenant_id = m.tenant_id
			LEFT JOIN warehouses w ON w.id = m.warehouse_id AND w.tenant_id = m.tenant_id
			WHERE m.tenant_id=$1
			ORDER BY m.id DESC LIMIT 500`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]StockMovement, 0)
		for rows.Next() {
			var m StockMovement
			if err := rows.Scan(&m.ID, &m.ProductID, &m.ProductName, &m.WarehouseID, &m.WarehouseName,
				&m.Quantity, &m.MovementType, &m.ReferenceType, &m.ReferenceID, &m.Note, &m.CreatedBy, &m.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, m)
		}
		return list, nil
	})
}
