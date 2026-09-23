package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BankTransfer struct {
	ID              int64     `json:"id"`
	TransferNumber  string    `json:"transfer_number"`
	FromAccountName string    `json:"from_account_name"`
	FromAccountType string    `json:"from_account_type"`
	ToAccountName   string    `json:"to_account_name"`
	ToAccountType   string    `json:"to_account_type"`
	Amount          float64   `json:"amount"`
	TransferDate    string    `json:"transfer_date"`
	Notes           string    `json:"notes"`
	Status          string    `json:"status"`
	CreatedBy       *int64    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BankTransferRepository struct {
	db *sql.DB
}

func NewBankTransferRepository(db *sql.DB) *BankTransferRepository {
	return &BankTransferRepository{db: db}
}

func (r *BankTransferRepository) ListAll(ctx context.Context) ([]BankTransfer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BankTransfer, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, transfer_number, from_account_name, from_account_type,
				to_account_name, to_account_type, amount, transfer_date::text,
				notes, status, created_by, created_at, updated_at
			FROM bank_transfers WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BankTransfer, 0)
		for rows.Next() {
			var bt BankTransfer
			if err := rows.Scan(&bt.ID, &bt.TransferNumber, &bt.FromAccountName, &bt.FromAccountType,
				&bt.ToAccountName, &bt.ToAccountType, &bt.Amount, &bt.TransferDate,
				&bt.Notes, &bt.Status, &bt.CreatedBy, &bt.CreatedAt, &bt.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, bt)
		}
		return list, nil
	})
}

func (r *BankTransferRepository) GetByID(ctx context.Context, id int64) (*BankTransfer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*BankTransfer, error) {
		bt := &BankTransfer{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, transfer_number, from_account_name, from_account_type,
				to_account_name, to_account_type, amount, transfer_date::text,
				notes, status, created_by, created_at, updated_at
			FROM bank_transfers WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&bt.ID, &bt.TransferNumber, &bt.FromAccountName, &bt.FromAccountType,
				&bt.ToAccountName, &bt.ToAccountType, &bt.Amount, &bt.TransferDate,
				&bt.Notes, &bt.Status, &bt.CreatedBy, &bt.CreatedAt, &bt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return bt, nil
	})
}

func (r *BankTransferRepository) Create(ctx context.Context, bt *BankTransfer) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO bank_transfers (transfer_number, from_account_name, from_account_type, to_account_name, to_account_type, amount, transfer_date, notes, status, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,COALESCE(NULLIF($7,'')::date, CURRENT_DATE),$8,$9,$10,$11) RETURNING id",
			bt.TransferNumber, bt.FromAccountName, bt.FromAccountType,
			bt.ToAccountName, bt.ToAccountType, bt.Amount, bt.TransferDate,
			bt.Notes, bt.Status, bt.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *BankTransferRepository) Update(ctx context.Context, bt *BankTransfer) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE bank_transfers SET from_account_name=$1, from_account_type=$2, to_account_name=$3, to_account_type=$4, amount=$5, transfer_date=COALESCE(NULLIF($6,'')::date, CURRENT_DATE), notes=$7, status=$8, updated_at=NOW() WHERE id=$9 AND tenant_id=$10",
			bt.FromAccountName, bt.FromAccountType, bt.ToAccountName, bt.ToAccountType,
			bt.Amount, bt.TransferDate, bt.Notes, bt.Status, bt.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankTransferRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM bank_transfers WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BankTransferRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bank_transfers WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("TFB-%05d", count+1), nil
	})
}
