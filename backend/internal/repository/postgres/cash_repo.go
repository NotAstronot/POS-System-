package postgres

import (
	"context"
	"database/sql"
	"time"
)

type CashTransaction struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Reference   string    `json:"reference"`
	CreatedBy   *int64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type CashTransactionRepository struct {
	db *sql.DB
}

func NewCashTransactionRepository(db *sql.DB) *CashTransactionRepository {
	return &CashTransactionRepository{db: db}
}

func (r *CashTransactionRepository) ListAll(ctx context.Context) ([]CashTransaction, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]CashTransaction, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, type, amount, description, category, reference, created_by, created_at FROM cash_transactions WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []CashTransaction
		for rows.Next() {
			var c CashTransaction
			if err := rows.Scan(&c.ID, &c.Type, &c.Amount, &c.Description, &c.Category, &c.Reference, &c.CreatedBy, &c.CreatedAt); err != nil {
				return nil, err
			}
			items = append(items, c)
		}
		return items, nil
	})
}

func (r *CashTransactionRepository) GetByID(ctx context.Context, id int64) (*CashTransaction, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*CashTransaction, error) {
		c := &CashTransaction{}
		err := tx.QueryRowContext(ctx, "SELECT id, type, amount, description, category, reference, created_by, created_at FROM cash_transactions WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&c.ID, &c.Type, &c.Amount, &c.Description, &c.Category, &c.Reference, &c.CreatedBy, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CashTransactionRepository) Create(ctx context.Context, c *CashTransaction) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO cash_transactions (type, amount, description, category, reference, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			c.Type, c.Amount, c.Description, c.Category, c.Reference, c.CreatedBy, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *CashTransactionRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM cash_transactions WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
