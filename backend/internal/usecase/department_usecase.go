package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type DepartmentUsecase struct {
	repo *postgres.DepartmentRepository
}

func NewDepartmentUsecase(repo *postgres.DepartmentRepository) *DepartmentUsecase {
	return &DepartmentUsecase{repo: repo}
}

func (u *DepartmentUsecase) ListAll(ctx context.Context) ([]postgres.Department, error) {
	return u.repo.ListAll(ctx)
}

func (u *DepartmentUsecase) GetByID(ctx context.Context, id int64) (*postgres.Department, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *DepartmentUsecase) Create(ctx context.Context, d *postgres.Department) error {
	if d.Name == "" {
		return errors.New("nama departemen wajib diisi")
	}
	id, err := u.repo.Create(ctx, d)
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func (u *DepartmentUsecase) Update(ctx context.Context, d *postgres.Department) error {
	if d.Name == "" {
		return errors.New("nama departemen wajib diisi")
	}
	return u.repo.Update(ctx, d)
}

func (u *DepartmentUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
