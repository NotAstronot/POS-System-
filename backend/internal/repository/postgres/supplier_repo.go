package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Supplier struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ContactName string    `json:"contact_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	NPWP        string    `json:"npwp"`
	Address     string    `json:"address"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) ListAll(ctx context.Context) ([]Supplier, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Supplier, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, contact_name, phone, email, npwp, address, is_active, created_at, updated_at FROM suppliers WHERE tenant_id=$1 ORDER BY id DESC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		items := make([]Supplier, 0)
		for rows.Next() {
			var s Supplier
			if err := rows.Scan(&s.ID, &s.Name, &s.ContactName, &s.Phone, &s.Email, &s.NPWP, &s.Address, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
				return nil, err
			}
			items = append(items, s)
		}
		return items, nil
	})
}

func (r *SupplierRepository) GetByID(ctx context.Context, id int64) (*Supplier, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Supplier, error) {
		s := &Supplier{}
		err := tx.QueryRowContext(ctx, "SELECT id, name, contact_name, phone, email, npwp, address, is_active, created_at, updated_at FROM suppliers WHERE id=$1 AND tenant_id=$2", id, tenantID).
			Scan(&s.ID, &s.Name, &s.ContactName, &s.Phone, &s.Email, &s.NPWP, &s.Address, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

func (r *SupplierRepository) Create(ctx context.Context, s *Supplier) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO suppliers (name, contact_name, phone, email, npwp, address, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
			s.Name, s.ContactName, s.Phone, s.Email, s.NPWP, s.Address, tenantID,
		).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *SupplierRepository) Update(ctx context.Context, s *Supplier) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE suppliers SET name=$1, contact_name=$2, phone=$3, email=$4, npwp=$5, address=$6, is_active=$7, updated_at=NOW() WHERE id=$8 AND tenant_id=$9",
			s.Name, s.ContactName, s.Phone, s.Email, s.NPWP, s.Address, s.IsActive, s.ID, tenantID,
		)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *SupplierRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM suppliers WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
