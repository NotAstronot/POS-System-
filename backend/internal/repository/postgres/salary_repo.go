package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Salary struct {
	ID              int64      `json:"id"`
	EmployeeID      int64      `json:"employee_id"`
	EmployeeName    string     `json:"employee_name"`
	PayPeriod       string     `json:"pay_period"`
	BaseSalary      float64    `json:"base_salary"`
	CommissionTotal float64    `json:"commission_total"`
	Deduction       float64    `json:"deduction"`
	NetSalary       float64    `json:"net_salary"`
	Status          string     `json:"status"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SalaryRepository struct {
	db *sql.DB
}

func NewSalaryRepository(db *sql.DB) *SalaryRepository {
	return &SalaryRepository{db: db}
}

func (r *SalaryRepository) ListAll(ctx context.Context) ([]Salary, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Salary, error) {
		rows, err := tx.QueryContext(ctx, `SELECT s.id, s.employee_id, COALESCE(e.name,'') as employee_name,
			s.pay_period, s.base_salary, s.commission_total, s.deduction, s.net_salary, s.status, s.paid_at, s.created_at, s.updated_at
			FROM salaries s LEFT JOIN employees e ON s.employee_id=e.id WHERE s.tenant_id=$1 ORDER BY s.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Salary
		for rows.Next() {
			var s Salary
			if err := rows.Scan(&s.ID, &s.EmployeeID, &s.EmployeeName, &s.PayPeriod, &s.BaseSalary, &s.CommissionTotal, &s.Deduction, &s.NetSalary, &s.Status, &s.PaidAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, s)
		}
		return items, nil
	})
}

func (r *SalaryRepository) GetByID(ctx context.Context, id int64) (*Salary, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Salary, error) {
		s := &Salary{}
		err := tx.QueryRowContext(ctx, `SELECT s.id, s.employee_id, COALESCE(e.name,'') as employee_name,
			s.pay_period, s.base_salary, s.commission_total, s.deduction, s.net_salary, s.status, s.paid_at, s.created_at, s.updated_at
			FROM salaries s LEFT JOIN employees e ON s.employee_id=e.id WHERE s.id=$1 AND s.tenant_id=$2`, id, tenantID).
			Scan(&s.ID, &s.EmployeeID, &s.EmployeeName, &s.PayPeriod, &s.BaseSalary, &s.CommissionTotal, &s.Deduction, &s.NetSalary, &s.Status, &s.PaidAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

func (r *SalaryRepository) Create(ctx context.Context, s *Salary) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		net := s.BaseSalary + s.CommissionTotal - s.Deduction
		err := tx.QueryRowContext(ctx,
			"INSERT INTO salaries (employee_id, pay_period, base_salary, commission_total, deduction, net_salary, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			s.EmployeeID, s.PayPeriod, s.BaseSalary, s.CommissionTotal, s.Deduction, net, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *SalaryRepository) Update(ctx context.Context, s *Salary) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		net := s.BaseSalary + s.CommissionTotal - s.Deduction
		_, err := tx.ExecContext(ctx,
			"UPDATE salaries SET employee_id=$1, pay_period=$2, base_salary=$3, commission_total=$4, deduction=$5, net_salary=$6, status=$7, updated_at=NOW() WHERE id=$8 AND tenant_id=$9",
			s.EmployeeID, s.PayPeriod, s.BaseSalary, s.CommissionTotal, s.Deduction, net, s.Status, s.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalaryRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM salaries WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalaryRepository) MarkPaid(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE salaries SET status='paid', paid_at=NOW(), updated_at=NOW() WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
