package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Withdrawal struct {
	ID            int64     `json:"id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	RequestedBy   *int64    `json:"requested_by"`
	ApprovedBy    *int64    `json:"approved_by"`
	RequesterName string    `json:"requester_name"`
	ApproverName  string    `json:"approver_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WithdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) ListAll(ctx context.Context) ([]Withdrawal, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Withdrawal, error) {
		rows, err := tx.QueryContext(ctx, `SELECT w.id, w.amount, w.description, w.method, w.status, w.requested_by, w.approved_by,
			COALESCE(u1.name,'') as requester_name, COALESCE(u2.name,'') as approver_name,
			w.created_at, w.updated_at
			FROM withdrawals w LEFT JOIN users u1 ON w.requested_by=u1.id LEFT JOIN users u2 ON w.approved_by=u2.id WHERE w.tenant_id=$1 ORDER BY w.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Withdrawal
		for rows.Next() {
			var w Withdrawal
			if err := rows.Scan(&w.ID, &w.Amount, &w.Description, &w.Method, &w.Status, &w.RequestedBy, &w.ApprovedBy, &w.RequesterName, &w.ApproverName, &w.CreatedAt, &w.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, w)
		}
		return items, nil
	})
}

func (r *WithdrawalRepository) GetByID(ctx context.Context, id int64) (*Withdrawal, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Withdrawal, error) {
		w := &Withdrawal{}
		err := tx.QueryRowContext(ctx, `SELECT w.id, w.amount, w.description, w.method, w.status, w.requested_by, w.approved_by,
			COALESCE(u1.name,'') as requester_name, COALESCE(u2.name,'') as approver_name,
			w.created_at, w.updated_at
			FROM withdrawals w LEFT JOIN users u1 ON w.requested_by=u1.id LEFT JOIN users u2 ON w.approved_by=u2.id WHERE w.id=$1 AND w.tenant_id=$2`, id, tenantID).
			Scan(&w.ID, &w.Amount, &w.Description, &w.Method, &w.Status, &w.RequestedBy, &w.ApprovedBy, &w.RequesterName, &w.ApproverName, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return w, nil
	})
}

func (r *WithdrawalRepository) Create(ctx context.Context, w *Withdrawal) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO withdrawals (amount, description, method, requested_by, tenant_id) VALUES ($1,$2,$3,$4,$5) RETURNING id",
			w.Amount, w.Description, w.Method, w.RequestedBy, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *WithdrawalRepository) Update(ctx context.Context, w *Withdrawal) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE withdrawals SET amount=$1, description=$2, method=$3, status=$4, approved_by=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			w.Amount, w.Description, w.Method, w.Status, w.ApprovedBy, w.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WithdrawalRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM withdrawals WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WithdrawalRepository) Approve(ctx context.Context, id int64, approvedBy int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE withdrawals SET status='approved', approved_by=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", approvedBy, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WithdrawalRepository) Reject(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE withdrawals SET status='rejected', updated_at=NOW() WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
