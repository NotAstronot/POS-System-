package usecase

import (
	"context"
	"pos-system/internal/repository/postgres"
)

type TransactionUsecase struct {
	transactionRepo *postgres.TransactionRepository
	orderRepo       *postgres.OrderRepository
}

func NewTransactionUsecase(transactionRepo *postgres.TransactionRepository, orderRepo *postgres.OrderRepository) *TransactionUsecase {
	return &TransactionUsecase{transactionRepo: transactionRepo, orderRepo: orderRepo}
}

func (u *TransactionUsecase) List(ctx context.Context) ([]postgres.Transaction, error) {
	return u.transactionRepo.List(ctx)
}

func (u *TransactionUsecase) GetByID(ctx context.Context, id int64) (*postgres.Transaction, error) {
	return u.transactionRepo.GetByID(ctx, id)
}

func (u *TransactionUsecase) Create(ctx context.Context, orderID int64, amount float64, paymentMethod string) (*postgres.Transaction, error) {
	return u.transactionRepo.Create(ctx, orderID, amount, paymentMethod)
}

func (u *TransactionUsecase) GetLogByID(ctx context.Context, id int64) (*postgres.TransactionLog, error) {
	return u.transactionRepo.GetLogByID(ctx, id)
}

func (u *TransactionUsecase) GetLogsByUsername(ctx context.Context, username string) ([]postgres.TransactionLog, error) {
	return u.transactionRepo.GetLogsByUsername(ctx, username)
}
