package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type CommissionUsecase struct {
	repo *postgres.CommissionRepository
}

func NewCommissionUsecase(repo *postgres.CommissionRepository) *CommissionUsecase {
	return &CommissionUsecase{repo: repo}
}

func (u *CommissionUsecase) ListAll(ctx context.Context) ([]postgres.Commission, error) {
	return u.repo.ListAll(ctx)
}

func (u *CommissionUsecase) GetByID(ctx context.Context, id int64) (*postgres.Commission, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CommissionUsecase) Create(ctx context.Context, c *postgres.Commission) error {
	if c.EmployeeID == 0 {
		return errors.New("karyawan wajib dipilih")
	}
	if c.Rate <= 0 {
		return errors.New("rate komisi harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, c)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (u *CommissionUsecase) Update(ctx context.Context, c *postgres.Commission) error {
	return u.repo.Update(ctx, c)
}

func (u *CommissionUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
