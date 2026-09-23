package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List(ctx context.Context) ([]Category, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Category, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, created_at, updated_at FROM categories ORDER BY id")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var list []Category
		for rows.Next() {
			var c Category
			if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, c)
		}
		return list, rows.Err()
	})
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*Category, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Category, error) {
		c := &Category{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM categories WHERE id=$1", id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CategoryRepository) Create(ctx context.Context, name string) (*Category, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Category, error) {
		c := &Category{}
		err := tx.QueryRowContext(ctx, "INSERT INTO categories (name, tenant_id) VALUES ($1, $2) RETURNING id, name, created_at, updated_at", name, tenantID).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CategoryRepository) Update(ctx context.Context, id int64, name string) (*Category, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Category, error) {
		c := &Category{}
		err := tx.QueryRowContext(ctx, "UPDATE categories SET name=$1, updated_at=NOW() WHERE id=$2 RETURNING id, name, created_at, updated_at", name, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM categories WHERE id=$1", id)
		return err
	})
}