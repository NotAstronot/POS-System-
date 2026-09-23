package postgres

import (
	"context"
	"database/sql"
	"time"
)

type ChartOfAccount struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	NormalBalance string    `json:"normal_balance"`
	Description   string    `json:"description"`
	IsActive      bool      `json:"is_active"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChartOfAccountRepository struct {
	db *sql.DB
}

func NewChartOfAccountRepository(db *sql.DB) *ChartOfAccountRepository {
	return &ChartOfAccountRepository{db: db}
}

func (r *ChartOfAccountRepository) ListAll(ctx context.Context) ([]ChartOfAccount, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ChartOfAccount, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT c.id, c.code, c.name, c.category, c.normal_balance, c.description, c.is_active,
				CASE WHEN c.normal_balance = 'kredit'
					THEN COALESCE(SUM(j.credit), 0) - COALESCE(SUM(j.debit), 0)
					ELSE COALESCE(SUM(j.debit), 0) - COALESCE(SUM(j.credit), 0)
				END AS balance,
				c.created_at, c.updated_at
			FROM chart_of_accounts c
			LEFT JOIN journal_entry_items j ON j.account_id = c.id AND j.tenant_id=$1
			WHERE c.tenant_id=$1
			GROUP BY c.id ORDER BY c.code`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]ChartOfAccount, 0)
		for rows.Next() {
			var c ChartOfAccount
			if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Category, &c.NormalBalance, &c.Description, &c.IsActive, &c.Balance, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, c)
		}
		return list, nil
	})
}

func (r *ChartOfAccountRepository) GetByID(ctx context.Context, id int64) (*ChartOfAccount, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*ChartOfAccount, error) {
		c := &ChartOfAccount{}
		err := tx.QueryRowContext(ctx, `
			SELECT c.id, c.code, c.name, c.category, c.normal_balance, c.description, c.is_active,
				CASE WHEN c.normal_balance = 'kredit'
					THEN COALESCE(SUM(j.credit), 0) - COALESCE(SUM(j.debit), 0)
					ELSE COALESCE(SUM(j.debit), 0) - COALESCE(SUM(j.credit), 0)
				END AS balance,
				c.created_at, c.updated_at
			FROM chart_of_accounts c
			LEFT JOIN journal_entry_items j ON j.account_id = c.id AND j.tenant_id=$2
			WHERE c.id=$1 AND c.tenant_id=$2
			GROUP BY c.id`, id, tenantID).
			Scan(&c.ID, &c.Code, &c.Name, &c.Category, &c.NormalBalance, &c.Description, &c.IsActive, &c.Balance, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *ChartOfAccountRepository) Create(ctx context.Context, c *ChartOfAccount) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO chart_of_accounts (code, name, category, normal_balance, description, is_active, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			c.Code, c.Name, c.Category, c.NormalBalance, c.Description, c.IsActive, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *ChartOfAccountRepository) Update(ctx context.Context, c *ChartOfAccount) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE chart_of_accounts SET code=$1, name=$2, category=$3, normal_balance=$4, description=$5, is_active=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8",
			c.Code, c.Name, c.Category, c.NormalBalance, c.Description, c.IsActive, c.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *ChartOfAccountRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM chart_of_accounts WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
