package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesReceiptUsecase struct {
	repo        *postgres.SalesReceiptRepository
	invoiceRepo *postgres.SalesInvoiceRepository
}

func NewSalesReceiptUsecase(repo *postgres.SalesReceiptRepository, invoiceRepo *postgres.SalesInvoiceRepository) *SalesReceiptUsecase {
	return &SalesReceiptUsecase{repo: repo, invoiceRepo: invoiceRepo}
}

func (u *SalesReceiptUsecase) ListAll(ctx context.Context) ([]postgres.SalesReceipt, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesReceiptUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesReceipt, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesReceiptUsecase) Create(ctx context.Context, rc *postgres.SalesReceipt) error {
	if rc.SalesInvoiceID == 0 {
		return errors.New("faktur wajib dipilih")
	}
	if rc.Amount <= 0 {
		return errors.New("jumlah pembayaran harus lebih dari 0")
	}

	invoice, err := u.invoiceRepo.GetByID(ctx, rc.SalesInvoiceID)
	if err != nil {
		return errors.New("faktur tidak ditemukan")
	}
	if invoice.Status == "paid" {
		return errors.New("faktur sudah lunas")
	}
	if invoice.Remaining < rc.Amount {
		return errors.New("jumlah pembayaran melebihi sisa tagihan")
	}

	number, err := u.repo.GenerateReceiptNumber(ctx)
	if err != nil {
		return err
	}
	rc.ReceiptNumber = number

	id, err := u.repo.Create(ctx, rc)
	if err != nil {
		return err
	}
	rc.ID = id

	newPaid := invoice.PaidAmount + rc.Amount
	newStatus := "partial"
	if newPaid >= invoice.TotalAmount {
		newStatus = "paid"
	}
	invoice.PaidAmount = newPaid
	invoice.Status = newStatus
	_ = u.invoiceRepo.UpdatePaidAmount(ctx, invoice.ID, newPaid)

	return nil
}

func (u *SalesReceiptUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
