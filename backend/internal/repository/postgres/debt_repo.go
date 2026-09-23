package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Debt struct {
	ID           int64     `json:"id"`
	CreditorName string    `json:"creditor_name"`
	Amount       float64   `json:"amount"`
	PaidAmount   float64   `json:"paid_amount"`
	Description  string    `json:"description"`
	DueDate      *string   `json:"due_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DebtRepository struct {
	db *sql.DB
}

func NewDebtRepository(db *sql.DB) *DebtRepository {
	return &DebtRepository{db: db}
}

func (r *DebtRepository) ListAll(ctx context.Context) ([]Debt, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Debt, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, creditor_name, amount, paid_amount, description, due_date::text, status, created_at, updated_at FROM debts WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Debt
		for rows.Next() {
			var d Debt
			if err := rows.Scan(&d.ID, &d.CreditorName, &d.Amount, &d.PaidAmount, &d.Description, &d.DueDate, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, d)
		}
		return items, nil
	})
}

func (r *DebtRepository) GetByID(ctx context.Context, id int64) (*Debt, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Debt, error) {
		d := &Debt{}
		err := tx.QueryRowContext(ctx, "SELECT id, creditor_name, amount, paid_amount, description, due_date::text, status, created_at, updated_at FROM debts WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&d.ID, &d.CreditorName, &d.Amount, &d.PaidAmount, &d.Description, &d.DueDate, &d.Status, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return d, nil
	})
}

func (r *DebtRepository) Create(ctx context.Context, d *Debt) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO debts (creditor_name, amount, description, due_date, tenant_id) VALUES ($1,$2,$3,$4,$5) RETURNING id",
			d.CreditorName, d.Amount, d.Description, d.DueDate, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *DebtRepository) Update(ctx context.Context, d *Debt) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE debts SET creditor_name=$1, amount=$2, paid_amount=$3, description=$4, due_date=$5, status=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8",
			d.CreditorName, d.Amount, d.PaidAmount, d.Description, d.DueDate, d.Status, d.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *DebtRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM debts WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
