package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type PurchaseInvoiceUsecase struct {
	repo *postgres.PurchaseInvoiceRepository
}

func NewPurchaseInvoiceUsecase(repo *postgres.PurchaseInvoiceRepository) *PurchaseInvoiceUsecase {
	return &PurchaseInvoiceUsecase{repo: repo}
}

func (u *PurchaseInvoiceUsecase) ListAll(ctx context.Context) ([]postgres.PurchaseInvoice, error) {
	return u.repo.ListAll(ctx)
}

func (u *PurchaseInvoiceUsecase) GetByID(ctx context.Context, id int64) (*postgres.PurchaseInvoice, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PurchaseInvoiceUsecase) Create(ctx context.Context, pi *postgres.PurchaseInvoice) error {
	if pi.PurchaseOrderID == 0 {
		return errors.New("pesanan pembelian wajib dipilih")
	}
	if pi.SupplierID == 0 {
		return errors.New("supplier wajib dipilih")
	}
	invoiceNumber, err := u.repo.GenerateInvoiceNumber(ctx)
	if err != nil {
		return err
	}
	pi.InvoiceNumber = invoiceNumber
	pi.Status = "unpaid"

	id, err := u.repo.Create(ctx, pi)
	if err != nil {
		return err
	}
	pi.ID = id
	return nil
}

func (u *PurchaseInvoiceUsecase) Update(ctx context.Context, pi *postgres.PurchaseInvoice) error {
	return u.repo.Update(ctx, pi)
}

func (u *PurchaseInvoiceUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
