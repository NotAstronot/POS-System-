package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type JournalEntryLine struct {
	ID             int64   `json:"id"`
	JournalEntryID int64   `json:"journal_entry_id"`
	AccountCode    string  `json:"account_code"`
	AccountName    string  `json:"account_name"`
	Debit          float64 `json:"debit"`
	Credit         float64 `json:"credit"`
	Description    string  `json:"description"`
}

type JournalEntry struct {
	ID            int64              `json:"id"`
	EntryNumber   string             `json:"entry_number"`
	EntryDate     string             `json:"entry_date"`
	ReferenceType string             `json:"reference_type"`
	ReferenceID   *int64             `json:"reference_id"`
	Reference     string             `json:"reference"`
	Description   string             `json:"description"`
	Status        string             `json:"status"`
	CreatedBy     *int64             `json:"created_by"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Lines         []JournalEntryLine `json:"lines"`
}

type JournalEntryItem struct {
	ID             int64   `json:"id"`
	JournalEntryID int64   `json:"journal_entry_id"`
	AccountID      int64   `json:"account_id"`
	AccountCode    string  `json:"account_code"`
	AccountName    string  `json:"account_name"`
	Debit          float64 `json:"debit"`
	Credit         float64 `json:"credit"`
	Description    string  `json:"description"`
}

type JournalEntryRepository struct {
	db *sql.DB
}

func NewJournalEntryRepository(db *sql.DB) *JournalEntryRepository {
	return &JournalEntryRepository{db: db}
}

func (r *JournalEntryRepository) ListAll(ctx context.Context) ([]JournalEntry, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]JournalEntry, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, entry_number, entry_date::text, reference_type, reference_id, description, status, created_by, created_at, updated_at
			FROM journal_entries WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]JournalEntry, 0)
		for rows.Next() {
			var je JournalEntry
			if err := rows.Scan(&je.ID, &je.EntryNumber, &je.EntryDate, &je.ReferenceType, &je.ReferenceID, &je.Description, &je.Status, &je.CreatedBy, &je.CreatedAt, &je.UpdatedAt); err != nil {
				return nil, err
			}
			lines, _ := r.GetLines(ctx, je.ID)
			je.Lines = lines
			if je.Lines == nil {
				je.Lines = make([]JournalEntryLine, 0)
			}
			list = append(list, je)
		}
		return list, nil
	})
}

func (r *JournalEntryRepository) GetByID(ctx context.Context, id int64) (*JournalEntry, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*JournalEntry, error) {
		je := &JournalEntry{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, entry_number, entry_date::text, reference_type, reference_id, description, status, created_by, created_at, updated_at
			FROM journal_entries WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&je.ID, &je.EntryNumber, &je.EntryDate, &je.ReferenceType, &je.ReferenceID, &je.Description, &je.Status, &je.CreatedBy, &je.CreatedAt, &je.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lines, err := r.GetLines(ctx, id)
		if err == nil {
			je.Lines = lines
		}
		return je, nil
	})
}

func (r *JournalEntryRepository) GetLines(ctx context.Context, entryID int64) ([]JournalEntryLine, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]JournalEntryLine, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT l.id, l.journal_entry_id, a.account_code, a.account_name, l.debit, l.credit, l.description
			FROM journal_entry_lines l
			JOIN chart_of_accounts a ON a.id = l.account_id AND a.tenant_id = l.tenant_id
			WHERE l.journal_entry_id=$1 AND l.tenant_id=$2 ORDER BY l.id`, entryID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]JournalEntryLine, 0)
		for rows.Next() {
			var l JournalEntryLine
			if err := rows.Scan(&l.ID, &l.JournalEntryID, &l.AccountCode, &l.AccountName, &l.Debit, &l.Credit, &l.Description); err != nil {
				return nil, err
			}
			list = append(list, l)
		}
		return list, nil
	})
}

func (r *JournalEntryRepository) Create(ctx context.Context, je *JournalEntry, items []JournalEntryItem) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO journal_entries (entry_number, entry_date, reference_type, reference_id, reference, description, created_by, tenant_id) VALUES ($1,COALESCE(NULLIF($2,'')::date, CURRENT_DATE),$3,$4,$5,$6,$7,$8) RETURNING id",
			je.EntryNumber, je.EntryDate, je.ReferenceType, je.ReferenceID, je.Reference, je.Description, je.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return err
		}
		for _, item := range items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO journal_entry_lines (journal_entry_id, account_id, debit, credit, description, tenant_id) VALUES ($1,$2,$3,$4,$5,$6)",
				id, item.AccountID, item.Debit, item.Credit, item.Description, tenantID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *JournalEntryRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("JE-%05d", count+1), nil
	})
}

func (r *JournalEntryRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE journal_entries SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *JournalEntryRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM journal_entries WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
