package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesQuotationUsecase struct {
	repo *postgres.SalesQuotationRepository
}

func NewSalesQuotationUsecase(repo *postgres.SalesQuotationRepository) *SalesQuotationUsecase {
	return &SalesQuotationUsecase{repo: repo}
}

func (u *SalesQuotationUsecase) ListAll(ctx context.Context) ([]postgres.SalesQuotation, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesQuotationUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesQuotation, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesQuotationUsecase) Create(ctx context.Context, q *postgres.SalesQuotation) error {
	if q.CustomerID == 0 {
		return errors.New("pelanggan wajib dipilih")
	}
	if len(q.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	number, err := u.repo.GenerateQuotationNumber(ctx)
	if err != nil {
		return err
	}
	q.QuotationNumber = number
	q.Status = "draft"

	total := float64(0)
	for i := range q.Items {
		q.Items[i].Subtotal = q.Items[i].Quantity * q.Items[i].UnitPrice
		total += q.Items[i].Subtotal
	}
	q.TotalAmount = total

	id, err := u.repo.Create(ctx, q)
	if err != nil {
		return err
	}
	q.ID = id
	return nil
}

func (u *SalesQuotationUsecase) Update(ctx context.Context, q *postgres.SalesQuotation) error {
	total := float64(0)
	for i := range q.Items {
		q.Items[i].Subtotal = q.Items[i].Quantity * q.Items[i].UnitPrice
		total += q.Items[i].Subtotal
	}
	q.TotalAmount = total
	return u.repo.Update(ctx, q)
}

func (u *SalesQuotationUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
