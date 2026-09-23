package usecase

import (
	"context"
	"pos-system/internal/repository/postgres"
)

type BranchUsecase struct {
	repo *postgres.BranchRepository
}

func NewBranchUsecase(repo *postgres.BranchRepository) *BranchUsecase {
	return &BranchUsecase{repo: repo}
}

func (u *BranchUsecase) List(ctx context.Context) ([]postgres.Branch, error) {
	return u.repo.ListAll(ctx)
}

func (u *BranchUsecase) GetByID(ctx context.Context, id int64) (*postgres.Branch, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *BranchUsecase) Create(ctx context.Context, b *postgres.Branch) (int64, error) {
	return u.repo.Create(ctx, b)
}

func (u *BranchUsecase) Update(ctx context.Context, b *postgres.Branch) error {
	return u.repo.Update(ctx, b)
}

func (u *BranchUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
