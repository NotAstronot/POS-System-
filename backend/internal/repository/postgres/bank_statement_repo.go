package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BankStatement struct {
	ID              int64     `json:"id"`
	BankAccountID   int64     `json:"bank_account_id"`
	BankAccountName string    `json:"bank_account_name"`
	StatementDate   string    `json:"statement_date"`
	Description     string    `json:"description"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	Balance         float64   `json:"balance"`
	IsReconciled    bool      `json:"is_reconciled"`
	CreatedBy       *int64    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type BankStatementImportLine struct {
	ID              int64   `json:"id"`
	ImportID        int64   `json:"import_id"`
	TransactionDate string  `json:"transaction_date"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	IsReconciled    bool    `json:"is_reconciled"`
	ReconcileID     *int64  `json:"reconcile_id"`
}

type BankStatementImport struct {
	ID            int64                     `json:"id"`
	ImportNumber  string                    `json:"import_number"`
	AccountName   string                    `json:"account_name"`
	AccountType   string                    `json:"account_type"`
	StatementDate string                    `json:"statement_date"`
	Source        string                    `json:"source"`
	Status        string                    `json:"status"`
	TotalRows     int                       `json:"total_rows"`
	CreatedBy     *int64                    `json:"created_by"`
	CreatedAt     time.Time                 `json:"created_at"`
	Lines         []BankStatementImportLine `json:"lines"`
}

type BankStatementRepository struct {
	db *sql.DB
}

func NewBankStatementRepository(db *sql.DB) *BankStatementRepository {
	return &BankStatementRepository{db: db}
}

func (r *BankStatementRepository) ListAll(ctx context.Context) ([]BankStatement, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BankStatement, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT bs.id, bs.bank_account_id, a.account_name, bs.statement_date::text, bs.description, bs.debit, bs.credit, bs.balance, bs.is_reconciled, bs.created_by, bs.created_at
			FROM bank_statements bs
			JOIN chart_of_accounts a ON a.id = bs.bank_account_id AND a.tenant_id = bs.tenant_id
			WHERE bs.tenant_id=$1
			ORDER BY bs.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BankStatement, 0)
		for rows.Next() {
			var b BankStatement
			if err := rows.Scan(&b.ID, &b.BankAccountID, &b.BankAccountName, &b.StatementDate, &b.Description, &b.Debit, &b.Credit, &b.Balance, &b.IsReconciled, &b.CreatedBy, &b.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, b)
		}
		return list, nil
	})
}

func (r *BankStatementRepository) GetBankStatementByID(ctx context.Context, id int64) (*BankStatement, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*BankStatement, error) {
		b := &BankStatement{}
		err := tx.QueryRowContext(ctx, `
			SELECT bs.id, bs.bank_account_id, a.account_name, bs.statement_date::text, bs.description, bs.debit, bs.credit, bs.balance, bs.is_reconciled, bs.created_by, bs.created_at
			FROM bank_statements bs
			JOIN chart_of_accounts a ON a.id = bs.bank_account_id AND a.tenant_id = bs.tenant_id
			WHERE bs.id=$1 AND bs.tenant_id=$2`, id, tenantID).
			Scan(&b.ID, &b.BankAccountID, &b.BankAccountName, &b.StatementDate, &b.Description, &b.Debit, &b.Credit, &b.Balance, &b.IsReconciled, &b.CreatedBy, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		return b, nil
	})
}

func (r *BankStatementRepository) CreateBankStatement(ctx context.Context, b *BankStatement) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO bank_statements (bank_account_id, statement_date, description, debit, credit, balance, created_by, tenant_id) VALUES ($1,COALESCE(NULLIF($2,'')::date, CURRENT_DATE),$3,$4,$5,$6,$7,$8) RETURNING id",
			b.BankAccountID, b.StatementDate, b.Description, b.Debit, b.Credit, b.Balance, b.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *BankStatementRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM bank_statements WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankStatementRepository) ListImports(ctx context.Context) ([]BankStatementImport, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BankStatementImport, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, import_number, account_name, account_type, statement_date::text, source, status, total_rows, created_by, created_at
			FROM bank_statement_imports WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BankStatementImport, 0)
		for rows.Next() {
			var b BankStatementImport
			if err := rows.Scan(&b.ID, &b.ImportNumber, &b.AccountName, &b.AccountType, &b.StatementDate, &b.Source, &b.Status, &b.TotalRows, &b.CreatedBy, &b.CreatedAt); err != nil {
				return nil, err
			}
			lines, _ := r.GetImportLines(ctx, b.ID)
			b.Lines = lines
			if b.Lines == nil {
				b.Lines = make([]BankStatementImportLine, 0)
			}
			list = append(list, b)
		}
		return list, nil
	})
}

func (r *BankStatementRepository) GetImportLines(ctx context.Context, importID int64) ([]BankStatementImportLine, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BankStatementImportLine, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, import_id, transaction_date, description, amount, type, is_reconciled, reconcile_id FROM bank_statement_import_lines WHERE import_id=$1 AND tenant_id=$2 ORDER BY id", importID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BankStatementImportLine, 0)
		for rows.Next() {
			var l BankStatementImportLine
			if err := rows.Scan(&l.ID, &l.ImportID, &l.TransactionDate, &l.Description, &l.Amount, &l.Type, &l.IsReconciled, &l.ReconcileID); err != nil {
				return nil, err
			}
			list = append(list, l)
		}
		return list, nil
	})
}

func (r *BankStatementRepository) CreateImport(ctx context.Context, b *BankStatementImport) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO bank_statement_imports (import_number, account_name, account_type, statement_date, source, status, total_rows, created_by, tenant_id) VALUES ($1,$2,$3,COALESCE(NULLIF($4,'')::date, CURRENT_DATE),$5,$6,$7,$8,$9) RETURNING id",
			b.ImportNumber, b.AccountName, b.AccountType, b.StatementDate, b.Source, b.Status, b.TotalRows, b.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, l := range b.Lines {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO bank_statement_import_lines (import_id, transaction_date, description, amount, type, tenant_id) VALUES ($1,$2,$3,$4,$5,$6)",
				id, l.TransactionDate, l.Description, l.Amount, l.Type, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *BankStatementRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bank_statement_imports WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("BSI-%05d", count+1), nil
	})
}

func (r *BankStatementRepository) MarkLineReconciled(ctx context.Context, lineID int64, reconcileID int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE bank_statement_import_lines SET is_reconciled=true, reconcile_id=$1 WHERE id=$2 AND tenant_id=$3", reconcileID, lineID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankStatementRepository) Create(ctx context.Context, b *BankStatementImport) (int64, error) {
	return r.CreateImport(ctx, b)
}

func (r *BankStatementRepository) GetByID(ctx context.Context, id int64) (*BankStatementImport, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*BankStatementImport, error) {
		b := &BankStatementImport{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, import_number, account_name, account_type, statement_date::text, source, status, total_rows, created_by, created_at
			FROM bank_statement_imports WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&b.ID, &b.ImportNumber, &b.AccountName, &b.AccountType, &b.StatementDate, &b.Source, &b.Status, &b.TotalRows, &b.CreatedBy, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		lines, _ := r.GetImportLines(ctx, id)
		b.Lines = lines
		if b.Lines == nil {
			b.Lines = make([]BankStatementImportLine, 0)
		}
		return b, nil
	})
}
