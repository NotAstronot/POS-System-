package postgres

import (
	"context"
	"database/sql"
)

type RevenueSummary struct {
	Total     float64 `json:"total"`
	Count     int     `json:"count"`
	Period    string  `json:"period"`
	Timestamp string  `json:"timestamp"`
}

type RecentOrder struct {
	ID            int64   `json:"id"`
	OrderNumber   string  `json:"order_number"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"payment_method"`
	CashierName   string  `json:"cashier_name"`
	ItemCount     int     `json:"item_count"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
}

type RevenueRepository struct {
	db *sql.DB
}

func NewRevenueRepository(db *sql.DB) *RevenueRepository {
	return &RevenueRepository{db: db}
}

func (r *RevenueRepository) Outlets(ctx context.Context) ([]RecentOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]RecentOrder, error) {
		return []RecentOrder{}, nil
	})
}

func (r *RevenueRepository) Summary(ctx context.Context, period string) (*RevenueSummary, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*RevenueSummary, error) {
		return &RevenueSummary{}, nil
	})
}