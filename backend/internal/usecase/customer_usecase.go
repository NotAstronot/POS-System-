package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type CustomerUsecase struct {
	repo *postgres.CustomerRepository
}

func NewCustomerUsecase(repo *postgres.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (u *CustomerUsecase) ListAll(ctx context.Context) ([]postgres.Customer, error) {
	return u.repo.ListAll(ctx)
}

func (u *CustomerUsecase) GetByID(ctx context.Context, id int64) (*postgres.Customer, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CustomerUsecase) Create(ctx context.Context, c *postgres.Customer) error {
	if c.Name == "" {
		return errors.New("nama pelanggan wajib diisi")
	}
	id, err := u.repo.Create(ctx, c)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (u *CustomerUsecase) Update(ctx context.Context, c *postgres.Customer) error {
	return u.repo.Update(ctx, c)
}

func (u *CustomerUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
