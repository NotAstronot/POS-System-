package usecase

import (
	"context"
	"errors"
	"pos-system/internal/domain"
	"pos-system/internal/repository/postgres"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo   *postgres.UserRepository
	tenantRepo *postgres.TenantRepository
}

func NewAuthUsecase(userRepo *postgres.UserRepository, tenantRepo *postgres.TenantRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, tenantRepo: tenantRepo}
}

func (u *AuthUsecase) Login(ctx context.Context, username, password string) (*postgres.User, error) {
	user, err := u.userRepo.FindByUsernameGlobal(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.New("akun telah dinonaktifkan")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (u *AuthUsecase) GetPermissions(ctx context.Context, tenantID, userID int64) ([]postgres.Permission, error) {
	return u.userRepo.GetPermissions(domain.WithTenantID(ctx, tenantID), userID)
}

const defaultOutletID = "b0000000-0000-0000-0000-000000000001"

type RegisterMerchantInput struct {
	MerchantName string
	Slug         string
	Username     string
	OwnerName    string
	Password     string
}

// RegisterMerchant creates a new tenant and its owner user in one shot
// (public endpoint — no tenant context yet, everything is system-scoped).
func (u *AuthUsecase) RegisterMerchant(ctx context.Context, in RegisterMerchantInput) (*postgres.Tenant, *postgres.User, error) {
	in.MerchantName = strings.TrimSpace(in.MerchantName)
	in.Slug = strings.TrimSpace(in.Slug)
	in.Username = strings.TrimSpace(in.Username)
	in.OwnerName = strings.TrimSpace(in.OwnerName)

	if in.MerchantName == "" {
		return nil, nil, errors.New("nama merchant wajib diisi")
	}
	if in.Username == "" {
		return nil, nil, errors.New("username wajib diisi")
	}
	if in.OwnerName == "" {
		in.OwnerName = in.Username
	}
	if !regexp.MustCompile(`^\d{6}$`).MatchString(in.Password) {
		return nil, nil, errors.New("password harus tepat 6 digit angka")
	}

	if existing, _ := u.userRepo.FindByUsernameGlobal(ctx, in.Username); existing != nil {
		return nil, nil, errors.New("username sudah digunakan")
	}

	slug := in.Slug
	explicitSlug := slug != ""
	if !explicitSlug {
		slug = slugify(in.MerchantName)
	}
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(slug) {
		return nil, nil, errors.New("slug hanya boleh huruf kecil, angka, dan strip")
	}
	if explicitSlug {
		if t, _ := u.tenantRepo.GetBySlug(ctx, slug); t != nil {
			return nil, nil, errors.New("slug sudah digunakan")
		}
	} else {
		base := slug
		for i := 2; ; i++ {
			t, _ := u.tenantRepo.GetBySlug(ctx, slug)
			if t == nil {
				break
			}
			if i > 100 {
				return nil, nil, errors.New("slug sudah digunakan, tentukan slug secara manual")
			}
			slug = base + "-" + strconv.Itoa(i)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	tenant := &postgres.Tenant{
		Name:             in.MerchantName,
		Slug:             slug,
		SubscriptionPlan: "free",
		MaxUsers:         5,
		MaxProducts:      500,
		MaxBranches:      1,
		Settings:         "{}",
		Status:           "active",
	}
	owner := &postgres.User{
		Username: in.Username,
		Name:     in.OwnerName,
		Password: string(hash),
		Role:     "admin",
		OutletID: defaultOutletID,
	}
	if err := u.tenantRepo.CreateWithOwner(ctx, tenant, owner); err != nil {
		return nil, nil, err
	}
	return tenant, owner, nil
}
