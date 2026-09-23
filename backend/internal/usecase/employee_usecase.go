package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type EmployeeUsecase struct {
	repo *postgres.EmployeeRepository
}

func NewEmployeeUsecase(repo *postgres.EmployeeRepository) *EmployeeUsecase {
	return &EmployeeUsecase{repo: repo}
}

func (u *EmployeeUsecase) ListAll(ctx context.Context) ([]postgres.Employee, error) {
	return u.repo.ListAll(ctx)
}

func (u *EmployeeUsecase) GetByID(ctx context.Context, id int64) (*postgres.Employee, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *EmployeeUsecase) Create(ctx context.Context, e *postgres.Employee) error {
	if e.NIK == "" {
		return errors.New("NIK wajib diisi")
	}
	if e.Name == "" {
		return errors.New("nama karyawan wajib diisi")
	}
	id, err := u.repo.Create(ctx, e)
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (u *EmployeeUsecase) Update(ctx context.Context, e *postgres.Employee) error {
	if e.Name == "" {
		return errors.New("nama karyawan wajib diisi")
	}
	return u.repo.Update(ctx, e)
}

func (u *EmployeeUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
