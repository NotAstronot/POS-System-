package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type DeliveryOrderUsecase struct {
	repo      *postgres.DeliveryOrderRepository
	orderRepo *postgres.SalesOrderRepository
	stockRepo *postgres.StockRepository
}

func NewDeliveryOrderUsecase(repo *postgres.DeliveryOrderRepository, orderRepo *postgres.SalesOrderRepository, stockRepo *postgres.StockRepository) *DeliveryOrderUsecase {
	return &DeliveryOrderUsecase{repo: repo, orderRepo: orderRepo, stockRepo: stockRepo}
}

func (u *DeliveryOrderUsecase) ListAll(ctx context.Context) ([]postgres.DeliveryOrder, error) {
	return u.repo.ListAll(ctx)
}

func (u *DeliveryOrderUsecase) GetByID(ctx context.Context, id int64) (*postgres.DeliveryOrder, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *DeliveryOrderUsecase) Create(ctx context.Context, d *postgres.DeliveryOrder) error {
	if d.SalesOrderID == 0 {
		return errors.New("sales order wajib dipilih")
	}
	if len(d.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}

	order, err := u.orderRepo.GetByID(ctx, d.SalesOrderID)
	if err != nil {
		return errors.New("sales order tidak ditemukan")
	}
	if order.Status == "cancelled" {
		return errors.New("sales order dibatalkan, tidak bisa dikirim")
	}

	number, err := u.repo.GenerateDONumber(ctx)
	if err != nil {
		return err
	}
	d.DONumber = number
	d.Status = "shipped"

	id, err := u.repo.Create(ctx, d)
	if err != nil {
		return err
	}
	d.ID = id

	for _, item := range d.Items {
		for _, existing := range order.Items {
			if existing.ProductName == item.ProductName {
				newQty := existing.DeliveredQty + item.Quantity
				if newQty > existing.Quantity {
					return errors.New("jumlah kirim melebihi jumlah pesanan: " + item.ProductName)
				}
				_ = u.orderRepo.UpdateItemDeliveredQty(ctx, existing.ID, newQty)
			}
		}
		if u.stockRepo != nil {
			productID, err := u.stockRepo.GetProductIDByName(ctx, item.ProductName)
			if err == nil {
				whID, _ := u.stockRepo.GetDefaultWarehouseID(ctx)
				if whID != 0 {
					cur, _ := u.stockRepo.GetStock(ctx, productID, whID)
					if cur < item.Quantity {
						return errors.New("stok tidak cukup di gudang utama: " + item.ProductName)
					}
					_ = u.stockRepo.AddStock(ctx, productID, whID, -item.Quantity, "delivery_out", "delivery_orders", &d.ID, "Pengiriman "+d.DONumber, d.CreatedBy)
				}
			}
		}
	}

	allDelivered := true
	updatedOrder, _ := u.orderRepo.GetByID(ctx, d.SalesOrderID)
	for _, item := range updatedOrder.Items {
		if item.DeliveredQty < item.Quantity {
			allDelivered = false
			break
		}
	}
	if allDelivered {
		updatedOrder.Status = "delivered"
	} else {
		updatedOrder.Status = "partial"
	}
	_ = u.orderRepo.Update(ctx, updatedOrder)

	return nil
}

func (u *DeliveryOrderUsecase) Update(ctx context.Context, d *postgres.DeliveryOrder) error {
	return u.repo.Update(ctx, d)
}

func (u *DeliveryOrderUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
