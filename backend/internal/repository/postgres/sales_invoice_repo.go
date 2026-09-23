package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SalesInvoice struct {
	ID            int64     `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	SalesOrderID  int64     `json:"sales_order_id"`
	OrderNumber   string    `json:"order_number"`
	CustomerID    int64     `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	InvoiceDate   string    `json:"invoice_date"`
	DueDate       *string   `json:"due_date"`
	TotalAmount   float64   `json:"total_amount"`
	PaidAmount    float64   `json:"paid_amount"`
	Remaining     float64   `json:"remaining_amount"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SalesInvoiceRepository struct {
	db *sql.DB
}

func NewSalesInvoiceRepository(db *sql.DB) *SalesInvoiceRepository {
	return &SalesInvoiceRepository{db: db}
}

func (r *SalesInvoiceRepository) ListAll(ctx context.Context) ([]SalesInvoice, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesInvoice, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT i.id, i.invoice_number, i.sales_order_id, o.order_number, i.customer_id, cu.name,
				i.invoice_date::text, i.due_date::text, i.total_amount, i.paid_amount, (i.total_amount - i.paid_amount),
				i.status, i.notes, i.created_at, i.updated_at
			FROM sales_invoices i
			JOIN sales_orders o ON o.id = i.sales_order_id
			JOIN customers cu ON cu.id = i.customer_id
			WHERE i.tenant_id=$1
			ORDER BY i.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]SalesInvoice, 0)
		for rows.Next() {
			var i SalesInvoice
			if err := rows.Scan(&i.ID, &i.InvoiceNumber, &i.SalesOrderID, &i.OrderNumber, &i.CustomerID, &i.CustomerName,
				&i.InvoiceDate, &i.DueDate, &i.TotalAmount, &i.PaidAmount, &i.Remaining, &i.Status, &i.Notes, &i.CreatedAt, &i.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, i)
		}
		return items, nil
	})
}

func (r *SalesInvoiceRepository) GetByID(ctx context.Context, id int64) (*SalesInvoice, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesInvoice, error) {
		i := &SalesInvoice{}
		err := tx.QueryRowContext(ctx, `
			SELECT i.id, i.invoice_number, i.sales_order_id, o.order_number, i.customer_id, cu.name,
				i.invoice_date::text, i.due_date::text, i.total_amount, i.paid_amount, (i.total_amount - i.paid_amount),
				i.status, i.notes, i.created_at, i.updated_at
			FROM sales_invoices i
			JOIN sales_orders o ON o.id = i.sales_order_id
			JOIN customers cu ON cu.id = i.customer_id
			WHERE i.id=$1 AND i.tenant_id=$2`, id, tenantID).
			Scan(&i.ID, &i.InvoiceNumber, &i.SalesOrderID, &i.OrderNumber, &i.CustomerID, &i.CustomerName,
				&i.InvoiceDate, &i.DueDate, &i.TotalAmount, &i.PaidAmount, &i.Remaining, &i.Status, &i.Notes, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return i, nil
	})
}

func (r *SalesInvoiceRepository) Create(ctx context.Context, i *SalesInvoice) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO sales_invoices (invoice_number, sales_order_id, customer_id, invoice_date, due_date, total_amount, notes, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
			i.InvoiceNumber, i.SalesOrderID, i.CustomerID, i.InvoiceDate, i.DueDate, i.TotalAmount, i.Notes, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *SalesInvoiceRepository) Update(ctx context.Context, i *SalesInvoice) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE sales_invoices SET sales_order_id=$1, customer_id=$2, invoice_date=$3, due_date=$4, total_amount=$5, status=$6, notes=$7, updated_at=NOW() WHERE id=$8 AND tenant_id=$9",
			i.SalesOrderID, i.CustomerID, i.InvoiceDate, i.DueDate, i.TotalAmount, i.Status, i.Notes, i.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesInvoiceRepository) UpdatePaidAmount(ctx context.Context, invoiceID int64, paidAmount float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_invoices SET paid_amount=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", paidAmount, invoiceID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesInvoiceRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_invoices WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesInvoiceRepository) GenerateInvoiceNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_invoices WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SI-%05d", count+1), nil
	})
}
