package database

import (
	"context"
	"database/sql"

	"pos-system/internal/domain"
)

// RequireTenant returns tenant_id from context.Context or ErrTenantRequired.
// Tenant must have been injected by the TenantContext middleware.
func RequireTenant(ctx context.Context) (int64, error) {
	id, ok := domain.GetTenantID(ctx)
	if !ok || id <= 0 {
		return 0, ErrTenantRequired
	}
	return id, nil
}

// BeginTenantTx begins a transaction and injects tenant_id from ctx into RLS.
// tenant_id is always taken from context — never from a separate parameter.
func BeginTenantTx(ctx context.Context, db *sql.DB) (*sql.Tx, int64, error) {
	tenantID, err := RequireTenant(ctx)
	if err != nil {
		return nil, 0, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	if err := SetRLSTenant(tx, tenantID); err != nil {
		_ = tx.Rollback()
		return nil, 0, err
	}
	return tx, tenantID, nil
}

// WithTenantTx runs fn in a transaction with RLS scoped to ctx's tenant_id.
func WithTenantTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx, tenantID int64) error) error {
	tx, tenantID, err := BeginTenantTx(ctx, db)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx, tenantID); err != nil {
		return err
	}
	return tx.Commit()
}

// WithTenantTx1 is WithTenantTx returning one value.
func WithTenantTx1[T any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx, tenantID int64) (T, error)) (T, error) {
	var result T
	err := WithTenantTx(ctx, db, func(tx *sql.Tx, tenantID int64) error {
		var err error
		result, err = fn(tx, tenantID)
		return err
	})
	return result, err
}

// WithTenantTx2 is WithTenantTx returning two values.
func WithTenantTx2[T1, T2 any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx, tenantID int64) (T1, T2, error)) (T1, T2, error) {
	var result1 T1
	var result2 T2
	err := WithTenantTx(ctx, db, func(tx *sql.Tx, tenantID int64) error {
		var err error
		result1, result2, err = fn(tx, tenantID)
		return err
	})
	return result1, result2, err
}

// WithExplicitTenantTx runs fn with RLS for an explicit tenantID.
// Use only for system paths (login bootstrap, super-admin targeting another tenant).
// Still requires a non-nil ctx for cancellation.
func WithExplicitTenantTx(ctx context.Context, db *sql.DB, tenantID int64, fn func(tx *sql.Tx) error) error {
	if tenantID <= 0 {
		return ErrTenantRequired
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := SetRLSTenant(tx, tenantID); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// WithExplicitTenantTx1 is WithExplicitTenantTx returning one value.
func WithExplicitTenantTx1[T any](ctx context.Context, db *sql.DB, tenantID int64, fn func(tx *sql.Tx) (T, error)) (T, error) {
	var result T
	err := WithExplicitTenantTx(ctx, db, tenantID, func(tx *sql.Tx) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}

// WithSystemTx runs fn without RLS (tables without tenant isolation, e.g. tenants catalog).
func WithSystemTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// WithSystemTx1 is WithSystemTx returning one value.
func WithSystemTx1[T any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) (T, error)) (T, error) {
	var result T
	err := WithSystemTx(ctx, db, func(tx *sql.Tx) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}
