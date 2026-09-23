package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SalesQuotationItem struct {
	ID               int64   `json:"id"`
	SalesQuotationID int64   `json:"sales_quotation_id"`
	ProductID        int64   `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         float64 `json:"quantity"`
	Unit             string  `json:"unit"`
	UnitPrice        float64 `json:"unit_price"`
	Subtotal         float64 `json:"subtotal"`
}

type SalesQuotation struct {
	ID              int64                `json:"id"`
	QuotationNumber string               `json:"quotation_number"`
	CustomerID      int64                `json:"customer_id"`
	CustomerName    string               `json:"customer_name"`
	QuotationDate   string               `json:"quotation_date"`
	ValidUntil      string               `json:"valid_until"`
	Status          string               `json:"status"`
	Notes           string               `json:"notes"`
	TotalAmount     float64              `json:"total_amount"`
	CreatedBy       *int64               `json:"created_by"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Items           []SalesQuotationItem `json:"items"`
}

type SalesQuotationRepository struct {
	db *sql.DB
}

func NewSalesQuotationRepository(db *sql.DB) *SalesQuotationRepository {
	return &SalesQuotationRepository{db: db}
}

func (r *SalesQuotationRepository) ListAll(ctx context.Context) ([]SalesQuotation, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesQuotation, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT sq.id, sq.quotation_number, sq.customer_id, COALESCE(c.name,''), sq.quotation_date::text, sq.valid_until::text, sq.status, sq.notes, sq.total_amount, sq.created_by, sq.created_at, sq.updated_at
			FROM sales_quotations sq
			LEFT JOIN customers c ON c.id = sq.customer_id AND c.tenant_id = sq.tenant_id
			WHERE sq.tenant_id=$1
			ORDER BY sq.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesQuotation, 0)
		for rows.Next() {
			var sq SalesQuotation
			if err := rows.Scan(&sq.ID, &sq.QuotationNumber, &sq.CustomerID, &sq.CustomerName, &sq.QuotationDate, &sq.ValidUntil, &sq.Status, &sq.Notes, &sq.TotalAmount, &sq.CreatedBy, &sq.CreatedAt, &sq.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, sq.ID)
			sq.Items = items
			if sq.Items == nil {
				sq.Items = make([]SalesQuotationItem, 0)
			}
			list = append(list, sq)
		}
		return list, nil
	})
}

func (r *SalesQuotationRepository) GetByID(ctx context.Context, id int64) (*SalesQuotation, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesQuotation, error) {
		sq := &SalesQuotation{}
		err := tx.QueryRowContext(ctx, `
			SELECT sq.id, sq.quotation_number, sq.customer_id, COALESCE(c.name,''), sq.quotation_date::text, sq.valid_until::text, sq.status, sq.notes, sq.total_amount, sq.created_by, sq.created_at, sq.updated_at
			FROM sales_quotations sq
			LEFT JOIN customers c ON c.id = sq.customer_id AND c.tenant_id = sq.tenant_id
			WHERE sq.id=$1 AND sq.tenant_id=$2`, id, tenantID).
			Scan(&sq.ID, &sq.QuotationNumber, &sq.CustomerID, &sq.CustomerName, &sq.QuotationDate, &sq.ValidUntil, &sq.Status, &sq.Notes, &sq.TotalAmount, &sq.CreatedBy, &sq.CreatedAt, &sq.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			sq.Items = items
		}
		return sq, nil
	})
}

func (r *SalesQuotationRepository) GetItems(ctx context.Context, quotationID int64) ([]SalesQuotationItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesQuotationItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, sales_quotation_id, product_id, product_name, quantity, unit, unit_price FROM sales_quotation_items WHERE sales_quotation_id=$1 AND tenant_id=$2 ORDER BY id", quotationID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]SalesQuotationItem, 0)
		for rows.Next() {
			var i SalesQuotationItem
			if err := rows.Scan(&i.ID, &i.SalesQuotationID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit, &i.UnitPrice); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *SalesQuotationRepository) Create(ctx context.Context, sq *SalesQuotation) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO sales_quotations (quotation_number, customer_id, quotation_date, valid_until, notes, total_amount, created_by, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,'')::date, CURRENT_DATE),$4,$5,$6,$7,$8) RETURNING id",
			sq.QuotationNumber, sq.CustomerID, sq.QuotationDate, sq.ValidUntil, sq.Notes, sq.TotalAmount, sq.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range sq.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO sales_quotation_items (sales_quotation_id, product_id, product_name, quantity, unit, unit_price, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$7)),$4,$5,$6,$7)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, item.UnitPrice, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *SalesQuotationRepository) Update(ctx context.Context, sq *SalesQuotation) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx,
			"UPDATE sales_quotations SET customer_id=$1, quotation_date=COALESCE(NULLIF($2,'')::date, CURRENT_DATE), valid_until=$3, notes=$4, total_amount=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			sq.CustomerID, sq.QuotationDate, sq.ValidUntil, sq.Notes, sq.TotalAmount, sq.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM sales_quotation_items WHERE sales_quotation_id=$1 AND tenant_id=$2", sq.ID, tenantID); err != nil {
			return err
		}
		for _, item := range sq.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO sales_quotation_items (sales_quotation_id, product_id, product_name, quantity, unit, unit_price, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$7)),$4,$5,$6,$7)",
				sq.ID, item.ProductID, item.ProductName, item.Quantity, item.Unit, item.UnitPrice, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *SalesQuotationRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_quotations SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesQuotationRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_quotations WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesQuotationRepository) GenerateQuotationNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_quotations WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SQ-%05d", count+1), nil
	})
}
