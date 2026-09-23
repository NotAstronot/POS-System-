package usecase

import (
	"context"
	"errors"
	"pos-system/internal/domain"
	"pos-system/internal/repository/postgres"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo *postgres.UserRepository
}

func NewAuthUsecase(userRepo *postgres.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
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
