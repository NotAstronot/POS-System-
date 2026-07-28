package postgres

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	Unit        string    `json:"unit"`
	ImageURL    string    `json:"image_url"`
	CategoryID  *int64    `json:"category_id"`
	HasVariants bool      `json:"has_variants"`
	TaxRate     float64   `json:"tax_rate"`
	SortOrder   int       `json:"sort_order"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductAvailability struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	Date        string    `json:"date"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

const productCols = "id, code, name, description, base_price, unit, image_url, category_id, has_variants, tax_rate, sort_order, stock, created_at, updated_at"

func scanProduct(p *Product, row scanRow) error {
	return row.Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) List() ([]Product, error) {
	rows, err := r.db.Query("SELECT "+productCols+" FROM products ORDER BY sort_order, name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Product
	for rows.Next() {
		var p Product
		if err := scanProduct(&p, rows); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ProductRepository) GetByID(id int64) (*Product, error) {
	p := &Product{}
	err := scanProduct(p, r.db.QueryRow("SELECT "+productCols+" FROM products WHERE id=$1", id))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) Search(q string) ([]Product, error) {
	like := "%" + q + "%"
	rows, err := r.db.Query("SELECT "+productCols+" FROM products WHERE name ILIKE $1 OR code ILIKE $1 ORDER BY sort_order, name", like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Product
	for rows.Next() {
		var p Product
		if err := scanProduct(&p, rows); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ProductRepository) Barcode(code string) (*Product, error) {
	p := &Product{}
	err := scanProduct(p, r.db.QueryRow("SELECT "+productCols+" FROM products WHERE code=$1", code))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) ListByCategory(categoryID int64) ([]Product, error) {
	rows, err := r.db.Query("SELECT "+productCols+" FROM products WHERE category_id=$1 ORDER BY sort_order, name", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Product
	for rows.Next() {
		var p Product
		if err := scanProduct(&p, rows); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ProductRepository) Create(code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int) (*Product, error) {
	p := &Product{}
	err := scanProduct(p, r.db.QueryRow(
		"INSERT INTO products (code, name, description, base_price, unit, image_url, category_id, has_variants, tax_rate, sort_order, stock) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "+productCols,
		code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) Update(id int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int) (*Product, error) {
	p := &Product{}
	err := scanProduct(p, r.db.QueryRow(
		"UPDATE products SET code=$1, name=$2, description=$3, base_price=$4, unit=$5, image_url=$6, category_id=$7, has_variants=$8, tax_rate=$9, sort_order=$10, stock=$11, updated_at=NOW() WHERE id=$12 RETURNING "+productCols,
		code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock, id))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}

func (r *ProductRepository) SetAvailability(productID int64, date string, isAvailable bool) (*ProductAvailability, error) {
	a := &ProductAvailability{}
	err := r.db.QueryRow(
		"INSERT INTO product_availability (product_id, date, is_available) VALUES ($1,$2,$3) ON CONFLICT (product_id, date) DO UPDATE SET is_available=$3 RETURNING id, product_id, date, is_available, created_at",
		productID, date, isAvailable).Scan(&a.ID, &a.ProductID, &a.Date, &a.IsAvailable, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *ProductRepository) GetAvailability(productID int64, date string) (*ProductAvailability, error) {
	a := &ProductAvailability{}
	err := r.db.QueryRow("SELECT id, product_id, date, is_available, created_at FROM product_availability WHERE product_id=$1 AND date=$2", productID, date).Scan(&a.ID, &a.ProductID, &a.Date, &a.IsAvailable, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *ProductRepository) ListAvailability(date string) ([]ProductAvailability, error) {
	rows, err := r.db.Query("SELECT id, product_id, date, is_available, created_at FROM product_availability WHERE date=$1", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ProductAvailability
	for rows.Next() {
		var a ProductAvailability
		if err := rows.Scan(&a.ID, &a.ProductID, &a.Date, &a.IsAvailable, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}
