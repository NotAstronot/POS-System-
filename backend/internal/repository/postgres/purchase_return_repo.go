package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PurchaseReturnItem struct {
	ID               int64   `json:"id"`
	PurchaseReturnID int64   `json:"purchase_return_id"`
	ProductID        int64   `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         float64 `json:"quantity"`
	Unit             string  `json:"unit"`
	UnitPrice        float64 `json:"unit_price"`
	Subtotal         float64 `json:"subtotal"`
	Reason           string  `json:"reason"`
}

type PurchaseReturn struct {
	ID              int64                `json:"id"`
	ReturnNumber    string               `json:"return_number"`
	PurchaseOrderID int64                `json:"purchase_order_id"`
	SupplierID      int64                `json:"supplier_id"`
	SupplierName    string               `json:"supplier_name"`
	ReturnDate      string               `json:"return_date"`
	Status          string               `json:"status"`
	Notes           string               `json:"notes"`
	Reason          string               `json:"reason"`
	CreatedBy       *int64               `json:"created_by"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Items           []PurchaseReturnItem `json:"items"`
}

type PurchaseReturnRepository struct {
	db *sql.DB
}

func NewPurchaseReturnRepository(db *sql.DB) *PurchaseReturnRepository {
	return &PurchaseReturnRepository{db: db}
}

func (r *PurchaseReturnRepository) ListAll(ctx context.Context) ([]PurchaseReturn, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchaseReturn, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT pr.id, pr.return_number, pr.supplier_id, COALESCE(s.name,''), pr.return_date::text, pr.status, pr.notes, pr.created_by, pr.created_at, pr.updated_at
			FROM purchase_returns pr
			LEFT JOIN suppliers s ON s.id = pr.supplier_id AND s.tenant_id = pr.tenant_id
			WHERE pr.tenant_id=$1
			ORDER BY pr.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PurchaseReturn, 0)
		for rows.Next() {
			var pr PurchaseReturn
			if err := rows.Scan(&pr.ID, &pr.ReturnNumber, &pr.SupplierID, &pr.SupplierName, &pr.ReturnDate, &pr.Status, &pr.Notes, &pr.CreatedBy, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, pr.ID)
			pr.Items = items
			if pr.Items == nil {
				pr.Items = make([]PurchaseReturnItem, 0)
			}
			list = append(list, pr)
		}
		return list, nil
	})
}

func (r *PurchaseReturnRepository) GetByID(ctx context.Context, id int64) (*PurchaseReturn, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PurchaseReturn, error) {
		pr := &PurchaseReturn{}
		err := tx.QueryRowContext(ctx, `
			SELECT pr.id, pr.return_number, pr.supplier_id, COALESCE(s.name,''), pr.return_date::text, pr.status, pr.notes, pr.created_by, pr.created_at, pr.updated_at
			FROM purchase_returns pr
			LEFT JOIN suppliers s ON s.id = pr.supplier_id AND s.tenant_id = pr.tenant_id
			WHERE pr.id=$1 AND pr.tenant_id=$2`, id, tenantID).
			Scan(&pr.ID, &pr.ReturnNumber, &pr.SupplierID, &pr.SupplierName, &pr.ReturnDate, &pr.Status, &pr.Notes, &pr.CreatedBy, &pr.CreatedAt, &pr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			pr.Items = items
		}
		return pr, nil
	})
}

func (r *PurchaseReturnRepository) GetItems(ctx context.Context, returnID int64) ([]PurchaseReturnItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchaseReturnItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, purchase_return_id, product_id, product_name, quantity, unit FROM purchase_return_items WHERE purchase_return_id=$1 AND tenant_id=$2 ORDER BY id", returnID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PurchaseReturnItem, 0)
		for rows.Next() {
			var i PurchaseReturnItem
			if err := rows.Scan(&i.ID, &i.PurchaseReturnID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *PurchaseReturnRepository) Create(ctx context.Context, pr *PurchaseReturn) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO purchase_returns (return_number, supplier_id, return_date, notes, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6) RETURNING id",
			pr.ReturnNumber, pr.SupplierID, pr.ReturnDate, pr.Notes, pr.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range pr.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO purchase_return_items (purchase_return_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *PurchaseReturnRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE purchase_returns SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseReturnRepository) Update(ctx context.Context, pr *PurchaseReturn) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		if _, err := tx.ExecContext(ctx,
			"UPDATE purchase_returns SET supplier_id=$1, return_date=COALESCE(NULLIF($2,'')::date, CURRENT_DATE), notes=$3, updated_at=NOW() WHERE id=$4 AND tenant_id=$5",
			pr.SupplierID, pr.ReturnDate, pr.Notes, pr.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM purchase_return_items WHERE purchase_return_id=$1 AND tenant_id=$2", pr.ID, tenantID); err != nil {
			return err
		}
		for _, item := range pr.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO purchase_return_items (purchase_return_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				pr.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PurchaseReturnRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM purchase_returns WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseReturnRepository) GenerateReturnNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchase_returns WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("PR-%05d", count+1), nil
	})
}
