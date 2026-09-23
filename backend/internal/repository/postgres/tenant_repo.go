package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Tenant struct {
	ID                    int64      `json:"id"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	Domain                string     `json:"domain"`
	LogoURL               string     `json:"logo_url"`
	SubscriptionPlan      string     `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	MaxUsers              int        `json:"max_users"`
	MaxProducts           int        `json:"max_products"`
	MaxBranches           int        `json:"max_branches"`
	Settings              string     `json:"settings"`
	Status                string     `json:"status"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

const tenantCols = "id, name, slug, domain, logo_url, subscription_plan, subscription_expires_at, max_users, max_products, max_branches, settings, status, created_at, updated_at"

func scanTenant(t *Tenant, row scanRow) error {
	return row.Scan(&t.ID, &t.Name, &t.Slug, &t.Domain, &t.LogoURL, &t.SubscriptionPlan, &t.SubscriptionExpiresAt, &t.MaxUsers, &t.MaxProducts, &t.MaxBranches, &t.Settings, &t.Status, &t.CreatedAt, &t.UpdatedAt)
}

// Tenants catalog has no RLS — system scope (super-admin only routes).

func (r *TenantRepository) List(ctx context.Context) ([]Tenant, error) {
	return withSystemTx1(ctx, r.db, func(tx *sql.Tx) ([]Tenant, error) {
		rows, err := tx.QueryContext(ctx, "SELECT "+tenantCols+" FROM tenants ORDER BY id ASC")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var list []Tenant
		for rows.Next() {
			var t Tenant
			if err := scanTenant(&t, rows); err != nil {
				return nil, err
			}
			list = append(list, t)
		}
		return list, rows.Err()
	})
}

func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*Tenant, error) {
	return withSystemTx1(ctx, r.db, func(tx *sql.Tx) (*Tenant, error) {
		t := &Tenant{}
		err := scanTenant(t, tx.QueryRowContext(ctx, "SELECT "+tenantCols+" FROM tenants WHERE id=$1", id))
		if err != nil {
			return nil, err
		}
		return t, nil
	})
}

func (r *TenantRepository) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	return withSystemTx1(ctx, r.db, func(tx *sql.Tx) (*Tenant, error) {
		t := &Tenant{}
		err := scanTenant(t, tx.QueryRowContext(ctx, "SELECT "+tenantCols+" FROM tenants WHERE slug=$1", slug))
		if err != nil {
			return nil, err
		}
		return t, nil
	})
}

func (r *TenantRepository) Create(ctx context.Context, t *Tenant) (*Tenant, error) {
	return withSystemTx1(ctx, r.db, func(tx *sql.Tx) (*Tenant, error) {
		err := tx.QueryRowContext(ctx,
			`INSERT INTO tenants (name, slug, domain, logo_url, subscription_plan, subscription_expires_at, max_users, max_products, max_branches, settings, status)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`,
			t.Name, t.Slug, t.Domain, t.LogoURL, t.SubscriptionPlan, t.SubscriptionExpiresAt,
			t.MaxUsers, t.MaxProducts, t.MaxBranches, t.Settings, t.Status,
		).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *TenantRepository) Update(ctx context.Context, t *Tenant) error {
	return withSystemTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`UPDATE tenants SET name=$1, slug=$2, domain=$3, logo_url=$4, subscription_plan=$5, subscription_expires_at=$6,
			 max_users=$7, max_products=$8, max_branches=$9, settings=$10, status=$11, updated_at=NOW() WHERE id=$12`,
			t.Name, t.Slug, t.Domain, t.LogoURL, t.SubscriptionPlan, t.SubscriptionExpiresAt,
			t.MaxUsers, t.MaxProducts, t.MaxBranches, t.Settings, t.Status, t.ID,
		)
		return err
	})
}

func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	return withSystemTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM tenants WHERE id=$1", id)
		return err
	})
}

// Usage counters: super-admin inspects a target tenant via explicit tenantID
// (not the caller's JWT tenant).

func (r *TenantRepository) CountUsers(ctx context.Context, tenantID int64) (int, error) {
	return databaseWithExplicitCount1(ctx, r.db, tenantID,
		"SELECT COUNT(*) FROM users WHERE tenant_id=$1")
}

func (r *TenantRepository) CountProducts(ctx context.Context, tenantID int64) (int, error) {
	return databaseWithExplicitCount1(ctx, r.db, tenantID,
		"SELECT COUNT(*) FROM products WHERE tenant_id=$1")
}

func (r *TenantRepository) CountBranches(ctx context.Context, tenantID int64) (int, error) {
	return databaseWithExplicitCount1(ctx, r.db, tenantID,
		"SELECT COUNT(*) FROM branches WHERE tenant_id=$1")
}

func databaseWithExplicitCount1(ctx context.Context, db *sql.DB, tenantID int64, query string) (int, error) {
	return withExplicitTenantTx1(ctx, db, tenantID, func(tx *sql.Tx) (int, error) {
		var count int
		err := tx.QueryRowContext(ctx, query, tenantID).Scan(&count)
		return count, err
	})
}
