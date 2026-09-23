package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type SalesReturnUsecase struct {
	repo      *postgres.SalesReturnRepository
	orderRepo *postgres.SalesOrderRepository
	stockRepo *postgres.StockRepository
}

func NewSalesReturnUsecase(repo *postgres.SalesReturnRepository, orderRepo *postgres.SalesOrderRepository, stockRepo *postgres.StockRepository) *SalesReturnUsecase {
	return &SalesReturnUsecase{repo: repo, orderRepo: orderRepo, stockRepo: stockRepo}
}

func (u *SalesReturnUsecase) ListAll(ctx context.Context) ([]postgres.SalesReturn, error) {
	return u.repo.ListAll(ctx)
}

func (u *SalesReturnUsecase) GetByID(ctx context.Context, id int64) (*postgres.SalesReturn, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SalesReturnUsecase) Create(ctx context.Context, rt *postgres.SalesReturn) error {
	if rt.SalesOrderID == 0 {
		return errors.New("sales order wajib dipilih")
	}
	if len(rt.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	number, err := u.repo.GenerateReturnNumber(ctx)
	if err != nil {
		return err
	}
	rt.ReturnNumber = number
	rt.Status = "pending"

	total := float64(0)
	for i := range rt.Items {
		rt.Items[i].Subtotal = rt.Items[i].Quantity * rt.Items[i].UnitPrice
		total += rt.Items[i].Subtotal
	}
	rt.TotalAmount = total

	id, err := u.repo.Create(ctx, rt)
	if err != nil {
		return err
	}
	rt.ID = id

	if u.stockRepo != nil {
		for _, item := range rt.Items {
			productID, err := u.stockRepo.GetProductIDByName(ctx, item.ProductName)
			if err != nil {
				continue
			}
			whID, err := u.stockRepo.GetDefaultWarehouseID(ctx)
			if err != nil {
				continue
			}
			_ = u.stockRepo.AddStock(ctx, productID, whID, item.Quantity, "sales_return", "sales_returns", &rt.ID, "Retur penjualan "+rt.ReturnNumber, rt.CreatedBy)
		}
	}
	return nil
}

func (u *SalesReturnUsecase) Update(ctx context.Context, rt *postgres.SalesReturn) error {
	total := float64(0)
	for i := range rt.Items {
		rt.Items[i].Subtotal = rt.Items[i].Quantity * rt.Items[i].UnitPrice
		total += rt.Items[i].Subtotal
	}
	rt.TotalAmount = total
	return u.repo.Update(ctx, rt)
}

func (u *SalesReturnUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
