package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SalesReturnItem struct {
	ID            int64   `json:"id"`
	SalesReturnID int64   `json:"sales_return_id"`
	ProductID     int64   `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
	UnitPrice     float64 `json:"unit_price"`
	Subtotal      float64 `json:"subtotal"`
}

type SalesReturn struct {
	ID           int64             `json:"id"`
	ReturnNumber string            `json:"return_number"`
	SalesOrderID int64             `json:"sales_order_id"`
	CustomerID   int64             `json:"customer_id"`
	CustomerName string            `json:"customer_name"`
	ReturnDate   string            `json:"return_date"`
	Status       string            `json:"status"`
	Notes        string            `json:"notes"`
	TotalAmount  float64           `json:"total_amount"`
	CreatedBy    *int64            `json:"created_by"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Items        []SalesReturnItem `json:"items"`
}

type SalesReturnRepository struct {
	db *sql.DB
}

func NewSalesReturnRepository(db *sql.DB) *SalesReturnRepository {
	return &SalesReturnRepository{db: db}
}

func (r *SalesReturnRepository) ListAll(ctx context.Context) ([]SalesReturn, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesReturn, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT sr.id, sr.return_number, sr.customer_id, COALESCE(c.name,''), sr.return_date::text, sr.status, sr.notes, sr.created_by, sr.created_at, sr.updated_at
			FROM sales_returns sr
			LEFT JOIN customers c ON c.id = sr.customer_id AND c.tenant_id = sr.tenant_id
			WHERE sr.tenant_id=$1
			ORDER BY sr.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesReturn, 0)
		for rows.Next() {
			var sr SalesReturn
			if err := rows.Scan(&sr.ID, &sr.ReturnNumber, &sr.CustomerID, &sr.CustomerName, &sr.ReturnDate, &sr.Status, &sr.Notes, &sr.CreatedBy, &sr.CreatedAt, &sr.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, sr.ID)
			sr.Items = items
			if sr.Items == nil {
				sr.Items = make([]SalesReturnItem, 0)
			}
			list = append(list, sr)
		}
		return list, nil
	})
}

func (r *SalesReturnRepository) GetByID(ctx context.Context, id int64) (*SalesReturn, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesReturn, error) {
		sr := &SalesReturn{}
		err := tx.QueryRowContext(ctx, `
			SELECT sr.id, sr.return_number, sr.customer_id, COALESCE(c.name,''), sr.return_date::text, sr.status, sr.notes, sr.created_by, sr.created_at, sr.updated_at
			FROM sales_returns sr
			LEFT JOIN customers c ON c.id = sr.customer_id AND c.tenant_id = sr.tenant_id
			WHERE sr.id=$1 AND sr.tenant_id=$2`, id, tenantID).
			Scan(&sr.ID, &sr.ReturnNumber, &sr.CustomerID, &sr.CustomerName, &sr.ReturnDate, &sr.Status, &sr.Notes, &sr.CreatedBy, &sr.CreatedAt, &sr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			sr.Items = items
		}
		return sr, nil
	})
}

func (r *SalesReturnRepository) GetItems(ctx context.Context, returnID int64) ([]SalesReturnItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesReturnItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, sales_return_id, product_id, product_name, quantity, unit FROM sales_return_items WHERE sales_return_id=$1 AND tenant_id=$2 ORDER BY id", returnID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesReturnItem, 0)
		for rows.Next() {
			var i SalesReturnItem
			if err := rows.Scan(&i.ID, &i.SalesReturnID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *SalesReturnRepository) Create(ctx context.Context, sr *SalesReturn) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO sales_returns (return_number, customer_id, return_date, notes, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6) RETURNING id",
			sr.ReturnNumber, sr.CustomerID, sr.ReturnDate, sr.Notes, sr.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range sr.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO sales_return_items (sales_return_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *SalesReturnRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_returns SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesReturnRepository) Update(ctx context.Context, rt *SalesReturn) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		if _, err := tx.ExecContext(ctx,
			"UPDATE sales_returns SET sales_order_id=$1, customer_id=$2, return_date=COALESCE(NULLIF($3,'')::date, CURRENT_DATE), notes=$4, total_amount=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			rt.SalesOrderID, rt.CustomerID, rt.ReturnDate, rt.Notes, rt.TotalAmount, rt.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM sales_return_items WHERE sales_return_id=$1 AND tenant_id=$2", rt.ID, tenantID); err != nil {
			return err
		}
		for _, item := range rt.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO sales_return_items (sales_return_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				rt.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *SalesReturnRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_returns WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesReturnRepository) GenerateReturnNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_returns WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SR-%05d", count+1), nil
	})
}
