package postgres

import (
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

func (r *ShiftRepository) GetActiveShift() (*Shift, error) {
	s := &Shift{}
	err := r.db.QueryRow("SELECT id, user_id, opened_at, closed_at, opening_balance, closing_balance, status FROM shifts WHERE status='open' LIMIT 1").Scan(&s.ID, &s.UserID, &s.OpenedAt, &s.ClosedAt, &s.OpeningBalance, &s.ClosingBalance, &s.Status)
	if err != nil {
		return nil, err
	}
	return s, nil
}
