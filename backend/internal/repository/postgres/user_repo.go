package postgres

import (
	"database/sql"
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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const userCols = "id, username, name, password, role, outlet_id, tenant_id, created_at, updated_at"

func scanUser(u *User, row scanRow) error {
	return row.Scan(&u.ID, &u.Username, &u.Name, &u.Password, &u.Role, &u.OutletID, &u.TenantID, &u.CreatedAt, &u.UpdatedAt)
}

type scanRow interface {
	Scan(dest ...any) error
}

func (r *UserRepository) FindByUsername(username string) (*User, error) {
	u := &User{}
	err := scanUser(u, r.db.QueryRow("SELECT "+userCols+" FROM users WHERE username=$1", username))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(id int64) (*User, error) {
	u := &User{}
	err := scanUser(u, r.db.QueryRow("SELECT "+userCols+" FROM users WHERE id=$1", id))
	if err != nil {
		return nil, err
	}
	return u, nil
}
