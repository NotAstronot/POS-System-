package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StockTransferItem struct {
	ID          int64   `json:"id"`
	TransferID  int64   `json:"transfer_id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

type StockTransfer struct {
	ID                int64               `json:"id"`
	TransferNumber    string              `json:"transfer_number"`
	FromWarehouseID   int64               `json:"from_warehouse_id"`
	FromWarehouseName string              `json:"from_warehouse_name"`
	ToWarehouseID     int64               `json:"to_warehouse_id"`
	ToWarehouseName   string              `json:"to_warehouse_name"`
	TransferDate      string              `json:"transfer_date"`
	Status            string              `json:"status"`
	Notes             string              `json:"notes"`
	CreatedBy         *int64              `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Items             []StockTransferItem `json:"items"`
}

type StockTransferRepository struct {
	db *sql.DB
}

func NewStockTransferRepository(db *sql.DB) *StockTransferRepository {
	return &StockTransferRepository{db: db}
}

func (r *StockTransferRepository) ListAll(ctx context.Context) ([]StockTransfer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]StockTransfer, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT t.id, t.transfer_number, t.from_warehouse_id, fw.name, t.to_warehouse_id, tw.name,
				t.transfer_date::text, t.status, t.notes, t.created_by, t.created_at, t.updated_at
			FROM stock_transfers t
			JOIN warehouses fw ON fw.id = t.from_warehouse_id
			JOIN warehouses tw ON tw.id = t.to_warehouse_id
			WHERE t.tenant_id=$1
			ORDER BY t.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]StockTransfer, 0)
		for rows.Next() {
			var t StockTransfer
			if err := rows.Scan(&t.ID, &t.TransferNumber, &t.FromWarehouseID, &t.FromWarehouseName, &t.ToWarehouseID, &t.ToWarehouseName,
				&t.TransferDate, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, t.ID)
			t.Items = items
			if t.Items == nil {
				t.Items = make([]StockTransferItem, 0)
			}
			list = append(list, t)
		}
		return list, nil
	})
}

func (r *StockTransferRepository) GetByID(ctx context.Context, id int64) (*StockTransfer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*StockTransfer, error) {
		t := &StockTransfer{}
		err := tx.QueryRowContext(ctx, `
			SELECT t.id, t.transfer_number, t.from_warehouse_id, fw.name, t.to_warehouse_id, tw.name,
				t.transfer_date::text, t.status, t.notes, t.created_by, t.created_at, t.updated_at
			FROM stock_transfers t
			JOIN warehouses fw ON fw.id = t.from_warehouse_id
			JOIN warehouses tw ON tw.id = t.to_warehouse_id
			WHERE t.id=$1 AND t.tenant_id=$2`, id, tenantID).
			Scan(&t.ID, &t.TransferNumber, &t.FromWarehouseID, &t.FromWarehouseName, &t.ToWarehouseID, &t.ToWarehouseName,
				&t.TransferDate, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			t.Items = items
		}
		return t, nil
	})
}

func (r *StockTransferRepository) GetItems(ctx context.Context, transferID int64) ([]StockTransferItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]StockTransferItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, transfer_id, product_id, product_name, quantity, unit FROM stock_transfer_items WHERE transfer_id=$1 AND tenant_id=$2 ORDER BY id", transferID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]StockTransferItem, 0)
		for rows.Next() {
			var i StockTransferItem
			if err := rows.Scan(&i.ID, &i.TransferID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *StockTransferRepository) Create(ctx context.Context, t *StockTransfer) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO stock_transfers (transfer_number, from_warehouse_id, to_warehouse_id, transfer_date, notes, created_by, tenant_id) VALUES ($1,$2,$3,COALESCE(NULLIF($4,'')::date, CURRENT_DATE),$5,$6,$7) RETURNING id",
			t.TransferNumber, t.FromWarehouseID, t.ToWarehouseID, t.TransferDate, t.Notes, t.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range t.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO stock_transfer_items (transfer_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *StockTransferRepository) Update(ctx context.Context, t *StockTransfer) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx,
			"UPDATE stock_transfers SET from_warehouse_id=$1, to_warehouse_id=$2, transfer_date=COALESCE(NULLIF($3,'')::date, CURRENT_DATE), notes=$4, updated_at=NOW() WHERE id=$5 AND tenant_id=$6",
			t.FromWarehouseID, t.ToWarehouseID, t.TransferDate, t.Notes, t.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM stock_transfer_items WHERE transfer_id=$1 AND tenant_id=$2", t.ID, tenantID); err != nil {
			return err
		}
		for _, item := range t.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO stock_transfer_items (transfer_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				t.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StockTransferRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE stock_transfers SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *StockTransferRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM stock_transfers WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *StockTransferRepository) GenerateTransferNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_transfers WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("ST-%05d", count+1), nil
	})
}
