package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesOrderUsecase struct {
	repo *postgres.SalesOrderRepository
}

func NewSalesOrderUsecase(repo *postgres.SalesOrderRepository) *SalesOrderUsecase {
	return &SalesOrderUsecase{repo: repo}
}

func (u *SalesOrderUsecase) ListAll(ctx context.Context) ([]postgres.SalesOrder, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesOrderUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesOrder, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesOrderUsecase) Create(ctx context.Context, o *postgres.SalesOrder) error {
	if o.CustomerID == 0 {
		return errors.New("pelanggan wajib dipilih")
	}
	if len(o.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	number, err := u.repo.GenerateOrderNumber(ctx)
	if err != nil {
		return err
	}
	o.OrderNumber = number
	o.Status = "draft"

	total := float64(0)
	for i := range o.Items {
		o.Items[i].Subtotal = o.Items[i].Quantity * o.Items[i].UnitPrice
		total += o.Items[i].Subtotal
	}
	o.TotalAmount = total

	id, err := u.repo.Create(ctx, o)
	if err != nil {
		return err
	}
	o.ID = id
	return nil
}

func (u *SalesOrderUsecase) Update(ctx context.Context, o *postgres.SalesOrder) error {
	total := float64(0)
	for i := range o.Items {
		o.Items[i].Subtotal = o.Items[i].Quantity * o.Items[i].UnitPrice
		total += o.Items[i].Subtotal
	}
	o.TotalAmount = total
	return u.repo.Update(ctx, o)
}

func (u *SalesOrderUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
