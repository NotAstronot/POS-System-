package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type PurchasePaymentUsecase struct {
	paymentRepo *postgres.PurchasePaymentRepository
	invoiceRepo *postgres.PurchaseInvoiceRepository
}

func NewPurchasePaymentUsecase(paymentRepo *postgres.PurchasePaymentRepository, invoiceRepo *postgres.PurchaseInvoiceRepository) *PurchasePaymentUsecase {
	return &PurchasePaymentUsecase{paymentRepo: paymentRepo, invoiceRepo: invoiceRepo}
}

func (u *PurchasePaymentUsecase) ListAll(ctx context.Context) ([]postgres.PurchasePayment, error) {
	return u.paymentRepo.ListAll(ctx)
}

func (u *PurchasePaymentUsecase) GetByID(ctx context.Context, id int64) (*postgres.PurchasePayment, error) {
	return u.paymentRepo.GetByID(ctx, id)
}

func (u *PurchasePaymentUsecase) Create(ctx context.Context, pp *postgres.PurchasePayment) error {
	if pp.PurchaseInvoiceID == 0 {
		return errors.New("faktur pembelian wajib dipilih")
	}
	if pp.Amount <= 0 {
		return errors.New("jumlah pembayaran harus lebih dari 0")
	}
	invoice, err := u.invoiceRepo.GetByID(ctx, pp.PurchaseInvoiceID)
	if err != nil {
		return errors.New("faktur tidak ditemukan")
	}
	if invoice.PaidAmount+pp.Amount > invoice.TotalAmount {
		return errors.New("jumlah pembayaran melebihi sisa tagihan")
	}

	paymentNumber, err := u.paymentRepo.GeneratePaymentNumber(ctx)
	if err != nil {
		return err
	}
	pp.PaymentNumber = paymentNumber

	id, err := u.paymentRepo.Create(ctx, pp)
	if err != nil {
		return err
	}
	pp.ID = id

	newPaidAmount := invoice.PaidAmount + pp.Amount
	status := "partial"
	if newPaidAmount >= invoice.TotalAmount {
		status = "paid"
	}
	_ = u.invoiceRepo.UpdatePaidAmount(ctx, pp.PurchaseInvoiceID, newPaidAmount, status)

	return nil
}

func (u *PurchasePaymentUsecase) Delete(ctx context.Context, id int64) error {
	pp, err := u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("pembayaran tidak ditemukan")
	}
	invoice, err := u.invoiceRepo.GetByID(ctx, pp.PurchaseInvoiceID)
	if err != nil {
		return errors.New("faktur tidak ditemukan")
	}
	newPaidAmount := invoice.PaidAmount - pp.Amount
	if newPaidAmount < 0 {
		newPaidAmount = 0
	}
	status := "unpaid"
	if newPaidAmount > 0 {
		status = "partial"
	}
	if newPaidAmount >= invoice.TotalAmount {
		status = "paid"
	}
	_ = u.invoiceRepo.UpdatePaidAmount(ctx, pp.PurchaseInvoiceID, newPaidAmount, status)
	return u.paymentRepo.Delete(ctx, id)
}
