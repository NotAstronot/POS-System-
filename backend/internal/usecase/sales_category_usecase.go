package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesCategoryUsecase struct {
	repo *postgres.SalesCategoryRepository
}

func NewSalesCategoryUsecase(repo *postgres.SalesCategoryRepository) *SalesCategoryUsecase {
	return &SalesCategoryUsecase{repo: repo}
}

func (u *SalesCategoryUsecase) ListAll(ctx context.Context) ([]postgres.SalesCategory, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesCategoryUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesCategory, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesCategoryUsecase) Create(ctx context.Context, c *postgres.SalesCategory) error {
	if c.Name == "" {
		return errors.New("nama kategori wajib diisi")
	}
	id, err := u.repo.Create(ctx, c)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (u *SalesCategoryUsecase) Update(ctx context.Context, c *postgres.SalesCategory) error {
	return u.repo.Update(ctx, c)
}

func (u *SalesCategoryUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
