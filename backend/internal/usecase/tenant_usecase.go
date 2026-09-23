package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"regexp"
	"strings"
)

type TenantUsecase struct {
	repo *postgres.TenantRepository
}

func NewTenantUsecase(repo *postgres.TenantRepository) *TenantUsecase {
	return &TenantUsecase{repo: repo}
}

func (u *TenantUsecase) List(ctx context.Context) ([]postgres.Tenant, error) {
	return u.repo.List(ctx)
}

func (u *TenantUsecase) GetByID(ctx context.Context, id int64) (*postgres.Tenant, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *TenantUsecase) Create(ctx context.Context, t *postgres.Tenant) (*postgres.Tenant, error) {
	if t.Name == "" {
		return nil, errors.New("nama tenant wajib diisi")
	}
	if t.Slug == "" {
		t.Slug = slugify(t.Name)
	}

	re := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !re.MatchString(t.Slug) {
		return nil, errors.New("slug hanya boleh huruf kecil, angka, dan strip")
	}

	existing, _ := u.repo.GetBySlug(ctx, t.Slug)
	if existing != nil {
		return nil, errors.New("slug sudah digunakan")
	}

	if t.SubscriptionPlan == "" {
		t.SubscriptionPlan = "free"
	}
	if t.MaxUsers == 0 {
		t.MaxUsers = 5
	}
	if t.MaxProducts == 0 {
		t.MaxProducts = 500
	}
	if t.MaxBranches == 0 {
		t.MaxBranches = 1
	}
	if t.Status == "" {
		t.Status = "active"
	}
	if t.Settings == "" {
		t.Settings = "{}"
	}

	return u.repo.Create(ctx, t)
}

func (u *TenantUsecase) Update(ctx context.Context, t *postgres.Tenant) (*postgres.Tenant, error) {
	existing, err := u.repo.GetByID(ctx, t.ID)
	if err != nil {
		return nil, errors.New("tenant tidak ditemukan")
	}

	if t.Name != "" {
		existing.Name = t.Name
	}
	if t.Slug != "" {
		re := regexp.MustCompile(`^[a-z0-9-]+$`)
		if !re.MatchString(t.Slug) {
			return nil, errors.New("slug hanya boleh huruf kecil, angka, dan strip")
		}
		existing.Slug = t.Slug
	}
	if t.Domain != "" {
		existing.Domain = t.Domain
	}
	if t.LogoURL != "" {
		existing.LogoURL = t.LogoURL
	}
	if t.SubscriptionPlan != "" {
		existing.SubscriptionPlan = t.SubscriptionPlan
	}
	if t.SubscriptionExpiresAt != nil {
		existing.SubscriptionExpiresAt = t.SubscriptionExpiresAt
	}
	if t.MaxUsers > 0 {
		existing.MaxUsers = t.MaxUsers
	}
	if t.MaxProducts > 0 {
		existing.MaxProducts = t.MaxProducts
	}
	if t.MaxBranches > 0 {
		existing.MaxBranches = t.MaxBranches
	}
	if t.Settings != "" {
		existing.Settings = t.Settings
	}
	if t.Status != "" {
		existing.Status = t.Status
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *TenantUsecase) Delete(ctx context.Context, id int64) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("tenant tidak ditemukan")
	}
	return u.repo.Delete(ctx, id)
}

func (u *TenantUsecase) GetUsage(ctx context.Context, tenantID int64) (map[string]int, error) {
	users, err := u.repo.CountUsers(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	products, err := u.repo.CountProducts(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	branches, err := u.repo.CountBranches(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return map[string]int{
		"users":    users,
		"products": products,
		"branches": branches,
	}, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`[^a-z0-9-]`)
	s = re.ReplaceAllString(s, "-")
	re2 := regexp.MustCompile(`-+`)
	s = re2.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "tenant"
	}
	return s
}
