package usecase

import (
	"pos-system/internal/repository/postgres"
)

type ShiftUsecase struct {
	shiftRepo    *postgres.ShiftRepository
	orderRepo    *postgres.OrderRepository
	orderItemRepo *postgres.OrderItemRepository
	productRepo  *postgres.ProductRepository
}

func NewShiftUsecase(shiftRepo *postgres.ShiftRepository, orderRepo *postgres.OrderRepository, orderItemRepo *postgres.OrderItemRepository, productRepo *postgres.ProductRepository) *ShiftUsecase {
	return &ShiftUsecase{shiftRepo: shiftRepo, orderRepo: orderRepo, orderItemRepo: orderItemRepo, productRepo: productRepo}
}

func (u *ShiftUsecase) GetActiveShift() (*postgres.Shift, error) {
	return u.shiftRepo.GetActiveShift()
}
