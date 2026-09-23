package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *postgres.UserRepository
}

func NewUserUsecase(userRepo *postgres.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) List(ctx context.Context) ([]postgres.User, error) {
	return u.userRepo.ListAll(ctx)
}

func (u *UserUsecase) GetByID(ctx context.Context, id int64) (*postgres.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *UserUsecase) Create(ctx context.Context, user *postgres.User, password string) error {
	if user.Username == "" {
		return errors.New("username wajib diisi")
	}
	if password == "" {
		return errors.New("password wajib diisi")
	}
	pinRe := regexp.MustCompile(`^\d{6}$`)
	if !pinRe.MatchString(password) {
		return errors.New("password harus tepat 6 digit angka")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	if user.Role == "" {
		user.Role = "cashier"
	}
	if user.OutletID == "" {
		user.OutletID = "b0000000-0000-0000-0000-000000000001"
	}
	if user.TenantID == "" {
		user.TenantID = "a0000000-0000-0000-0000-000000000001"
	}
	user.IsActive = true
	_, err = u.userRepo.Create(ctx, user)
	return err
}

func (u *UserUsecase) Update(ctx context.Context, user *postgres.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *UserUsecase) UpdatePassword(ctx context.Context, id int64, password string) error {
	if password == "" {
		return errors.New("password wajib diisi")
	}
	pinRe := regexp.MustCompile(`^\d{6}$`)
	if !pinRe.MatchString(password) {
		return errors.New("password harus tepat 6 digit angka")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePassword(ctx, id, string(hash))
}

func (u *UserUsecase) Delete(ctx context.Context, id int64) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUsecase) GetPermissions(ctx context.Context, userID int64) ([]postgres.Permission, error) {
	return u.userRepo.GetPermissions(ctx, userID)
}

func (u *UserUsecase) SetPermissions(ctx context.Context, userID int64, permissionIDs []int64) error {
	return u.userRepo.SetPermissions(ctx, userID, permissionIDs)
}

func (u *UserUsecase) GetAllPermissions(ctx context.Context) ([]postgres.Permission, error) {
	return u.userRepo.GetAllPermissions(ctx)
}
