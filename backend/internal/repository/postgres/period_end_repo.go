package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PeriodClosing struct {
	ID        int64      `json:"id"`
	Period    string     `json:"period"`
	Status    string     `json:"status"`
	OpenedBy  *int64     `json:"opened_by"`
	OpenedAt  *time.Time `json:"opened_at"`
	ClosedBy  *int64     `json:"closed_by"`
	ClosedAt  *time.Time `json:"closed_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type FxDifference struct {
	ID          int64     `json:"id"`
	Period      string    `json:"period"`
	Currency    string    `json:"currency"`
	RateDiff    float64   `json:"rate_diff"`
	GainAmount  float64   `json:"gain_amount"`
	LossAmount  float64   `json:"loss_amount"`
	Notes       string    `json:"notes"`
	Description string    `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type DepreciationResult struct {
	TotalDepreciation float64 `json:"total_depreciation"`
	AssetsProcessed   int     `json:"assets_processed"`
}

type PeriodEnd struct {
	ID                int64      `json:"id"`
	PeriodStartDate   string     `json:"period_start_date"`
	PeriodEndDate     string     `json:"period_end_date"`
	Status            string     `json:"status"`
	TotalRevenue      float64    `json:"total_revenue"`
	TotalExpense      float64    `json:"total_expense"`
	NetIncome         float64    `json:"net_income"`
	DepreciationTotal float64    `json:"depreciation_total"`
	ClosedBy          *int64     `json:"closed_by"`
	ClosedAt          *time.Time `json:"closed_at"`
	CreatedBy         *int64     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type PeriodEndRepository struct {
	db *sql.DB
}

func NewPeriodEndRepository(db *sql.DB) *PeriodEndRepository {
	return &PeriodEndRepository{db: db}
}

func (r *PeriodEndRepository) ListAll(ctx context.Context) ([]PeriodEnd, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PeriodEnd, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, period_start_date::text, period_end_date::text, status, total_revenue, total_expense, net_income, depreciation_total, closed_by, closed_at, created_by, created_at, updated_at
			FROM period_ends WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PeriodEnd, 0)
		for rows.Next() {
			var p PeriodEnd
			if err := rows.Scan(&p.ID, &p.PeriodStartDate, &p.PeriodEndDate, &p.Status, &p.TotalRevenue, &p.TotalExpense, &p.NetIncome, &p.DepreciationTotal, &p.ClosedBy, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, p)
		}
		return list, nil
	})
}

func (r *PeriodEndRepository) GetByID(ctx context.Context, id int64) (*PeriodEnd, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PeriodEnd, error) {
		p := &PeriodEnd{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, period_start_date::text, period_end_date::text, status, total_revenue, total_expense, net_income, depreciation_total, closed_by, closed_at, created_by, created_at, updated_at
			FROM period_ends WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&p.ID, &p.PeriodStartDate, &p.PeriodEndDate, &p.Status, &p.TotalRevenue, &p.TotalExpense, &p.NetIncome, &p.DepreciationTotal, &p.ClosedBy, &p.ClosedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return p, nil
	})
}

func (r *PeriodEndRepository) Create(ctx context.Context, p *PeriodEnd) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO period_ends (period_start_date, period_end_date, total_revenue, total_expense, net_income, depreciation_total, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
			p.PeriodStartDate, p.PeriodEndDate, p.TotalRevenue, p.TotalExpense, p.NetIncome, p.DepreciationTotal, p.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *PeriodEndRepository) RunDepreciation(ctx context.Context, period string) (*DepreciationResult, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*DepreciationResult, error) {
		_, err := tx.ExecContext(ctx, `
			UPDATE fixed_assets SET accumulated_depreciation = accumulated_depreciation + monthly_depreciation, updated_at = NOW()
			WHERE is_active = true AND tenant_id = $1 AND accumulated_depreciation < acquisition_value`, tenantID)
		if err != nil {
			return nil, err
		}
		var totalDep float64
		var assetsProcessed int
		err = tx.QueryRowContext(ctx, "SELECT COALESCE(SUM(monthly_depreciation), 0), COUNT(*) FROM fixed_assets WHERE is_active = true AND tenant_id = $1 AND accumulated_depreciation < acquisition_value", tenantID).Scan(&totalDep, &assetsProcessed)
		if err != nil {
			return nil, err
		}
		return &DepreciationResult{TotalDepreciation: totalDep, AssetsProcessed: assetsProcessed}, nil
	})
}

func (r *PeriodEndRepository) ClosePeriod(ctx context.Context, period string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, `
			UPDATE period_closings SET status='closed', closed_by=$1, closed_at=NOW()
			WHERE period=$2 AND tenant_id=$3`, tenantID, period, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PeriodEndRepository) ListPeriods(ctx context.Context) ([]PeriodClosing, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PeriodClosing, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, period, status, opened_by, opened_at, closed_by, closed_at, created_at
			FROM period_closings WHERE tenant_id=$1 ORDER BY period DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PeriodClosing, 0)
		for rows.Next() {
			var p PeriodClosing
			if err := rows.Scan(&p.ID, &p.Period, &p.Status, &p.OpenedBy, &p.OpenedAt, &p.ClosedBy, &p.ClosedAt, &p.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, p)
		}
		return list, nil
	})
}

func (r *PeriodEndRepository) GetPeriod(ctx context.Context, id int64) (*PeriodClosing, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PeriodClosing, error) {
		p := &PeriodClosing{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, period, status, opened_by, opened_at, closed_by, closed_at, created_at
			FROM period_closings WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&p.ID, &p.Period, &p.Status, &p.OpenedBy, &p.OpenedAt, &p.ClosedBy, &p.ClosedAt, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		return p, nil
	})
}

func (r *PeriodEndRepository) OpenPeriod(ctx context.Context, period string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		var exists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM period_closings WHERE period=$1 AND tenant_id=$2)", period, tenantID).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("periode %s sudah ada", period)
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO period_closings (period, status, opened_by, tenant_id) VALUES ($1, 'open', $2, $3)", period, tenantID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PeriodEndRepository) ListFx(ctx context.Context, period string) ([]FxDifference, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]FxDifference, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, period, currency, rate_diff, gain_amount, loss_amount, notes, created_by, created_at
			FROM fx_differences WHERE period=$1 AND tenant_id=$2 ORDER BY id`, period, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]FxDifference, 0)
		for rows.Next() {
			var f FxDifference
			if err := rows.Scan(&f.ID, &f.Period, &f.Currency, &f.RateDiff, &f.GainAmount, &f.LossAmount, &f.Notes, &f.CreatedBy, &f.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, f)
		}
		return list, nil
	})
}

func (r *PeriodEndRepository) CreateFx(ctx context.Context, f *FxDifference) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO fx_differences (period, currency, rate_diff, gain_amount, loss_amount, notes, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
			f.Period, f.Currency, f.RateDiff, f.GainAmount, f.LossAmount, f.Notes, f.CreatedBy, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
