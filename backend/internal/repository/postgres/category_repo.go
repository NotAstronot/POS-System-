package postgres

import (
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

func (r *CategoryRepository) List() ([]Category, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM categories ORDER BY id")
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
	return list, nil
}

func (r *CategoryRepository) GetByID(id int64) (*Category, error) {
	c := &Category{}
	err := r.db.QueryRow("SELECT id, name, created_at, updated_at FROM categories WHERE id=$1", id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoryRepository) Create(name string) (*Category, error) {
	c := &Category{}
	err := r.db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at", name).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoryRepository) Update(id int64, name string) (*Category, error) {
	c := &Category{}
	err := r.db.QueryRow("UPDATE categories SET name=$1, updated_at=NOW() WHERE id=$2 RETURNING id, name, created_at, updated_at", name, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoryRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id=$1", id)
	return err
}
