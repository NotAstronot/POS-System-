package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type CashUsecase struct {
	repo *postgres.CashTransactionRepository
}

func NewCashUsecase(repo *postgres.CashTransactionRepository) *CashUsecase {
	return &CashUsecase{repo: repo}
}

func (u *CashUsecase) ListAll(ctx context.Context) ([]postgres.CashTransaction, error) {
	return u.repo.ListAll(ctx)
}

func (u *CashUsecase) GetByID(ctx context.Context, id int64) (*postgres.CashTransaction, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CashUsecase) Create(ctx context.Context, c *postgres.CashTransaction) error {
	if c.Amount <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, c)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (u *CashUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
