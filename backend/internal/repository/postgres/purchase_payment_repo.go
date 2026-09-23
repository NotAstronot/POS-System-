package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PurchasePayment struct {
	ID                int64     `json:"id"`
	PaymentNumber     string    `json:"payment_number"`
	PurchaseInvoiceID int64     `json:"purchase_invoice_id"`
	InvoiceNumber     string    `json:"invoice_number"`
	SupplierID        int64     `json:"supplier_id"`
	SupplierName      string    `json:"supplier_name"`
	PaymentDate       string    `json:"payment_date"`
	Amount            float64   `json:"amount"`
	PaymentMethod     string    `json:"payment_method"`
	Notes             string    `json:"notes"`
	CreatedBy         *int64    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
}

type PurchasePaymentRepository struct {
	db *sql.DB
}

func NewPurchasePaymentRepository(db *sql.DB) *PurchasePaymentRepository {
	return &PurchasePaymentRepository{db: db}
}

func (r *PurchasePaymentRepository) ListAll(ctx context.Context) ([]PurchasePayment, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchasePayment, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT pp.id, pp.payment_number, pp.purchase_invoice_id, pi.invoice_number, pp.supplier_id, s.name,
				pp.payment_date::text, pp.amount, pp.payment_method, pp.notes, pp.created_by, pp.created_at
			FROM purchase_payments pp
			JOIN purchase_invoices pi ON pi.id = pp.purchase_invoice_id
			JOIN suppliers s ON s.id = pp.supplier_id
			WHERE pp.tenant_id=$1
			ORDER BY pp.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]PurchasePayment, 0)
		for rows.Next() {
			var pp PurchasePayment
			if err := rows.Scan(&pp.ID, &pp.PaymentNumber, &pp.PurchaseInvoiceID, &pp.InvoiceNumber, &pp.SupplierID, &pp.SupplierName,
				&pp.PaymentDate, &pp.Amount, &pp.PaymentMethod, &pp.Notes, &pp.CreatedBy, &pp.CreatedAt); err != nil {
				return nil, err
			}
			items = append(items, pp)
		}
		return items, nil
	})
}

func (r *PurchasePaymentRepository) GetByID(ctx context.Context, id int64) (*PurchasePayment, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PurchasePayment, error) {
		pp := &PurchasePayment{}
		err := tx.QueryRowContext(ctx, `
			SELECT pp.id, pp.payment_number, pp.purchase_invoice_id, pi.invoice_number, pp.supplier_id, s.name,
				pp.payment_date::text, pp.amount, pp.payment_method, pp.notes, pp.created_by, pp.created_at
			FROM purchase_payments pp
			JOIN purchase_invoices pi ON pi.id = pp.purchase_invoice_id
			JOIN suppliers s ON s.id = pp.supplier_id
			WHERE pp.id=$1 AND pp.tenant_id=$2`, id, tenantID).
			Scan(&pp.ID, &pp.PaymentNumber, &pp.PurchaseInvoiceID, &pp.InvoiceNumber, &pp.SupplierID, &pp.SupplierName,
				&pp.PaymentDate, &pp.Amount, &pp.PaymentMethod, &pp.Notes, &pp.CreatedBy, &pp.CreatedAt)
		if err != nil {
			return nil, err
		}
		return pp, nil
	})
}

func (r *PurchasePaymentRepository) Create(ctx context.Context, pp *PurchasePayment) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO purchase_payments (payment_number, purchase_invoice_id, supplier_id, payment_date, amount, payment_method, notes, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id",
			pp.PaymentNumber, pp.PurchaseInvoiceID, pp.SupplierID, pp.PaymentDate, pp.Amount, pp.PaymentMethod, pp.Notes, pp.CreatedBy, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *PurchasePaymentRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM purchase_payments WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchasePaymentRepository) GeneratePaymentNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchase_payments WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("PAY-%05d", count+1), nil
	})
}
