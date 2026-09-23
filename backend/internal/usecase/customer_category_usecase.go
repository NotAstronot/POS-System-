package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type CustomerCategoryUsecase struct {
	repo *postgres.CustomerCategoryRepository
}

func NewCustomerCategoryUsecase(repo *postgres.CustomerCategoryRepository) *CustomerCategoryUsecase {
	return &CustomerCategoryUsecase{repo: repo}
}

func (u *CustomerCategoryUsecase) ListAll(ctx context.Context) ([]postgres.CustomerCategory, error) {
	return u.repo.ListAll(ctx)
}

func (u *CustomerCategoryUsecase) GetByID(ctx context.Context, id int64) (*postgres.CustomerCategory, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CustomerCategoryUsecase) Create(ctx context.Context, c *postgres.CustomerCategory) error {
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

func (u *CustomerCategoryUsecase) Update(ctx context.Context, c *postgres.CustomerCategory) error {
	return u.repo.Update(ctx, c)
}

func (u *CustomerCategoryUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
