package usecase

import (
	"pos-system/internal/repository/postgres"
)

type TransactionUsecase struct {
	transactionRepo *postgres.TransactionRepository
	orderRepo       *postgres.OrderRepository
}

func NewTransactionUsecase(transactionRepo *postgres.TransactionRepository, orderRepo *postgres.OrderRepository) *TransactionUsecase {
	return &TransactionUsecase{transactionRepo: transactionRepo, orderRepo: orderRepo}
}

func (u *TransactionUsecase) List() ([]postgres.Transaction, error) {
	return u.transactionRepo.List()
}

func (u *TransactionUsecase) GetByID(id int64) (*postgres.Transaction, error) {
	return u.transactionRepo.GetByID(id)
}

func (u *TransactionUsecase) Create(orderID int64, amount float64, paymentMethod string) (*postgres.Transaction, error) {
	return u.transactionRepo.Create(orderID, amount, paymentMethod)
}
