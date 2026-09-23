package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PurchaseInvoice struct {
	ID              int64     `json:"id"`
	InvoiceNumber   string    `json:"invoice_number"`
	PurchaseOrderID int64     `json:"purchase_order_id"`
	PONumber        string    `json:"po_number"`
	SupplierID      int64     `json:"supplier_id"`
	SupplierName    string    `json:"supplier_name"`
	InvoiceDate     string    `json:"invoice_date"`
	DueDate         *string   `json:"due_date"`
	TotalAmount     float64   `json:"total_amount"`
	PaidAmount      float64   `json:"paid_amount"`
	RemainingAmount float64   `json:"remaining_amount"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PurchaseInvoiceRepository struct {
	db *sql.DB
}

func NewPurchaseInvoiceRepository(db *sql.DB) *PurchaseInvoiceRepository {
	return &PurchaseInvoiceRepository{db: db}
}

func (r *PurchaseInvoiceRepository) ListAll(ctx context.Context) ([]PurchaseInvoice, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PurchaseInvoice, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT pi.id, pi.invoice_number, pi.purchase_order_id, po.po_number, pi.supplier_id, s.name,
				pi.invoice_date::text, pi.due_date::text, pi.total_amount, pi.paid_amount,
				(pi.total_amount - pi.paid_amount) as remaining_amount,
				pi.status, pi.notes, pi.created_at, pi.updated_at
			FROM purchase_invoices pi
			JOIN purchase_orders po ON po.id = pi.purchase_order_id
			JOIN suppliers s ON s.id = pi.supplier_id
			WHERE pi.tenant_id=$1
			ORDER BY pi.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]PurchaseInvoice, 0)
		for rows.Next() {
			var pi PurchaseInvoice
			if err := rows.Scan(&pi.ID, &pi.InvoiceNumber, &pi.PurchaseOrderID, &pi.PONumber, &pi.SupplierID, &pi.SupplierName,
				&pi.InvoiceDate, &pi.DueDate, &pi.TotalAmount, &pi.PaidAmount, &pi.RemainingAmount,
				&pi.Status, &pi.Notes, &pi.CreatedAt, &pi.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, pi)
		}
		return items, nil
	})
}

func (r *PurchaseInvoiceRepository) GetByID(ctx context.Context, id int64) (*PurchaseInvoice, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PurchaseInvoice, error) {
		pi := &PurchaseInvoice{}
		err := tx.QueryRowContext(ctx, `
			SELECT pi.id, pi.invoice_number, pi.purchase_order_id, po.po_number, pi.supplier_id, s.name,
				pi.invoice_date::text, pi.due_date::text, pi.total_amount, pi.paid_amount,
				(pi.total_amount - pi.paid_amount) as remaining_amount,
				pi.status, pi.notes, pi.created_at, pi.updated_at
			FROM purchase_invoices pi
			JOIN purchase_orders po ON po.id = pi.purchase_order_id
			JOIN suppliers s ON s.id = pi.supplier_id
			WHERE pi.id=$1 AND pi.tenant_id=$2`, id, tenantID).
			Scan(&pi.ID, &pi.InvoiceNumber, &pi.PurchaseOrderID, &pi.PONumber, &pi.SupplierID, &pi.SupplierName,
				&pi.InvoiceDate, &pi.DueDate, &pi.TotalAmount, &pi.PaidAmount, &pi.RemainingAmount,
				&pi.Status, &pi.Notes, &pi.CreatedAt, &pi.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return pi, nil
	})
}

func (r *PurchaseInvoiceRepository) Create(ctx context.Context, pi *PurchaseInvoice) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO purchase_invoices (invoice_number, purchase_order_id, supplier_id, invoice_date, due_date, total_amount, notes, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
			pi.InvoiceNumber, pi.PurchaseOrderID, pi.SupplierID, pi.InvoiceDate, pi.DueDate, pi.TotalAmount, pi.Notes, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *PurchaseInvoiceRepository) Update(ctx context.Context, pi *PurchaseInvoice) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE purchase_invoices SET invoice_date=$1, due_date=$2, total_amount=$3, paid_amount=$4, status=$5, notes=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8",
			pi.InvoiceDate, pi.DueDate, pi.TotalAmount, pi.PaidAmount, pi.Status, pi.Notes, pi.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseInvoiceRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM purchase_invoices WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseInvoiceRepository) UpdatePaidAmount(ctx context.Context, id int64, paidAmount float64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE purchase_invoices SET paid_amount=$1, status=$2, updated_at=NOW() WHERE id=$3 AND tenant_id=$4", paidAmount, status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PurchaseInvoiceRepository) GenerateInvoiceNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchase_invoices WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("INV-%05d", count+1), nil
	})
}
