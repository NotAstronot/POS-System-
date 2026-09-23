package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type PurchaseOrderUsecase struct {
	repo      *postgres.PurchaseOrderRepository
	stockRepo *postgres.StockRepository
}

func NewPurchaseOrderUsecase(repo *postgres.PurchaseOrderRepository, stockRepo *postgres.StockRepository) *PurchaseOrderUsecase {
	return &PurchaseOrderUsecase{repo: repo, stockRepo: stockRepo}
}

func (u *PurchaseOrderUsecase) ListAll(ctx context.Context) ([]postgres.PurchaseOrder, error) {
	return u.repo.ListAll(ctx)
}

func (u *PurchaseOrderUsecase) GetByID(ctx context.Context, id int64) (*postgres.PurchaseOrder, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PurchaseOrderUsecase) Create(ctx context.Context, po *postgres.PurchaseOrder) error {
	if po.SupplierID == 0 {
		return errors.New("supplier wajib dipilih")
	}
	if len(po.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	poNumber, err := u.repo.GeneratePONumber(ctx)
	if err != nil {
		return err
	}
	po.PONumber = poNumber
	po.Status = "draft"

	total := float64(0)
	for i := range po.Items {
		po.Items[i].Subtotal = po.Items[i].Quantity * po.Items[i].UnitPrice
		total += po.Items[i].Subtotal
	}
	po.TotalAmount = total

	id, err := u.repo.Create(ctx, po)
	if err != nil {
		return err
	}
	po.ID = id
	return nil
}

func (u *PurchaseOrderUsecase) Update(ctx context.Context, po *postgres.PurchaseOrder) error {
	total := float64(0)
	for i := range po.Items {
		po.Items[i].Subtotal = po.Items[i].Quantity * po.Items[i].UnitPrice
		total += po.Items[i].Subtotal
	}
	po.TotalAmount = total
	return u.repo.Update(ctx, po)
}

func (u *PurchaseOrderUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *PurchaseOrderUsecase) ReceiveItems(ctx context.Context, orderID, warehouseID int64, items []postgres.PurchaseItem) error {
	po, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}
	if po.Status == "received" || po.Status == "cancelled" {
		return errors.New("pesanan tidak dapat diterima")
	}
	for _, item := range items {
		for _, existing := range po.Items {
			if existing.ID == item.ID {
				newQty := existing.ReceivedQty + item.Quantity
				if newQty > existing.Quantity {
					return errors.New("jumlah melebihi jumlah yang dipesan")
				}
				_ = u.repo.UpdateItemReceivedQty(ctx, item.ID, newQty)

				if u.stockRepo != nil {
					productID, err := u.stockRepo.GetProductIDByName(ctx, existing.ProductName)
					if err != nil {
						continue
					}
					targetWh := warehouseID
					if targetWh == 0 {
						targetWh, _ = u.stockRepo.GetDefaultWarehouseID(ctx)
					}
					if targetWh != 0 {
						_ = u.stockRepo.AddStock(ctx, productID, targetWh, item.Quantity, "purchase_receive", "purchase_orders", &orderID, "Penerimaan PO "+po.PONumber, po.CreatedBy)
					}
				}
			}
		}
	}
	allReceived := true
	updatedPo, _ := u.repo.GetByID(ctx, orderID)
	for _, item := range updatedPo.Items {
		if item.ReceivedQty < item.Quantity {
			allReceived = false
			break
		}
	}
	if allReceived {
		updatedPo.Status = "received"
		_ = u.repo.Update(ctx, updatedPo)
	} else {
		updatedPo.Status = "partial"
		_ = u.repo.Update(ctx, updatedPo)
	}
	return nil
}
