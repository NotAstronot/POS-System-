package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Commission struct {
	ID             int64     `json:"id"`
	EmployeeID     int64     `json:"employee_id"`
	EmployeeName   string    `json:"employee_name"`
	CommissionType string    `json:"commission_type"`
	Rate           float64   `json:"rate"`
	Description    string    `json:"description"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CommissionRepository struct {
	db *sql.DB
}

func NewCommissionRepository(db *sql.DB) *CommissionRepository {
	return &CommissionRepository{db: db}
}

func (r *CommissionRepository) ListAll(ctx context.Context) ([]Commission, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Commission, error) {
		rows, err := tx.QueryContext(ctx, `SELECT c.id, c.employee_id, COALESCE(e.name,'') as employee_name,
			c.commission_type, c.rate, c.description, c.is_active, c.created_at, c.updated_at
			FROM commissions c LEFT JOIN employees e ON c.employee_id=e.id WHERE c.tenant_id=$1 ORDER BY c.id ASC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Commission
		for rows.Next() {
			var c Commission
			if err := rows.Scan(&c.ID, &c.EmployeeID, &c.EmployeeName, &c.CommissionType, &c.Rate, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, c)
		}
		return items, nil
	})
}

func (r *CommissionRepository) GetByID(ctx context.Context, id int64) (*Commission, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Commission, error) {
		c := &Commission{}
		err := tx.QueryRowContext(ctx, `SELECT c.id, c.employee_id, COALESCE(e.name,'') as employee_name,
			c.commission_type, c.rate, c.description, c.is_active, c.created_at, c.updated_at
			FROM commissions c LEFT JOIN employees e ON c.employee_id=e.id WHERE c.id=$1 AND c.tenant_id=$2`, id, tenantID).
			Scan(&c.ID, &c.EmployeeID, &c.EmployeeName, &c.CommissionType, &c.Rate, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CommissionRepository) Create(ctx context.Context, c *Commission) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO commissions (employee_id, commission_type, rate, description, tenant_id) VALUES ($1,$2,$3,$4,$5) RETURNING id",
			c.EmployeeID, c.CommissionType, c.Rate, c.Description, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *CommissionRepository) Update(ctx context.Context, c *Commission) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE commissions SET employee_id=$1, commission_type=$2, rate=$3, description=$4, is_active=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			c.EmployeeID, c.CommissionType, c.Rate, c.Description, c.IsActive, c.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *CommissionRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM commissions WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
