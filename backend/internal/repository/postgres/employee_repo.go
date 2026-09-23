package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Employee struct {
	ID           int64     `json:"id"`
	NIK          string    `json:"nik"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	DepartmentID *int64    `json:"department_id"`
	DeptName     string    `json:"dept_name"`
	Position     string    `json:"position"`
	HireDate     string    `json:"hire_date"`
	SalaryBase   float64   `json:"salary_base"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

const employeeCols = `e.id, e.nik, e.name, e.email, e.phone, e.department_id,
	COALESCE(d.name,'') as dept_name, e.position, e.hire_date::text, e.salary_base,
	e.is_active, e.created_at, e.updated_at`

func (r *EmployeeRepository) ListAll(ctx context.Context) ([]Employee, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Employee, error) {
		rows, err := tx.QueryContext(ctx, "SELECT "+employeeCols+" FROM employees e LEFT JOIN departments d ON e.department_id=d.id WHERE e.tenant_id=$1 ORDER BY e.id ASC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Employee
		for rows.Next() {
			var e Employee
			if err := rows.Scan(&e.ID, &e.NIK, &e.Name, &e.Email, &e.Phone, &e.DepartmentID, &e.DeptName, &e.Position, &e.HireDate, &e.SalaryBase, &e.IsActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, e)
		}
		return items, nil
	})
}

func (r *EmployeeRepository) GetByID(ctx context.Context, id int64) (*Employee, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Employee, error) {
		e := &Employee{}
		err := tx.QueryRowContext(ctx, "SELECT "+employeeCols+" FROM employees e LEFT JOIN departments d ON e.department_id=d.id WHERE e.id=$1 AND e.tenant_id=$2", id, tenantID).
			Scan(&e.ID, &e.NIK, &e.Name, &e.Email, &e.Phone, &e.DepartmentID, &e.DeptName, &e.Position, &e.HireDate, &e.SalaryBase, &e.IsActive, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return e, nil
	})
}

func (r *EmployeeRepository) Create(ctx context.Context, e *Employee) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO employees (nik, name, email, phone, department_id, position, hire_date, salary_base, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id",
			e.NIK, e.Name, e.Email, e.Phone, e.DepartmentID, e.Position, e.HireDate, e.SalaryBase, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *EmployeeRepository) Update(ctx context.Context, e *Employee) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE employees SET name=$1, email=$2, phone=$3, department_id=$4, position=$5, hire_date=$6, salary_base=$7, is_active=$8, updated_at=NOW() WHERE id=$9 AND tenant_id=$10",
			e.Name, e.Email, e.Phone, e.DepartmentID, e.Position, e.HireDate, e.SalaryBase, e.IsActive, e.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM employees WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
