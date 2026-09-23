package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type WithdrawalUsecase struct {
	repo *postgres.WithdrawalRepository
}

func NewWithdrawalUsecase(repo *postgres.WithdrawalRepository) *WithdrawalUsecase {
	return &WithdrawalUsecase{repo: repo}
}

func (u *WithdrawalUsecase) ListAll(ctx context.Context) ([]postgres.Withdrawal, error) {
	return u.repo.ListAll(ctx)
}

func (u *WithdrawalUsecase) GetByID(ctx context.Context, id int64) (*postgres.Withdrawal, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *WithdrawalUsecase) Create(ctx context.Context, w *postgres.Withdrawal) error {
	if w.Amount <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, w)
	if err != nil {
		return err
	}
	w.ID = id
	return nil
}

func (u *WithdrawalUsecase) Update(ctx context.Context, w *postgres.Withdrawal) error {
	return u.repo.Update(ctx, w)
}

func (u *WithdrawalUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *WithdrawalUsecase) Approve(ctx context.Context, id int64, approvedBy int64) error {
	return u.repo.Approve(ctx, id, approvedBy)
}

func (u *WithdrawalUsecase) Reject(ctx context.Context, id int64) error {
	return u.repo.Reject(ctx, id)
}
