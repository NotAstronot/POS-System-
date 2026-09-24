package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Product struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	BasePrice     float64   `json:"base_price"`
	Unit          string    `json:"unit"`
	ImageURL      string    `json:"image_url"`
	CategoryID    *int64    `json:"category_id"`
	HasVariants   bool      `json:"has_variants"`
	TaxRate       float64   `json:"tax_rate"`
	SortOrder     int       `json:"sort_order"`
	Stock         int       `json:"stock"`
	Barcode       string    `json:"barcode"`
	PurchasePrice float64   `json:"purchase_price"`
	MinStock      int       `json:"min_stock"`
	IsService     bool      `json:"is_service"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

const productCols = "id, code, name, description, base_price, unit, image_url, category_id, has_variants, tax_rate, sort_order, stock, barcode, purchase_price, min_stock, is_service, created_at, updated_at"

func scanProduct(p *Product, row scanRow) error {
	return row.Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.Barcode, &p.PurchasePrice, &p.MinStock, &p.IsService, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) List(ctx context.Context) ([]Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Product, error) {
		rows, err := tx.QueryContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.tenant_id=$1 ORDER BY p.sort_order, p.name`, tenantID)
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
		return list, rows.Err()
	})
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Product, error) {
		p := &Product{}
		err := tx.QueryRowContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.id=$1 AND p.tenant_id=$2`, id, tenantID).Scan(
			&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.Barcode, &p.PurchasePrice, &p.MinStock, &p.IsService, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return p, nil
	})
}

func (r *ProductRepository) Search(ctx context.Context, q string) ([]Product, error) {
	like := "%" + q + "%"
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Product, error) {
		rows, err := tx.QueryContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.tenant_id=$1 AND (p.name ILIKE $2 OR p.code ILIKE $2) ORDER BY p.sort_order, p.name`, tenantID, like)
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
		return list, rows.Err()
	})
}

func (r *ProductRepository) Barcode(ctx context.Context, code string) (*Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Product, error) {
		p := &Product{}
		err := tx.QueryRowContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.code=$1 AND p.tenant_id=$2`, code, tenantID).Scan(
			&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.Barcode, &p.PurchasePrice, &p.MinStock, &p.IsService, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return p, nil
	})
}

func (r *ProductRepository) ListByCategory(ctx context.Context, categoryID int64) ([]Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Product, error) {
		rows, err := tx.QueryContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.category_id=$1 AND p.tenant_id=$2 ORDER BY p.sort_order, p.name`, categoryID, tenantID)
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
		return list, rows.Err()
	})
}

func (r *ProductRepository) Create(ctx context.Context, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Product, error) {
		p := &Product{}
		err := tx.QueryRowContext(ctx,
			"INSERT INTO products (tenant_id, code, name, description, base_price, unit, image_url, category_id, has_variants, tax_rate, sort_order, stock, barcode, purchase_price, min_stock, is_service) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING "+productCols,
			tenantID, code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock, barcode, purchasePrice, minStock, isService).Scan(
			&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.Barcode, &p.PurchasePrice, &p.MinStock, &p.IsService, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if stock > 0 {
			if err := r.syncStockTableTx(tx, p.ID, stock, tenantID); err != nil {
				return nil, err
			}
		}
		return p, nil
	})
}

func (r *ProductRepository) Update(ctx context.Context, id int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Product, error) {
		p := &Product{}
		err := tx.QueryRowContext(ctx,
			"UPDATE products SET code=$1, name=$2, description=$3, base_price=$4, unit=$5, image_url=$6, category_id=$7, has_variants=$8, tax_rate=$9, sort_order=$10, stock=$11, barcode=$12, purchase_price=$13, min_stock=$14, is_service=$15, updated_at=NOW() WHERE id=$16 AND tenant_id=$17 RETURNING "+productCols,
			code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock, barcode, purchasePrice, minStock, isService, id, tenantID).Scan(
			&p.ID, &p.Code, &p.Name, &p.Description, &p.BasePrice, &p.Unit, &p.ImageURL, &p.CategoryID, &p.HasVariants, &p.TaxRate, &p.SortOrder, &p.Stock, &p.Barcode, &p.PurchasePrice, &p.MinStock, &p.IsService, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return p, nil
	})
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM products WHERE id=$1 AND tenant_id=$2", id, tenantID)
		return err
	})
}

func (r *ProductRepository) syncStockTableTx(tx *sql.Tx, productID int64, stock int, tenantID int64) error {
	_, err := tx.ExecContext(context.Background(),
		`INSERT INTO stock (product_id, warehouse_id, quantity, tenant_id, updated_at)
		 SELECT $1, id, $2, $3, NOW() FROM warehouses WHERE is_active = true AND tenant_id = $3 ORDER BY id LIMIT 1`,
		productID, stock, tenantID)
	return err
}

func (r *ProductRepository) SyncStockTable(ctx context.Context, productID int64, stock int) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		return r.syncStockTableTx(tx, productID, stock, tenantID)
	})
}

func (r *ProductRepository) SetAvailability(ctx context.Context, productID int64, date string, isAvailable bool) (*ProductAvailability, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*ProductAvailability, error) {
		a := &ProductAvailability{}
		err := tx.QueryRowContext(ctx,
			"INSERT INTO product_availability (product_id, date, is_available, tenant_id) VALUES ($1,$2,$3,$4) ON CONFLICT (product_id, date) DO UPDATE SET is_available=$3 RETURNING id, product_id, date, is_available, created_at",
			productID, date, isAvailable, tenantID).Scan(&a.ID, &a.ProductID, &a.Date, &a.IsAvailable, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		return a, nil
	})
}

func (r *ProductRepository) GetAvailability(ctx context.Context, productID int64, date string) (*ProductAvailability, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*ProductAvailability, error) {
		a := &ProductAvailability{}
		err := tx.QueryRowContext(ctx, "SELECT id, product_id, date, is_available, created_at FROM product_availability WHERE product_id=$1 AND date=$2 AND tenant_id=$3", productID, date, tenantID).Scan(&a.ID, &a.ProductID, &a.Date, &a.IsAvailable, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		return a, nil
	})
}

func (r *ProductRepository) ListAvailability(ctx context.Context, date string) ([]ProductAvailability, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ProductAvailability, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, product_id, date, is_available, created_at FROM product_availability WHERE date=$1 AND tenant_id=$2", date, tenantID)
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
		return list, rows.Err()
	})
}

func (r *ProductRepository) ListLowStock(ctx context.Context, threshold int) ([]Product, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Product, error) {
		rows, err := tx.QueryContext(ctx, `SELECT p.id, p.code, p.name, p.description, p.base_price, p.unit, p.image_url, p.category_id, p.has_variants, p.tax_rate, p.sort_order,
			COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0)::int AS stock,
			p.barcode, p.purchase_price, p.min_stock, p.is_service, p.created_at, p.updated_at
			FROM products p WHERE p.tenant_id=$1 AND COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0) <= $2`, tenantID, threshold)
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
		return list, rows.Err()
	})
}
