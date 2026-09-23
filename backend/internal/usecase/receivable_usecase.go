package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type ReceivableUsecase struct {
	repo *postgres.ReceivableRepository
}

func NewReceivableUsecase(repo *postgres.ReceivableRepository) *ReceivableUsecase {
	return &ReceivableUsecase{repo: repo}
}

func (u *ReceivableUsecase) ListAll(ctx context.Context) ([]postgres.Receivable, error) {
	return u.repo.ListAll(ctx)
}

func (u *ReceivableUsecase) GetByID(ctx context.Context, id int64) (*postgres.Receivable, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ReceivableUsecase) Create(ctx context.Context, rc *postgres.Receivable) error {
	if rc.DebtorName == "" {
		return errors.New(" nama debitur wajib diisi")
	}
	if rc.Amount <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, rc)
	if err != nil {
		return err
	}
	rc.ID = id
	return nil
}

func (u *ReceivableUsecase) Update(ctx context.Context, rc *postgres.Receivable) error {
	return u.repo.Update(ctx, rc)
}

func (u *ReceivableUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
