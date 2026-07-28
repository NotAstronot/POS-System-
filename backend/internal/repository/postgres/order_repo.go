package postgres

import (
	"database/sql"
	"time"
)

type Order struct {
	ID        int64     `json:"id"`
	ShiftID   int64     `json:"shift_id"`
	UserID    int64     `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(id int64) (*Order, error) {
	o := &Order{}
	err := r.db.QueryRow("SELECT id, shift_id, user_id, total, status, created_at, updated_at FROM orders WHERE id=$1", id).Scan(&o.ID, &o.ShiftID, &o.UserID, &o.Total, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}
