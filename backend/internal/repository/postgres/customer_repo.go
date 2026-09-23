package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Customer struct {
	ID                   int64     `json:"id"`
	Name                 string    `json:"name"`
	Phone                string    `json:"phone"`
	Email                string    `json:"email"`
	NPWP                 string    `json:"npwp"`
	Address              string    `json:"address"`
	CustomerCategoryID   *int64    `json:"customer_category_id"`
	CustomerCategoryName string    `json:"customer_category_name"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) ListAll(ctx context.Context) ([]Customer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Customer, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT c.id, c.name, c.phone, c.email, c.npwp, c.address, c.customer_category_id, COALESCE(cc.name, ''),
				c.is_active, c.created_at, c.updated_at
			FROM customers c
			LEFT JOIN customer_categories cc ON cc.id = c.customer_category_id
			WHERE c.tenant_id=$1
			ORDER BY c.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]Customer, 0)
		for rows.Next() {
			var c Customer
			if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.NPWP, &c.Address, &c.CustomerCategoryID, &c.CustomerCategoryName,
				&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, c)
		}
		return items, nil
	})
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*Customer, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Customer, error) {
		c := &Customer{}
		err := tx.QueryRowContext(ctx, `
			SELECT c.id, c.name, c.phone, c.email, c.npwp, c.address, c.customer_category_id, COALESCE(cc.name, ''),
				c.is_active, c.created_at, c.updated_at
			FROM customers c
			LEFT JOIN customer_categories cc ON cc.id = c.customer_category_id
			WHERE c.id=$1 AND c.tenant_id=$2`, id, tenantID).
			Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.NPWP, &c.Address, &c.CustomerCategoryID, &c.CustomerCategoryName,
				&c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}

func (r *CustomerRepository) Create(ctx context.Context, c *Customer) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, "INSERT INTO customers (name, phone, email, npwp, address, customer_category_id, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			c.Name, c.Phone, c.Email, c.NPWP, c.Address, c.CustomerCategoryID, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *CustomerRepository) Update(ctx context.Context, c *Customer) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE customers SET name=$1, phone=$2, email=$3, npwp=$4, address=$5, customer_category_id=$6, is_active=$7, updated_at=NOW() WHERE id=$8 AND tenant_id=$9",
			c.Name, c.Phone, c.Email, c.NPWP, c.Address, c.CustomerCategoryID, c.IsActive, c.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM customers WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
