package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type PurchaseReturnUsecase struct {
	repo      *postgres.PurchaseReturnRepository
	stockRepo *postgres.StockRepository
}

func NewPurchaseReturnUsecase(repo *postgres.PurchaseReturnRepository, stockRepo *postgres.StockRepository) *PurchaseReturnUsecase {
	return &PurchaseReturnUsecase{repo: repo, stockRepo: stockRepo}
}

func (u *PurchaseReturnUsecase) ListAll(ctx context.Context) ([]postgres.PurchaseReturn, error) {
	return u.repo.ListAll(ctx)
}

func (u *PurchaseReturnUsecase) GetByID(ctx context.Context, id int64) (*postgres.PurchaseReturn, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PurchaseReturnUsecase) Create(ctx context.Context, pr *postgres.PurchaseReturn) error {
	if pr.PurchaseOrderID == 0 {
		return errors.New("pesanan pembelian wajib dipilih")
	}
	if pr.SupplierID == 0 {
		return errors.New("supplier wajib dipilih")
	}
	if len(pr.Items) == 0 {
		return errors.New("minimal harus ada 1 item retur")
	}
	returnNumber, err := u.repo.GenerateReturnNumber(ctx)
	if err != nil {
		return err
	}
	pr.ReturnNumber = returnNumber
	pr.Status = "pending"

	total := float64(0)
	for i := range pr.Items {
		pr.Items[i].Subtotal = pr.Items[i].Quantity * pr.Items[i].UnitPrice
		total += pr.Items[i].Subtotal
	}

	id, err := u.repo.Create(ctx, pr)
	if err != nil {
		return err
	}
	pr.ID = id

	if u.stockRepo != nil {
		for _, item := range pr.Items {
			productID, err := u.stockRepo.GetProductIDByName(ctx, item.ProductName)
			if err != nil {
				continue
			}
			whID, err := u.stockRepo.GetDefaultWarehouseID(ctx)
			if err != nil {
				continue
			}
			cur, _ := u.stockRepo.GetStock(ctx, productID, whID)
			if cur < item.Quantity {
				return errors.New("stok tidak cukup untuk retur: " + item.ProductName)
			}
			_ = u.stockRepo.AddStock(ctx, productID, whID, -item.Quantity, "purchase_return", "purchase_returns", &pr.ID, "Retur pembelian "+pr.ReturnNumber, pr.CreatedBy)
		}
	}
	return nil
}

func (u *PurchaseReturnUsecase) Update(ctx context.Context, pr *postgres.PurchaseReturn) error {
	return u.repo.Update(ctx, pr)
}

func (u *PurchaseReturnUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
