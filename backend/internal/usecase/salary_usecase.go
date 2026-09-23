package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalaryUsecase struct {
	repo *postgres.SalaryRepository
}

func NewSalaryUsecase(repo *postgres.SalaryRepository) *SalaryUsecase {
	return &SalaryUsecase{repo: repo}
}

func (u *SalaryUsecase) ListAll(ctx context.Context) ([]postgres.Salary, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalaryUsecase) GetByID(ctx context.Context, id int64) (*postgres.Salary, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalaryUsecase) Create(ctx context.Context, s *postgres.Salary) error {
	if s.EmployeeID == 0 {
		return errors.New("karyawan wajib dipilih")
	}
	if s.PayPeriod == "" {
		return errors.New("periode gaji wajib diisi")
	}
	if s.BaseSalary <= 0 {
		return errors.New("gaji pokok harus lebih dari 0")
	}
	id, err := u.repo.Create(ctx, s)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (u *SalaryUsecase) Update(ctx context.Context, s *postgres.Salary) error {
	return u.repo.Update(ctx, s)
}

func (u *SalaryUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *SalaryUsecase) MarkPaid(ctx context.Context, id int64) error {
	return u.repo.MarkPaid(ctx, id)
}
