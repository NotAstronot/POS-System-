package postgres

import (
	"context"
	"database/sql"

	"pos-system/internal/database"
)

// beginTenantTx starts a transaction with tenant_id taken from ctx (RLS injected).
// Prefer this (or database.WithTenantTx*) in repository methods.
func beginTenantTx(ctx context.Context, db *sql.DB) (*sql.Tx, int64, error) {
	return database.BeginTenantTx(ctx, db)
}

// withTenantTx runs fn with a context-scoped tenant transaction.
func withTenantTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx, tenantID int64) error) error {
	return database.WithTenantTx(ctx, db, fn)
}

// withTenantTx1 is withTenantTx returning one value.
func withTenantTx1[T any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx, tenantID int64) (T, error)) (T, error) {
	return database.WithTenantTx1(ctx, db, fn)
}

// withSystemTx runs fn without RLS (tenants catalog and other global tables).
func withSystemTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	return database.WithSystemTx(ctx, db, fn)
}

// withSystemTx1 is withSystemTx returning one value.
func withSystemTx1[T any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) (T, error)) (T, error) {
	return database.WithSystemTx1(ctx, db, fn)
}

// withExplicitTenantTx runs fn for an explicit tenant (super-admin / login bootstrap).
func withExplicitTenantTx(ctx context.Context, db *sql.DB, tenantID int64, fn func(tx *sql.Tx) error) error {
	return database.WithExplicitTenantTx(ctx, db, tenantID, fn)
}

// withExplicitTenantTx1 is withExplicitTenantTx returning one value.
func withExplicitTenantTx1[T any](ctx context.Context, db *sql.DB, tenantID int64, fn func(tx *sql.Tx) (T, error)) (T, error) {
	return database.WithExplicitTenantTx1(ctx, db, tenantID, fn)
}
