package usecase

import (
	"errors"
	"pos-system/internal/repository/postgres"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo *postgres.UserRepository
}

func NewAuthUsecase(userRepo *postgres.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (u *AuthUsecase) Login(username, password string) (*postgres.User, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
