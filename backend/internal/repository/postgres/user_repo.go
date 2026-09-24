package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	OutletID  string    `json:"outlet_id"`
	TenantID  string    `json:"tenant_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permission struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const userCols = "id, username, name, password, role, outlet_id, tenant_id, is_active, created_at, updated_at"

func scanUser(u *User, row scanRow) error {
	return row.Scan(&u.ID, &u.Username, &u.Name, &u.Password, &u.Role, &u.OutletID, &u.TenantID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
}

type scanRow interface {
	Scan(dest ...any) error
}

// FindByUsername finds a user by username within the current tenant context.
// users has no RLS, so tenant_id is filtered explicitly from ctx.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*User, error) {
		u := &User{}
		err := scanUser(u, tx.QueryRowContext(ctx,
			"SELECT "+userCols+" FROM users WHERE username=$1 AND tenant_id=$2::text", username, tenantID))
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}

// FindByUsernameGlobal searches for a user by username across ALL tenants (no RLS).
// Used during login when tenant is not yet known.
func (r *UserRepository) FindByUsernameGlobal(ctx context.Context, username string) (*User, error) {
	return withSystemTx1(ctx, r.db, func(tx *sql.Tx) (*User, error) {
		u := &User{}
		err := scanUser(u, tx.QueryRowContext(ctx, "SELECT "+userCols+" FROM users WHERE username=$1", username))
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*User, error) {
		u := &User{}
		err := scanUser(u, tx.QueryRowContext(ctx,
			"SELECT "+userCols+" FROM users WHERE id=$1 AND tenant_id=$2::text", id, tenantID))
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}

func (r *UserRepository) ListAll(ctx context.Context) ([]User, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]User, error) {
		rows, err := tx.QueryContext(ctx,
			"SELECT "+userCols+" FROM users WHERE tenant_id=$1::text ORDER BY id ASC", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Password, &u.Role, &u.OutletID, &u.TenantID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
				return nil, err
			}
			users = append(users, u)
		}
		return users, rows.Err()
	})
}

func (r *UserRepository) Create(ctx context.Context, u *User) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO users (username, name, password, role, outlet_id, tenant_id) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
			u.Username, u.Name, u.Password, u.Role, u.OutletID, tenantID,
		).Scan(&id)
		return id, err
	})
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE users SET name=$1, role=$2, outlet_id=$3, tenant_id=$4::text, is_active=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7::text",
			u.Name, u.Role, u.OutletID, tenantID, u.IsActive, u.ID, tenantID,
		)
		return err
	})
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE users SET password=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3::text",
			password, id, tenantID)
		return err
	})
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"DELETE FROM users WHERE id=$1 AND tenant_id=$2::text", id, tenantID)
		return err
	})
}

// GetPermissions returns permissions for a user within the current tenant.
// user_permissions/permissions have no tenant_id, so membership is scoped via users.
func (r *UserRepository) GetPermissions(ctx context.Context, userID int64) ([]Permission, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Permission, error) {
		rows, err := tx.QueryContext(ctx,
			`SELECT p.id, p.name, p.description
			 FROM permissions p
			 INNER JOIN user_permissions up ON p.id=up.permission_id
			 INNER JOIN users u ON u.id=up.user_id
			 WHERE up.user_id=$1 AND u.tenant_id=$2::text
			 ORDER BY p.id`,
			userID, tenantID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var perms []Permission
		for rows.Next() {
			var p Permission
			if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
				return nil, err
			}
			perms = append(perms, p)
		}
		return perms, rows.Err()
	})
}

// GetPermissionsForTenant returns permissions for a user in an explicit tenant (super-admin / login bootstrap).
func (r *UserRepository) GetPermissionsForTenant(ctx context.Context, tenantID int64, userID int64) ([]Permission, error) {
	if tenantID <= 0 {
		return nil, errors.New("tenant_id is required")
	}
	return withExplicitTenantTx1(ctx, r.db, tenantID, func(tx *sql.Tx) ([]Permission, error) {
		rows, err := tx.QueryContext(ctx,
			`SELECT p.id, p.name, p.description FROM permissions p INNER JOIN user_permissions up ON p.id=up.permission_id WHERE up.user_id=$1 ORDER BY p.id`,
			userID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var perms []Permission
		for rows.Next() {
			var p Permission
			if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
				return nil, err
			}
			perms = append(perms, p)
		}
		return perms, rows.Err()
	})
}

func (r *UserRepository) SetPermissions(ctx context.Context, userID int64, permissionIDs []int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		var ownerTenant string
		err := tx.QueryRowContext(ctx,
			"SELECT tenant_id FROM users WHERE id=$1 AND tenant_id=$2::text",
			userID, tenantID).Scan(&ownerTenant)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("user not found in this tenant")
			}
			return err
		}
		if _, err := tx.ExecContext(ctx,
			"DELETE FROM user_permissions WHERE user_id=$1", userID); err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO user_permissions (user_id, permission_id) VALUES ($1,$2)",
				userID, pid); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepository) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Permission, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name, description FROM permissions ORDER BY id")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var perms []Permission
		for rows.Next() {
			var p Permission
			if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
				return nil, err
			}
			perms = append(perms, p)
		}
		return perms, rows.Err()
	})
}
