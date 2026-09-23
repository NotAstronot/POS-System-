package repository

import "pos-system/internal/domain/entity"

type ProductRepository interface {
	List(tenantID int64) ([]entity.Product, error)
	GetByID(id int64, tenantID int64) (*entity.Product, error)
	Search(q string, tenantID int64) ([]entity.Product, error)
	Barcode(code string, tenantID int64) (*entity.Product, error)
	ListByCategory(categoryID int64, tenantID int64) ([]entity.Product, error)
	Create(tenantID int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*entity.Product, error)
	Update(id int64, tenantID int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*entity.Product, error)
	Delete(id int64, tenantID int64) error
	SetAvailability(productID int64, date string, isAvailable bool, tenantID int64) (*entity.ProductAvailability, error)
	GetAvailability(productID int64, date string, tenantID int64) (*entity.ProductAvailability, error)
	ListAvailability(date string, tenantID int64) ([]entity.ProductAvailability, error)
	ListLowStock(threshold int, tenantID int64) ([]entity.Product, error)
}
