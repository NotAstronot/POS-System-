package postgres

import (
	"context"
	"database/sql"
	"time"
)

type RevenueSummary struct {
	Total        float64       `json:"total"`
	Count        int           `json:"count"`
	Period       string        `json:"period"`
	Timestamp    string        `json:"timestamp"`
	RecentOrders []RecentOrder `json:"recent_orders"`
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
		var total float64
		var count int
		var trunc string
		switch period {
		case "weekly":
			trunc = "IYYY-IW"
		case "monthly":
			trunc = "YYYY-MM"
		default:
			trunc = "YYYY-MM-DD"
		}
		err := tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(total),0), COUNT(*)
			FROM orders
			WHERE tenant_id=$1 AND status='completed'
			  AND TO_CHAR(created_at, $2) = TO_CHAR(NOW(), $2)`,
			tenantID, trunc).Scan(&total, &count)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		return &RevenueSummary{
			Total:     total,
			Count:     count,
			Period:    period,
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	})
}
