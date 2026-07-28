package usecase

import "pos-system/internal/repository/postgres"

type ProductUsecase struct {
	repo *postgres.ProductRepository
}

func NewProductUsecase(repo *postgres.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) List() ([]postgres.Product, error) {
	return u.repo.List()
}

func (u *ProductUsecase) GetByID(id int64) (*postgres.Product, error) {
	return u.repo.GetByID(id)
}

func (u *ProductUsecase) Search(q string) ([]postgres.Product, error) {
	return u.repo.Search(q)
}

func (u *ProductUsecase) Barcode(code string) (*postgres.Product, error) {
	return u.repo.Barcode(code)
}

func (u *ProductUsecase) ListByCategory(categoryID int64) ([]postgres.Product, error) {
	return u.repo.ListByCategory(categoryID)
}

func (u *ProductUsecase) Create(code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int) (*postgres.Product, error) {
	return u.repo.Create(code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock)
}

func (u *ProductUsecase) Update(id int64, code, name, description string, basePrice float64, unit, imageURL string, categoryID *int64, hasVariants bool, taxRate float64, sortOrder, stock int) (*postgres.Product, error) {
	return u.repo.Update(id, code, name, description, basePrice, unit, imageURL, categoryID, hasVariants, taxRate, sortOrder, stock)
}

func (u *ProductUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}

func (u *ProductUsecase) SetAvailability(productID int64, date string, isAvailable bool) (*postgres.ProductAvailability, error) {
	return u.repo.SetAvailability(productID, date, isAvailable)
}

func (u *ProductUsecase) GetAvailability(productID int64, date string) (*postgres.ProductAvailability, error) {
	return u.repo.GetAvailability(productID, date)
}

func (u *ProductUsecase) ListAvailability(date string) ([]postgres.ProductAvailability, error) {
	return u.repo.ListAvailability(date)
}
