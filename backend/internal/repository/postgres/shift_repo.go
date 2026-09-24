package postgres

import (
	"context"
	"database/sql"

	"time"
)

type Shift struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	OpenedAt       time.Time  `json:"opened_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	OpeningBalance float64    `json:"opening_balance"`
	ClosingBalance *float64   `json:"closing_balance"`
	Status         string     `json:"status"`
}

type ShiftRepository struct {
	db *sql.DB
}

func NewShiftRepository(db *sql.DB) *ShiftRepository {
	return &ShiftRepository{db: db}
}

func (r *ShiftRepository) GetActiveShift(ctx context.Context) (*Shift, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Shift, error) {
		s := &Shift{}
		err := tx.QueryRowContext(ctx,
			"SELECT id, user_id, opened_at, closed_at, opening_balance, closing_balance, status FROM shifts WHERE status='open' AND tenant_id=$1 LIMIT 1",
			tenantID).Scan(&s.ID, &s.UserID, &s.OpenedAt, &s.ClosedAt, &s.OpeningBalance, &s.ClosingBalance, &s.Status)
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

func (r *ShiftRepository) Create(ctx context.Context, userID int64, openingBalance float64) (*Shift, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Shift, error) {
		s := &Shift{}
		err := tx.QueryRowContext(ctx,
			"INSERT INTO shifts (tenant_id, user_id, opening_balance, status) VALUES ($1,$2,$3,'open') RETURNING id, user_id, opened_at, closed_at, opening_balance, closing_balance, status",
			tenantID, userID, openingBalance).Scan(&s.ID, &s.UserID, &s.OpenedAt, &s.ClosedAt, &s.OpeningBalance, &s.ClosingBalance, &s.Status)
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

func (r *ShiftRepository) Close(ctx context.Context, id int64, closingBalance float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE shifts SET closing_balance=$1, closed_at=NOW(), status='closed' WHERE id=$2 AND status='open' AND tenant_id=$3",
			closingBalance, id, tenantID)
		return err
	})
}
