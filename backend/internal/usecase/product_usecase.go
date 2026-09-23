package usecase

import (
	"context"
	"pos-system/internal/repository/postgres"
)

type ProductUsecase struct {
	repo *postgres.ProductRepository
}

func NewProductUsecase(repo *postgres.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) List(ctx context.Context) ([]postgres.Product, error) {
	return u.repo.List(ctx)
}

func (u *ProductUsecase) GetByID(ctx context.Context, id int64) (*postgres.Product, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ProductUsecase) Search(ctx context.Context, q string) ([]postgres.Product, error) {
	return u.repo.Search(ctx, q)
}

func (u *ProductUsecase) Barcode(ctx context.Context, code string) (*postgres.Product, error) {
	return u.repo.Barcode(ctx, code)
}

func (u *ProductUsecase) ListByCategory(ctx context.Context, categoryID int64) ([]postgres.Product, error) {
	return u.repo.ListByCategory(ctx, categoryID)
}

func (u *ProductUsecase) Create(ctx context.Context, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*postgres.Product, error) {
	return u.repo.Create(ctx, code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock, barcode, purchasePrice, minStock, isService)
}

func (u *ProductUsecase) Update(ctx context.Context, id int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int, barcode string, purchasePrice float64, minStock int, isService bool) (*postgres.Product, error) {
	return u.repo.Update(ctx, id, code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock, barcode, purchasePrice, minStock, isService)
}

func (u *ProductUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) SetAvailability(ctx context.Context, productID int64, date string, isAvailable bool) (*postgres.ProductAvailability, error) {
	return u.repo.SetAvailability(ctx, productID, date, isAvailable)
}

func (u *ProductUsecase) GetAvailability(ctx context.Context, productID int64, date string) (*postgres.ProductAvailability, error) {
	return u.repo.GetAvailability(ctx, productID, date)
}

func (u *ProductUsecase) ListAvailability(ctx context.Context, date string) ([]postgres.ProductAvailability, error) {
	return u.repo.ListAvailability(ctx, date)
}

func (u *ProductUsecase) ListLowStock(ctx context.Context, threshold int) ([]postgres.Product, error) {
	return u.repo.ListLowStock(ctx, threshold)
}
