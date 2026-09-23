package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesInvoiceUsecase struct {
	repo *postgres.SalesInvoiceRepository
}

func NewSalesInvoiceUsecase(repo *postgres.SalesInvoiceRepository) *SalesInvoiceUsecase {
	return &SalesInvoiceUsecase{repo: repo}
}

func (u *SalesInvoiceUsecase) ListAll(ctx context.Context) ([]postgres.SalesInvoice, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesInvoiceUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesInvoice, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesInvoiceUsecase) Create(ctx context.Context, i *postgres.SalesInvoice) error {
	if i.SalesOrderID == 0 {
		return errors.New("sales order wajib dipilih")
	}
	if i.CustomerID == 0 {
		return errors.New("pelanggan wajib dipilih")
	}
	if i.TotalAmount <= 0 {
		return errors.New("total faktur harus lebih dari 0")
	}
	number, err := u.repo.GenerateInvoiceNumber(ctx)
	if err != nil {
		return err
	}
	i.InvoiceNumber = number
	i.Status = "unpaid"
	i.PaidAmount = 0

	id, err := u.repo.Create(ctx, i)
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}

func (u *SalesInvoiceUsecase) Update(ctx context.Context, i *postgres.SalesInvoice) error {
	if i.PaidAmount >= i.TotalAmount {
		i.Status = "paid"
	} else if i.PaidAmount > 0 {
		i.Status = "partial"
	} else {
		i.Status = "unpaid"
	}
	return u.repo.Update(ctx, i)
}

func (u *SalesInvoiceUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
