package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Branch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) ListAll(ctx context.Context) ([]Branch, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Branch, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, address, phone, is_active, created_at, updated_at FROM branches WHERE tenant_id=$1 ORDER BY id ASC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var branches []Branch
		for rows.Next() {
			var b Branch
			if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
				return nil, err
			}
			branches = append(branches, b)
		}
		return branches, nil
	})
}

func (r *BranchRepository) GetByID(ctx context.Context, id int64) (*Branch, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Branch, error) {
		b := &Branch{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, address, phone, is_active, created_at, updated_at FROM branches WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return b, nil
	})
}

func (r *BranchRepository) Create(ctx context.Context, b *Branch) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO branches (name, address, phone, tenant_id) VALUES ($1,$2,$3,$4) RETURNING id",
			b.Name, b.Address, b.Phone, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *BranchRepository) Update(ctx context.Context, b *Branch) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE branches SET name=$1, address=$2, phone=$3, is_active=$4, updated_at=NOW() WHERE id=$5 AND tenant_id=$6",
			b.Name, b.Address, b.Phone, b.IsActive, b.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BranchRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM branches WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
