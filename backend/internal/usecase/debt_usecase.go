package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type DebtUsecase struct {
	repo *postgres.DebtRepository
}

func NewDebtUsecase(repo *postgres.DebtRepository) *DebtUsecase {
	return &DebtUsecase{repo: repo}
}

func (u *DebtUsecase) ListAll(ctx context.Context) ([]postgres.Debt, error) {
	return u.repo.ListAll(ctx)
}

func (u *DebtUsecase) GetByID(ctx context.Context, id int64) (*postgres.Debt, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *DebtUsecase) Create(ctx context.Context, d *postgres.Debt) error {
	if d.CreditorName == "" {
		return errors.New("nama kreditur wajib diisi")
	}
	if d.Amount <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, d)
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func (u *DebtUsecase) Update(ctx context.Context, d *postgres.Debt) error {
	return u.repo.Update(ctx, d)
}

func (u *DebtUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
