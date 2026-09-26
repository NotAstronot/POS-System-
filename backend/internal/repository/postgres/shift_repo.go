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

// ShiftHistory adalah satu baris rekap shift kasir beserta ringkasan
// transaksi & omzet yang terjadi selama shift tersebut.
type ShiftHistory struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	UserName       string     `json:"user_name"`
	OpenedAt       time.Time  `json:"opened_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	OpeningBalance float64    `json:"opening_balance"`
	ClosingBalance *float64   `json:"closing_balance"`
	Status         string     `json:"status"`
	Transactions   int        `json:"transactions"`
	TotalSales     float64    `json:"total_sales"`
}

// History mengembalikan daftar shift terbaru (beserta kasir, jumlah
// transaksi, dan total omzet order completed per shift).
func (r *ShiftRepository) History(ctx context.Context, limit int) ([]ShiftHistory, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ShiftHistory, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT s.id, s.user_id, COALESCE(u.name,''), s.opened_at, s.closed_at,
				s.opening_balance, s.closing_balance, s.status,
				COUNT(o.id), COALESCE(SUM(o.total),0)
			FROM shifts s
			LEFT JOIN users u ON u.id = s.user_id AND u.tenant_id = s.tenant_id::text
			LEFT JOIN orders o ON o.shift_id = s.id AND o.tenant_id = s.tenant_id AND o.status='completed'
			WHERE s.tenant_id=$1
			GROUP BY s.id, s.user_id, u.name, s.opened_at, s.closed_at,
				s.opening_balance, s.closing_balance, s.status
			ORDER BY s.opened_at DESC LIMIT $2`, tenantID, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := []ShiftHistory{}
		for rows.Next() {
			var h ShiftHistory
			if err := rows.Scan(&h.ID, &h.UserID, &h.UserName, &h.OpenedAt, &h.ClosedAt,
				&h.OpeningBalance, &h.ClosingBalance, &h.Status, &h.Transactions, &h.TotalSales); err != nil {
				return nil, err
			}
			list = append(list, h)
		}
		return list, rows.Err()
	})
}
