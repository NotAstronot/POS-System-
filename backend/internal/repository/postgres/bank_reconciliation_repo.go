package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BankReconciliation struct {
	ID                   int64                    `json:"id"`
	ReconciliationNumber string                   `json:"reconciliation_number"`
	BankAccountID        int64                    `json:"bank_account_id"`
	BankAccountName      string                   `json:"bank_account_name"`
	AccountName          string                   `json:"account_name"`
	AccountType          string                   `json:"account_type"`
	StatementDate        string                   `json:"statement_date"`
	PeriodStart          string                   `json:"period_start"`
	PeriodEnd            string                   `json:"period_end"`
	StatementBalance     float64                  `json:"statement_balance"`
	BookBalance          float64                  `json:"book_balance"`
	SystemBalance        float64                  `json:"system_balance"`
	Difference           float64                  `json:"difference"`
	Status               string                   `json:"status"`
	Notes                string                   `json:"notes"`
	CreatedBy            *int64                   `json:"created_by"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
	Items                []BankReconciliationItem `json:"items"`
}

type BankReconciliationItem struct {
	ID               int64   `json:"id"`
	ReconciliationID int64   `json:"reconciliation_id"`
	Description      string  `json:"description"`
	TransactionDate  string  `json:"transaction_date"`
	Amount           float64 `json:"amount"`
	Type             string  `json:"type"`
	IsMatched        bool    `json:"is_matched"`
}

type BankReconciliationRepository struct {
	db *sql.DB
}

func NewBankReconciliationRepository(db *sql.DB) *BankReconciliationRepository {
	return &BankReconciliationRepository{db: db}
}

func (r *BankReconciliationRepository) ListAll(ctx context.Context) ([]BankReconciliation, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BankReconciliation, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT br.id, br.bank_account_id, a.account_name, br.statement_date::text, br.statement_balance, br.system_balance, br.status, br.notes, br.created_by, br.created_at, br.updated_at
			FROM bank_reconciliations br
			JOIN chart_of_accounts a ON a.id = br.bank_account_id AND a.tenant_id = br.tenant_id
			WHERE br.tenant_id=$1
			ORDER BY br.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BankReconciliation, 0)
		for rows.Next() {
			var b BankReconciliation
			if err := rows.Scan(&b.ID, &b.BankAccountID, &b.BankAccountName, &b.StatementDate, &b.StatementBalance, &b.SystemBalance, &b.Status, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, b)
		}
		return list, nil
	})
}

func (r *BankReconciliationRepository) GetByID(ctx context.Context, id int64) (*BankReconciliation, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*BankReconciliation, error) {
		b := &BankReconciliation{}
		err := tx.QueryRowContext(ctx, `
			SELECT br.id, br.bank_account_id, a.account_name, br.statement_date::text, br.statement_balance, br.system_balance, br.status, br.notes, br.created_by, br.created_at, br.updated_at
			FROM bank_reconciliations br
			JOIN chart_of_accounts a ON a.id = br.bank_account_id AND a.tenant_id = br.tenant_id
			WHERE br.id=$1 AND br.tenant_id=$2`, id, tenantID).
			Scan(&b.ID, &b.BankAccountID, &b.BankAccountName, &b.StatementDate, &b.StatementBalance, &b.SystemBalance, &b.Status, &b.Notes, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return b, nil
	})
}

func (r *BankReconciliationRepository) Create(ctx context.Context, b *BankReconciliation) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO bank_reconciliations (bank_account_id, statement_date, statement_balance, system_balance, notes, created_by, tenant_id) VALUES ($1,COALESCE(NULLIF($2,'')::date, CURRENT_DATE),$3,$4,$5,$6,$7) RETURNING id",
			b.BankAccountID, b.StatementDate, b.StatementBalance, b.SystemBalance, b.Notes, b.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *BankReconciliationRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE bank_reconciliations SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankReconciliationRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM bank_reconciliations WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankReconciliationRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bank_reconciliations WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("BR-%05d", count+1), nil
	})
}
