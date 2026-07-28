package usecase

import "pos-system/internal/repository/postgres"

type CategoryUsecase struct {
	repo *postgres.CategoryRepository
}

func NewCategoryUsecase(repo *postgres.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (u *CategoryUsecase) List() ([]postgres.Category, error) {
	return u.repo.List()
}

func (u *CategoryUsecase) GetByID(id int64) (*postgres.Category, error) {
	return u.repo.GetByID(id)
}

func (u *CategoryUsecase) Create(name string) (*postgres.Category, error) {
	return u.repo.Create(name)
}

func (u *CategoryUsecase) Update(id int64, name string) (*postgres.Category, error) {
	return u.repo.Update(id, name)
}

func (u *CategoryUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}
