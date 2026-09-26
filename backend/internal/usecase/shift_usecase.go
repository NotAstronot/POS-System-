package usecase

import (
	"context"
	"pos-system/internal/repository/postgres"
)

type ShiftUsecase struct {
	shiftRepo     *postgres.ShiftRepository
	orderRepo     *postgres.OrderRepository
	orderItemRepo *postgres.OrderItemRepository
	productRepo   *postgres.ProductRepository
}

func NewShiftUsecase(shiftRepo *postgres.ShiftRepository, orderRepo *postgres.OrderRepository, orderItemRepo *postgres.OrderItemRepository, productRepo *postgres.ProductRepository) *ShiftUsecase {
	return &ShiftUsecase{shiftRepo: shiftRepo, orderRepo: orderRepo, orderItemRepo: orderItemRepo, productRepo: productRepo}
}

func (u *ShiftUsecase) GetActiveShift(ctx context.Context) (*postgres.Shift, error) {
	return u.shiftRepo.GetActiveShift(ctx)
}

func (u *ShiftUsecase) History(ctx context.Context, limit int) ([]postgres.ShiftHistory, error) {
	return u.shiftRepo.History(ctx, limit)
}
