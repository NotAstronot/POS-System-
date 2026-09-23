package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SalesReceipt struct {
	ID             int64     `json:"id"`
	ReceiptNumber  string    `json:"receipt_number"`
	SalesInvoiceID int64     `json:"sales_invoice_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	CustomerID     int64     `json:"customer_id"`
	CustomerName   string    `json:"customer_name"`
	ReceiptDate    string    `json:"receipt_date"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	Notes          string    `json:"notes"`
	CreatedBy      *int64    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type SalesReceiptRepository struct {
	db *sql.DB
}

func NewSalesReceiptRepository(db *sql.DB) *SalesReceiptRepository {
	return &SalesReceiptRepository{db: db}
}

func (r *SalesReceiptRepository) ListAll(ctx context.Context) ([]SalesReceipt, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesReceipt, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT rc.id, rc.receipt_number, rc.sales_invoice_id, i.invoice_number, rc.customer_id, cu.name,
				rc.receipt_date::text, rc.amount, rc.payment_method, rc.notes, rc.created_by, rc.created_at
			FROM sales_receipts rc
			JOIN sales_invoices i ON i.id = rc.sales_invoice_id
			JOIN customers cu ON cu.id = rc.customer_id
			WHERE rc.tenant_id=$1
			ORDER BY rc.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]SalesReceipt, 0)
		for rows.Next() {
			var rc SalesReceipt
			if err := rows.Scan(&rc.ID, &rc.ReceiptNumber, &rc.SalesInvoiceID, &rc.InvoiceNumber, &rc.CustomerID, &rc.CustomerName,
				&rc.ReceiptDate, &rc.Amount, &rc.PaymentMethod, &rc.Notes, &rc.CreatedBy, &rc.CreatedAt); err != nil {
				return nil, err
			}
			items = append(items, rc)
		}
		return items, nil
	})
}

func (r *SalesReceiptRepository) GetByID(ctx context.Context, id int64) (*SalesReceipt, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesReceipt, error) {
		rc := &SalesReceipt{}
		err := tx.QueryRowContext(ctx, `
			SELECT rc.id, rc.receipt_number, rc.sales_invoice_id, i.invoice_number, rc.customer_id, cu.name,
				rc.receipt_date::text, rc.amount, rc.payment_method, rc.notes, rc.created_by, rc.created_at
			FROM sales_receipts rc
			JOIN sales_invoices i ON i.id = rc.sales_invoice_id
			JOIN customers cu ON cu.id = rc.customer_id
			WHERE rc.id=$1 AND rc.tenant_id=$2`, id, tenantID).
			Scan(&rc.ID, &rc.ReceiptNumber, &rc.SalesInvoiceID, &rc.InvoiceNumber, &rc.CustomerID, &rc.CustomerName,
				&rc.ReceiptDate, &rc.Amount, &rc.PaymentMethod, &rc.Notes, &rc.CreatedBy, &rc.CreatedAt)
		if err != nil {
			return nil, err
		}
		return rc, nil
	})
}

func (r *SalesReceiptRepository) Create(ctx context.Context, rc *SalesReceipt) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO sales_receipts (receipt_number, sales_invoice_id, customer_id, receipt_date, amount, payment_method, notes, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id",
			rc.ReceiptNumber, rc.SalesInvoiceID, rc.CustomerID, rc.ReceiptDate, rc.Amount, rc.PaymentMethod, rc.Notes, rc.CreatedBy, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *SalesReceiptRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_receipts WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesReceiptRepository) GenerateReceiptNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_receipts WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SR-%05d", count+1), nil
	})
}
