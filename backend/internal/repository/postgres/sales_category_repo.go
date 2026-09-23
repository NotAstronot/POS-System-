package postgres

import (
	"context"
	"database/sql"
	"time"
)

type SalesCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SalesCategoryRepository struct {
	db *sql.DB
}

func NewSalesCategoryRepository(db *sql.DB) *SalesCategoryRepository {
	return &SalesCategoryRepository{db: db}
}

func (r *SalesCategoryRepository) ListAll(ctx context.Context) ([]SalesCategory, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]SalesCategory, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, description, is_active, created_at, updated_at FROM sales_categories WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]SalesCategory, 0)
		for rows.Next() {
			var c SalesCategory
			if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, c)
		}
		return items, nil
	})
}

func (r *SalesCategoryRepository) GetByID(ctx context.Context, id int64) (*SalesCategory, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*SalesCategory, error) {
		c := &SalesCategory{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, description, is_active, created_at, updated_at FROM sales_categories WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&c.ID, &c.Name, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *SalesCategoryRepository) Create(ctx context.Context, c *SalesCategory) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "INSERT INTO sales_categories (name, description, tenant_id) VALUES ($1,$2,$3) RETURNING id",
			c.Name, c.Description, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *SalesCategoryRepository) Update(ctx context.Context, c *SalesCategory) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE sales_categories SET name=$1, description=$2, is_active=$3, updated_at=NOW() WHERE id=$4 AND tenant_id=$5",
			c.Name, c.Description, c.IsActive, c.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SalesCategoryRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sales_categories WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
