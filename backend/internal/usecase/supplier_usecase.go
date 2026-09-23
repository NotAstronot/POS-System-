package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SupplierUsecase struct {
	repo *postgres.SupplierRepository
}

func NewSupplierUsecase(repo *postgres.SupplierRepository) *SupplierUsecase {
	return &SupplierUsecase{repo: repo}
}

func (u *SupplierUsecase) ListAll(ctx context.Context) ([]postgres.Supplier, error) {
	return u.repo.ListAll(ctx)
}

func (u *SupplierUsecase) GetByID(ctx context.Context, id int64) (*postgres.Supplier, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SupplierUsecase) Create(ctx context.Context, s *postgres.Supplier) error {
	if s.Name == "" {
		return errors.New("nama supplier wajib diisi")
	}
	id, err := u.repo.Create(ctx, s)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (u *SupplierUsecase) Update(ctx context.Context, s *postgres.Supplier) error {
	return u.repo.Update(ctx, s)
}

func (u *SupplierUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
