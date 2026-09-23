package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"time"
)

type PeriodEndUsecase struct {
	repo *postgres.PeriodEndRepository
}

func NewPeriodEndUsecase(repo *postgres.PeriodEndRepository) *PeriodEndUsecase {
	return &PeriodEndUsecase{repo: repo}
}

func (u *PeriodEndUsecase) ListPeriods(ctx context.Context) ([]postgres.PeriodClosing, error) {
	return u.repo.ListPeriods(ctx)
}

func (u *PeriodEndUsecase) GetPeriod(ctx context.Context, id int64) (*postgres.PeriodClosing, error) {
	return u.repo.GetPeriod(ctx, id)
}

func (u *PeriodEndUsecase) OpenPeriod(ctx context.Context, period string) error {
	if len(period) != 7 || period[4] != '-' {
		return errors.New("format periode harus YYYY-MM")
	}
	if _, err := time.Parse("2006-01", period); err != nil {
		return errors.New("format periode harus YYYY-MM")
	}
	return u.repo.OpenPeriod(ctx, period)
}

func (u *PeriodEndUsecase) ListFx(ctx context.Context, period string) ([]postgres.FxDifference, error) {
	return u.repo.ListFx(ctx, period)
}

func (u *PeriodEndUsecase) CreateFx(ctx context.Context, f *postgres.FxDifference) error {
	if f.Currency == "" {
		return errors.New("mata uang wajib diisi")
	}
	if f.GainAmount < 0 || f.LossAmount < 0 {
		return errors.New("nilai selisih tidak boleh negatif")
	}
	if f.GainAmount == 0 && f.LossAmount == 0 {
		return errors.New("isi nilai keuntungan atau kerugian selisih kurs")
	}
	return u.repo.CreateFx(ctx, f)
}

func (u *PeriodEndUsecase) RunDepreciation(ctx context.Context, id int64) (*postgres.DepreciationResult, error) {
	p, err := u.repo.GetPeriod(ctx, id)
	if err != nil {
		return nil, errors.New("periode belum dibuka")
	}
	return u.repo.RunDepreciation(ctx, p.Period)
}

func (u *PeriodEndUsecase) ClosePeriod(ctx context.Context, id int64) error {
	p, err := u.repo.GetPeriod(ctx, id)
	if err != nil {
		return errors.New("periode belum dibuka")
	}
	return u.repo.ClosePeriod(ctx, p.Period)
}
