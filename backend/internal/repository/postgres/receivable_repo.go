package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Receivable struct {
	ID          int64     `json:"id"`
	DebtorName  string    `json:"debtor_name"`
	Amount      float64   `json:"amount"`
	PaidAmount  float64   `json:"paid_amount"`
	Description string    `json:"description"`
	DueDate     *string   `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReceivableRepository struct {
	db *sql.DB
}

func NewReceivableRepository(db *sql.DB) *ReceivableRepository {
	return &ReceivableRepository{db: db}
}

func (r *ReceivableRepository) ListAll(ctx context.Context) ([]Receivable, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Receivable, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, debtor_name, amount, paid_amount, description, due_date::text, status, created_at, updated_at FROM receivables WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Receivable
		for rows.Next() {
			var rc Receivable
			if err := rows.Scan(&rc.ID, &rc.DebtorName, &rc.Amount, &rc.PaidAmount, &rc.Description, &rc.DueDate, &rc.Status, &rc.CreatedAt, &rc.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, rc)
		}
		return items, nil
	})
}

func (r *ReceivableRepository) GetByID(ctx context.Context, id int64) (*Receivable, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Receivable, error) {
		rc := &Receivable{}
		err := tx.QueryRowContext(ctx, "SELECT id, debtor_name, amount, paid_amount, description, due_date::text, status, created_at, updated_at FROM receivables WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&rc.ID, &rc.DebtorName, &rc.Amount, &rc.PaidAmount, &rc.Description, &rc.DueDate, &rc.Status, &rc.CreatedAt, &rc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return rc, nil
	})
}

func (r *ReceivableRepository) Create(ctx context.Context, rc *Receivable) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO receivables (debtor_name, amount, description, due_date, tenant_id) VALUES ($1,$2,$3,$4,$5) RETURNING id",
			rc.DebtorName, rc.Amount, rc.Description, rc.DueDate, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *ReceivableRepository) Update(ctx context.Context, rc *Receivable) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE receivables SET debtor_name=$1, amount=$2, paid_amount=$3, description=$4, due_date=$5, status=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8",
			rc.DebtorName, rc.Amount, rc.PaidAmount, rc.Description, rc.DueDate, rc.Status, rc.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *ReceivableRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM receivables WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
