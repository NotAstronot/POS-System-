package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type ChartOfAccountsUsecase struct {
	repo *postgres.ChartOfAccountRepository
}

func NewChartOfAccountsUsecase(repo *postgres.ChartOfAccountRepository) *ChartOfAccountsUsecase {
	return &ChartOfAccountsUsecase{repo: repo}
}

func (u *ChartOfAccountsUsecase) ListAll(ctx context.Context) ([]postgres.ChartOfAccount, error) {
	return u.repo.ListAll(ctx)
}

func (u *ChartOfAccountsUsecase) GetByID(ctx context.Context, id int64) (*postgres.ChartOfAccount, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ChartOfAccountsUsecase) Create(ctx context.Context, c *postgres.ChartOfAccount) error {
	switch c.Category {
	case "aset", "kewajiban", "ekuitas", "pendapatan", "beban":
	default:
		return errors.New("kategori akun tidak valid")
	}
	if c.NormalBalance != "debit" && c.NormalBalance != "kredit" {
		return errors.New("saldo normal harus debit atau kredit")
	}
	id, err := u.repo.Create(ctx, c)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (u *ChartOfAccountsUsecase) Update(ctx context.Context, c *postgres.ChartOfAccount) error {
	switch c.Category {
	case "aset", "kewajiban", "ekuitas", "pendapatan", "beban":
	default:
		return errors.New("kategori akun tidak valid")
	}
	if c.NormalBalance != "debit" && c.NormalBalance != "kredit" {
		return errors.New("saldo normal harus debit atau kredit")
	}
	return u.repo.Update(ctx, c)
}

func (u *ChartOfAccountsUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
