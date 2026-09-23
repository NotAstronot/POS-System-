package postgres

import (
	"context"
	"database/sql"
	"time"
)

type CustomerCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceScheme string    `json:"price_scheme"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CustomerCategoryRepository struct {
	db *sql.DB
}

func NewCustomerCategoryRepository(db *sql.DB) *CustomerCategoryRepository {
	return &CustomerCategoryRepository{db: db}
}

func (r *CustomerCategoryRepository) ListAll(ctx context.Context) ([]CustomerCategory, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]CustomerCategory, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, description, price_scheme, is_active, created_at, updated_at FROM customer_categories WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]CustomerCategory, 0)
		for rows.Next() {
			var c CustomerCategory
			if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.PriceScheme, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, c)
		}
		return items, nil
	})
}

func (r *CustomerCategoryRepository) GetByID(ctx context.Context, id int64) (*CustomerCategory, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*CustomerCategory, error) {
		c := &CustomerCategory{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, description, price_scheme, is_active, created_at, updated_at FROM customer_categories WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&c.ID, &c.Name, &c.Description, &c.PriceScheme, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CustomerCategoryRepository) Create(ctx context.Context, c *CustomerCategory) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "INSERT INTO customer_categories (name, description, price_scheme, tenant_id) VALUES ($1,$2,$3,$4) RETURNING id",
			c.Name, c.Description, c.PriceScheme, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *CustomerCategoryRepository) Update(ctx context.Context, c *CustomerCategory) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE customer_categories SET name=$1, description=$2, price_scheme=$3, is_active=$4, updated_at=NOW() WHERE id=$5 AND tenant_id=$6",
			c.Name, c.Description, c.PriceScheme, c.IsActive, c.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *CustomerCategoryRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM customer_categories WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
