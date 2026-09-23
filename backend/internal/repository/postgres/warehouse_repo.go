package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Warehouse struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	BranchID   *int64    `json:"branch_id"`
	BranchName string    `json:"branch_name"`
	Address    string    `json:"address"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) List(ctx context.Context) ([]Warehouse, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Warehouse, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT w.id, w.code, w.name, w.branch_id, COALESCE(b.name, ''), w.address, w.is_active, w.created_at, w.updated_at
			FROM warehouses w
			LEFT JOIN branches b ON b.id = w.branch_id AND b.tenant_id = w.tenant_id
			WHERE w.tenant_id=$1
			ORDER BY w.id`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]Warehouse, 0)
		for rows.Next() {
			var w Warehouse
			if err := rows.Scan(&w.ID, &w.Code, &w.Name, &w.BranchID, &w.BranchName, &w.Address, &w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, w)
		}
		return list, nil
	})
}

func (r *WarehouseRepository) GetByID(ctx context.Context, id int64) (*Warehouse, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Warehouse, error) {
		w := &Warehouse{}
		err := tx.QueryRowContext(ctx, `
			SELECT w.id, w.code, w.name, w.branch_id, COALESCE(b.name, ''), w.address, w.is_active, w.created_at, w.updated_at
			FROM warehouses w
			LEFT JOIN branches b ON b.id = w.branch_id AND b.tenant_id = w.tenant_id
			WHERE w.id=$1 AND w.tenant_id=$2`, id, tenantID).
			Scan(&w.ID, &w.Code, &w.Name, &w.BranchID, &w.BranchName, &w.Address, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return w, nil
	})
}

func (r *WarehouseRepository) Create(ctx context.Context, w *Warehouse) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		err := tx.QueryRowContext(ctx,
			"INSERT INTO warehouses (code, name, branch_id, address, tenant_id) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at",
			w.Code, w.Name, w.BranchID, w.Address, tenantID).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WarehouseRepository) Update(ctx context.Context, w *Warehouse) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE warehouses SET code=$1, name=$2, branch_id=$3, address=$4, is_active=$5, updated_at=NOW() WHERE id=$6 AND tenant_id=$7",
			w.Code, w.Name, w.BranchID, w.Address, w.IsActive, w.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WarehouseRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM warehouses WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

type ItemStock struct {
	WarehouseID   int64   `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Quantity      float64 `json:"quantity"`
}

type ItemVariant struct {
	ID              int64   `json:"id"`
	ProductID       int64   `json:"product_id"`
	Name            string  `json:"name"`
	Barcode         string  `json:"barcode"`
	AdditionalPrice float64 `json:"additional_price"`
	Stock           int     `json:"stock"`
	IsActive        bool    `json:"is_active"`
}

type Item struct {
	ID               int64         `json:"id"`
	Code             string        `json:"code"`
	Name             string        `json:"name"`
	Barcode          string        `json:"barcode"`
	Description      string        `json:"description"`
	CategoryID       *int64        `json:"category_id"`
	CategoryName     string        `json:"category_name"`
	PurchasePrice    float64       `json:"purchase_price"`
	BasePrice        float64       `json:"base_price"`
	Unit             string        `json:"unit"`
	MinStock         int           `json:"min_stock"`
	IsService        bool          `json:"is_service"`
	HasVariants      bool          `json:"has_variants"`
	TaxRate          float64       `json:"tax_rate"`
	Stock            float64       `json:"stock"`
	InitialStock     float64       `json:"initial_stock"`
	ImageURL         string        `json:"image_url"`
	Variants         []ItemVariant `json:"variants"`
	StockByWarehouse []ItemStock   `json:"stock_by_warehouse"`
}

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) List(ctx context.Context) ([]Item, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Item, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT p.id, p.code, p.name, p.barcode, p.description, p.category_id, COALESCE(c.name, ''),
				p.purchase_price, p.base_price, p.unit, p.min_stock, p.is_service, p.has_variants, p.tax_rate,
				COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0) AS total_stock,
				p.image_url
			FROM products p
			LEFT JOIN categories c ON c.id = p.category_id AND c.tenant_id = p.tenant_id
			WHERE p.tenant_id=$1
			ORDER BY p.name`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]Item, 0)
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Code, &it.Name, &it.Barcode, &it.Description, &it.CategoryID, &it.CategoryName,
				&it.PurchasePrice, &it.BasePrice, &it.Unit, &it.MinStock, &it.IsService, &it.HasVariants, &it.TaxRate,
				&it.Stock, &it.ImageURL); err != nil {
				return nil, err
			}
			list = append(list, it)
		}
		return list, nil
	})
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*Item, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Item, error) {
		it := &Item{}
		err := tx.QueryRowContext(ctx, `
			SELECT p.id, p.code, p.name, p.barcode, p.description, p.category_id, COALESCE(c.name, ''),
				p.purchase_price, p.base_price, p.unit, p.min_stock, p.is_service, p.has_variants, p.tax_rate,
				COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id = p.id AND s.tenant_id = p.tenant_id), 0) AS total_stock,
				p.image_url
			FROM products p
			LEFT JOIN categories c ON c.id = p.category_id AND c.tenant_id = p.tenant_id
			WHERE p.id=$1 AND p.tenant_id=$2`, id, tenantID).
			Scan(&it.ID, &it.Code, &it.Name, &it.Barcode, &it.Description, &it.CategoryID, &it.CategoryName,
				&it.PurchasePrice, &it.BasePrice, &it.Unit, &it.MinStock, &it.IsService, &it.HasVariants, &it.TaxRate,
				&it.Stock, &it.ImageURL)
		if err != nil {
			return nil, err
		}
		it.Variants, _ = r.ListVariants(ctx, id)
		it.StockByWarehouse, _ = r.ListStockByWarehouse(ctx, id)
		return it, nil
	})
}

func (r *ItemRepository) ListVariants(ctx context.Context, productID int64) ([]ItemVariant, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ItemVariant, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, product_id, name, barcode, additional_price, stock, is_active FROM item_variants WHERE product_id=$1 AND tenant_id=$2 ORDER BY id", productID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]ItemVariant, 0)
		for rows.Next() {
			var v ItemVariant
			if err := rows.Scan(&v.ID, &v.ProductID, &v.Name, &v.Barcode, &v.AdditionalPrice, &v.Stock, &v.IsActive); err != nil {
				return nil, err
			}
			list = append(list, v)
		}
		return list, nil
	})
}

func (r *ItemRepository) ListStockByWarehouse(ctx context.Context, productID int64) ([]ItemStock, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ItemStock, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT s.warehouse_id, w.name, s.quantity FROM stock s
			JOIN warehouses w ON w.id = s.warehouse_id AND w.tenant_id = s.tenant_id
			WHERE s.product_id=$1 AND s.tenant_id=$2 ORDER BY w.id`, productID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]ItemStock, 0)
		for rows.Next() {
			var is ItemStock
			if err := rows.Scan(&is.WarehouseID, &is.WarehouseName, &is.Quantity); err != nil {
				return nil, err
			}
			list = append(list, is)
		}
		return list, nil
	})
}

func (r *ItemRepository) ListStockByWarehouseID(ctx context.Context, warehouseID int64) ([]ItemStock, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]ItemStock, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT s.warehouse_id, w.name, s.quantity FROM stock s
			JOIN warehouses w ON w.id = s.warehouse_id AND w.tenant_id = s.tenant_id
			WHERE s.warehouse_id=$1 AND s.tenant_id=$2 AND s.quantity <> 0 ORDER BY w.id`, warehouseID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]ItemStock, 0)
		for rows.Next() {
			var is ItemStock
			if err := rows.Scan(&is.WarehouseID, &is.WarehouseName, &is.Quantity); err != nil {
				return nil, err
			}
			list = append(list, is)
		}
		return list, nil
	})
}

func (r *ItemRepository) ReplaceVariants(ctx context.Context, productID int64, variants []ItemVariant) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {

		if _, err := tx.ExecContext(ctx, "DELETE FROM item_variants WHERE product_id=$1 AND tenant_id=$2", productID, tenantID); err != nil {
			return err
		}
		for _, v := range variants {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO item_variants (product_id, name, barcode, additional_price, stock, is_active, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7)",
				productID, v.Name, v.Barcode, v.AdditionalPrice, v.Stock, v.IsActive, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}
