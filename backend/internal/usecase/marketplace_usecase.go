package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"time"
)

type MarketplaceUsecase struct {
	repo           *postgres.MarketplaceRepository
	salesOrderRepo *postgres.SalesOrderRepository
}

func NewMarketplaceUsecase(repo *postgres.MarketplaceRepository, salesOrderRepo *postgres.SalesOrderRepository) *MarketplaceUsecase {
	return &MarketplaceUsecase{repo: repo, salesOrderRepo: salesOrderRepo}
}

func (u *MarketplaceUsecase) ListConnections(ctx context.Context) ([]postgres.MarketplaceConnection, error) {
	return u.repo.ListConnections(ctx)
}

func (u *MarketplaceUsecase) GetConnection(ctx context.Context, id int64) (*postgres.MarketplaceConnection, error) {
	return u.repo.GetConnection(ctx, id)
}

func (u *MarketplaceUsecase) CreateConnection(ctx context.Context, c *postgres.MarketplaceConnection) error {
	if c.Platform == "" {
		return errors.New("platform marketplace wajib diisi")
	}
	if c.ShopName == "" {
		return errors.New("nama toko wajib diisi")
	}
	if c.ShippingFeeAccount == "" {
		c.ShippingFeeAccount = "Biaya Ongkir"
	}
	if c.CommissionAccount == "" {
		c.CommissionAccount = "Komisi Platform"
	}
	_, err := u.repo.CreateConnection(ctx, c)
	return err
}

func (u *MarketplaceUsecase) UpdateConnection(ctx context.Context, c *postgres.MarketplaceConnection) error {
	if c.Platform == "" {
		return errors.New("platform marketplace wajib diisi")
	}
	if c.ShopName == "" {
		return errors.New("nama toko wajib diisi")
	}
	return u.repo.UpdateConnection(ctx, c)
}

func (u *MarketplaceUsecase) DeleteConnection(ctx context.Context, id int64) error {
	return u.repo.DeleteConnection(ctx, id)
}

func (u *MarketplaceUsecase) ListOrders(ctx context.Context) ([]postgres.MarketplaceOrder, error) {
	return u.repo.ListOrders(ctx)
}

func (u *MarketplaceUsecase) ImportOrders(ctx context.Context, connectionID int64, orders []postgres.MarketplaceOrder, createdBy *int64) (int, error) {
	if connectionID == 0 {
		return 0, errors.New("koneksi marketplace wajib dipilih")
	}
	conn, err := u.repo.GetConnection(ctx, connectionID)
	if err != nil {
		return 0, errors.New("koneksi marketplace tidak ditemukan")
	}
	if !conn.IsActive {
		return 0, errors.New("koneksi marketplace nonaktif")
	}
	if len(orders) == 0 {
		return 0, errors.New("tidak ada pesanan untuk diimpor")
	}

	imported := 0
	for _, o := range orders {
		if o.MarketplaceOrderID == "" {
			continue
		}
		exists, err := u.repo.OrderExists(ctx, connectionID, o.MarketplaceOrderID)
		if err != nil {
			return imported, err
		}
		if exists {
			continue
		}
		o.ConnectionID = connectionID
		if o.OrderDate == "" {
			o.OrderDate = time.Now().Format("2006-01-02")
		}
		o.Subtotal = o.Quantity * o.UnitPrice
		o.GrandTotal = o.Subtotal + o.ShippingFee + o.PlatformFee
		o.Status = "imported"
		o.CreatedBy = createdBy
		if err := u.repo.CreateOrder(ctx, &o); err != nil {
			return imported, err
		}
		imported++
	}

	if imported > 0 {
		_ = u.repo.UpdateLastSync(ctx, connectionID)
	}
	return imported, nil
}

func (u *MarketplaceUsecase) PostToSalesOrder(ctx context.Context, orderID int64, createdBy *int64) error {
	o, err := u.repo.GetOrder(ctx, orderID)
	if err != nil {
		return errors.New("pesanan marketplace tidak ditemukan")
	}
	if o.Status != "imported" {
		return errors.New("pesanan sudah diposting")
	}
	conn, err := u.repo.GetConnection(ctx, o.ConnectionID)
	if err != nil {
		return err
	}
	if conn.CustomerID == nil {
		return errors.New("koneksi marketplace belum dipetakan ke pelanggan (Customer)")
	}

	so := &postgres.SalesOrder{
		CustomerID:      *conn.CustomerID,
		SalesCategoryID: conn.SalesCategoryID,
		OrderDate:       o.OrderDate,
		Status:          "draft",
		Notes:           "Impor otomatis dari " + o.Platform + " (" + o.ShopName + ") - Order #" + o.MarketplaceOrderID,
		CreatedBy:       createdBy,
	}
	so.Items = append(so.Items, postgres.SalesOrderItem{
		ProductName: o.ProductName,
		Quantity:    o.Quantity,
		Unit:        "pcs",
		UnitPrice:   o.UnitPrice,
	})
	if o.ShippingFee > 0 {
		so.Items = append(so.Items, postgres.SalesOrderItem{
			ProductName: conn.ShippingFeeAccount + " (Ongkir)",
			Quantity:    1,
			Unit:        "ls",
			UnitPrice:   o.ShippingFee,
		})
	}
	if o.PlatformFee > 0 {
		so.Items = append(so.Items, postgres.SalesOrderItem{
			ProductName: conn.CommissionAccount + " (Komisi)",
			Quantity:    1,
			Unit:        "ls",
			UnitPrice:   o.PlatformFee,
		})
	}

	number, err := u.salesOrderRepo.GenerateOrderNumber(ctx)
	if err != nil {
		return err
	}
	so.OrderNumber = number
	total := float64(0)
	for i := range so.Items {
		so.Items[i].Subtotal = so.Items[i].Quantity * so.Items[i].UnitPrice
		total += so.Items[i].Subtotal
	}
	so.TotalAmount = total

	soID, err := u.salesOrderRepo.Create(ctx, so)
	if err != nil {
		return err
	}
	return u.repo.LinkSalesOrder(ctx, orderID, soID)
}

func (u *MarketplaceUsecase) DeleteOrder(ctx context.Context, id int64) error {
	return u.repo.DeleteOrder(ctx, id)
}
