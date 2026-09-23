package usecase

import (
	"context"
	"pos-system/internal/domain"
	"pos-system/internal/repository/postgres"
)

type CategoryUsecase struct {
	repo *postgres.CategoryRepository
}

func NewCategoryUsecase(repo *postgres.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (u *CategoryUsecase) List(ctx context.Context) ([]postgres.Category, error) {
	return u.repo.List(ctx)
}

func (u *CategoryUsecase) GetByID(ctx context.Context, id int64) (*postgres.Category, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CategoryUsecase) Create(ctx context.Context, name string) (*postgres.Category, error) {
	return u.repo.Create(ctx, name)
}

func (u *CategoryUsecase) Update(ctx context.Context, id int64, name string) (*postgres.Category, error) {
	return u.repo.Update(ctx, id, name)
}

func (u *CategoryUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func GetTenantID(c interface{ Value(any) any }) int64 {
	if ctx, ok := c.(interface {
		Value(key interface{}) interface{}
	}); ok {
		if v := ctx.Value(domain.TenantIDKey); v != nil {
			if id, ok := v.(int64); ok {
				return id
			}
		}
	}
	return 0
}
