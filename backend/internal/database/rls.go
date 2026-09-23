package database

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
)

// ErrTenantRequired is returned when a tenant-scoped operation runs without a tenant.
var ErrTenantRequired = errors.New("tenant_id is required for RLS")

// SetRLSTenant sets the PostgreSQL session variable for RLS within a transaction.
// Must be called after BEGIN and before any queries.
// Uses set_config(..., is_local=true) so the value resets on COMMIT/ROLLBACK.
func SetRLSTenant(tx *sql.Tx, tenantID int64) error {
	if tenantID <= 0 {
		return ErrTenantRequired
	}
	_, err := tx.ExecContext(context.Background(),
		`SELECT set_config('app.current_tenant', $1, true)`,
		strconv.FormatInt(tenantID, 10),
	)
	return err
}

// InjectRLSSession pins a dedicated connection from the pool and sets
// app.current_tenant at session scope (is_local=false).
// Caller must Close() the returned connection to release it back to the pool
// (and clear the session GUC).
func InjectRLSSession(ctx context.Context, db *sql.DB, tenantID int64) (*sql.Conn, error) {
	if tenantID <= 0 {
		return nil, ErrTenantRequired
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := conn.ExecContext(ctx,
		`SELECT set_config('app.current_tenant', $1, false)`,
		strconv.FormatInt(tenantID, 10),
	); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

type rlsConnKey struct{}

// WithRLSConn stores the request-scoped RLS connection in ctx.
func WithRLSConn(ctx context.Context, conn *sql.Conn) context.Context {
	return context.WithValue(ctx, rlsConnKey{}, conn)
}

// GetRLSConn returns the request-scoped RLS connection if present.
func GetRLSConn(ctx context.Context) (*sql.Conn, bool) {
	conn, ok := ctx.Value(rlsConnKey{}).(*sql.Conn)
	return conn, ok
}

// WithRLSTenant executes a function within a transaction with RLS tenant context set.
// Prefer WithTenantTx / BeginTenantTx which take tenant_id from context.Context.
// This variant is kept for explicit-tenant call sites that pass tenantID directly.
func WithRLSTenant(db *sql.DB, tenantID int64, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
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

// WithRLSTenant1 executes a function within a transaction with RLS tenant context set and returns one value.
func WithRLSTenant1[T any](db *sql.DB, tenantID int64, fn func(*sql.Tx) (T, error)) (T, error) {
	var result T
	err := WithRLSTenant(db, tenantID, func(tx *sql.Tx) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}

// WithRLSTenant2 executes a function within a transaction with RLS tenant context set and returns two values.
func WithRLSTenant2[T1, T2 any](db *sql.DB, tenantID int64, fn func(*sql.Tx) (T1, T2, error)) (T1, T2, error) {
	var result1 T1
	var result2 T2
	err := WithRLSTenant(db, tenantID, func(tx *sql.Tx) error {
		var err error
		result1, result2, err = fn(tx)
		return err
	})
	return result1, result2, err
}
