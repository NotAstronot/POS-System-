package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Department struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) ListAll(ctx context.Context) ([]Department, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Department, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, description, is_active, created_at, updated_at FROM departments WHERE tenant_id=$1 ORDER BY id ASC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []Department
		for rows.Next() {
			var d Department
			if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, d)
		}
		return items, nil
	})
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id int64) (*Department, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Department, error) {
		d := &Department{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, description, is_active, created_at, updated_at FROM departments WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&d.ID, &d.Name, &d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return d, nil
	})
}

func (r *DepartmentRepository) Create(ctx context.Context, d *Department) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "INSERT INTO departments (name, description, tenant_id) VALUES ($1,$2,$3) RETURNING id", d.Name, d.Description, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *DepartmentRepository) Update(ctx context.Context, d *Department) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE departments SET name=$1, description=$2, is_active=$3, updated_at=NOW() WHERE id=$4 AND tenant_id=$5", d.Name, d.Description, d.IsActive, d.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM departments WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
