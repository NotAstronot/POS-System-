package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type StockTransferUsecase struct {
	repo      *postgres.StockTransferRepository
	stockRepo *postgres.StockRepository
}

func NewStockTransferUsecase(repo *postgres.StockTransferRepository, stockRepo *postgres.StockRepository) *StockTransferUsecase {
	return &StockTransferUsecase{repo: repo, stockRepo: stockRepo}
}

func (u *StockTransferUsecase) ListAll(ctx context.Context) ([]postgres.StockTransfer, error) {
	return u.repo.ListAll(ctx)
}

func (u *StockTransferUsecase) GetByID(ctx context.Context, id int64) (*postgres.StockTransfer, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *StockTransferUsecase) Create(ctx context.Context, t *postgres.StockTransfer) error {
	if t.FromWarehouseID == 0 || t.ToWarehouseID == 0 {
		return errors.New("gudang asal dan tujuan wajib dipilih")
	}
	if t.FromWarehouseID == t.ToWarehouseID {
		return errors.New("gudang asal dan tujuan tidak boleh sama")
	}
	if len(t.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	number, err := u.repo.GenerateTransferNumber(ctx)
	if err != nil {
		return err
	}
	t.TransferNumber = number
	t.Status = "draft"
	id, err := u.repo.Create(ctx, t)
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (u *StockTransferUsecase) Update(ctx context.Context, t *postgres.StockTransfer) error {
	if t.FromWarehouseID == t.ToWarehouseID {
		return errors.New("gudang asal dan tujuan tidak boleh sama")
	}
	return u.repo.Update(ctx, t)
}

func (u *StockTransferUsecase) Send(ctx context.Context, id int64) error {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("transfer tidak ditemukan")
	}
	if t.Status == "sent" || t.Status == "received" {
		return errors.New("transfer sudah dikirim")
	}
	if t.Status == "cancelled" {
		return errors.New("transfer dibatalkan")
	}
	for _, item := range t.Items {
		cur, err := u.stockRepo.GetStock(ctx, item.ProductID, t.FromWarehouseID)
		if err != nil {
			cur = 0
		}
		if cur < item.Quantity {
			return errors.New("stok tidak cukup di gudang asal: " + item.ProductName)
		}
	}
	for _, item := range t.Items {
		_ = u.stockRepo.AddStock(ctx, item.ProductID, t.FromWarehouseID, -item.Quantity, "transfer_out", "stock_transfers", &t.ID, "Transfer keluar "+t.TransferNumber, t.CreatedBy)
	}
	return u.repo.UpdateStatus(ctx, id, "sent")
}

func (u *StockTransferUsecase) Receive(ctx context.Context, id int64) error {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("transfer tidak ditemukan")
	}
	if t.Status == "received" {
		return errors.New("transfer sudah diterima")
	}
	if t.Status == "cancelled" {
		return errors.New("transfer dibatalkan")
	}
	if t.Status != "sent" {
		return errors.New("transfer harus dikirim terlebih dahulu")
	}
	for _, item := range t.Items {
		_ = u.stockRepo.AddStock(ctx, item.ProductID, t.ToWarehouseID, item.Quantity, "transfer_in", "stock_transfers", &t.ID, "Transfer masuk "+t.TransferNumber, t.CreatedBy)
	}
	return u.repo.UpdateStatus(ctx, id, "received")
}

func (u *StockTransferUsecase) Cancel(ctx context.Context, id int64) error {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("transfer tidak ditemukan")
	}
	if t.Status == "received" {
		return errors.New("transfer sudah diterima, tidak bisa dibatalkan")
	}
	if t.Status == "sent" {
		for _, item := range t.Items {
			_ = u.stockRepo.AddStock(ctx, item.ProductID, t.FromWarehouseID, item.Quantity, "transfer_cancel", "stock_transfers", &t.ID, "Pembatalan transfer "+t.TransferNumber, t.CreatedBy)
		}
	}
	return u.repo.UpdateStatus(ctx, id, "cancelled")
}

func (u *StockTransferUsecase) Delete(ctx context.Context, id int64) error {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("transfer tidak ditemukan")
	}
	if t.Status == "sent" || t.Status == "received" {
		return errors.New("transfer yang sudah dikirim tidak bisa dihapus")
	}
	return u.repo.Delete(ctx, id)
}

type StockAdjustmentUsecase struct {
	repo      *postgres.StockAdjustmentRepository
	stockRepo *postgres.StockRepository
}

func NewStockAdjustmentUsecase(repo *postgres.StockAdjustmentRepository, stockRepo *postgres.StockRepository) *StockAdjustmentUsecase {
	return &StockAdjustmentUsecase{repo: repo, stockRepo: stockRepo}
}

func (u *StockAdjustmentUsecase) ListAll(ctx context.Context) ([]postgres.StockAdjustment, error) {
	return u.repo.ListAll(ctx)
}

func (u *StockAdjustmentUsecase) GetByID(ctx context.Context, id int64) (*postgres.StockAdjustment, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *StockAdjustmentUsecase) Create(ctx context.Context, a *postgres.StockAdjustment) error {
	if a.WarehouseID == 0 {
		return errors.New("gudang wajib dipilih")
	}
	if len(a.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	number, err := u.repo.GenerateAdjustmentNumber(ctx)
	if err != nil {
		return err
	}
	a.AdjustmentNumber = number
	a.Status = "applied"
	id, err := u.repo.Create(ctx, a)
	if err != nil {
		return err
	}
	a.ID = id

	for _, item := range a.Items {
		diff := item.ActualQty - item.SystemQty
		if diff == 0 {
			continue
		}
		mtype := "adjustment"
		if diff > 0 {
			mtype = "adjustment_plus"
		} else {
			mtype = "adjustment_minus"
		}
		_ = u.stockRepo.AddStock(ctx, item.ProductID, a.WarehouseID, diff, mtype, "stock_adjustments", &a.ID, a.Notes, a.CreatedBy)
	}
	return nil
}

func (u *StockAdjustmentUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

type PriceChangeUsecase struct {
	repo *postgres.PriceChangeRepository
}

type StockUsecase struct {
	repo *postgres.StockRepository
}

func NewStockUsecase(repo *postgres.StockRepository) *StockUsecase {
	return &StockUsecase{repo: repo}
}

func (u *StockUsecase) ListMovements(ctx context.Context) ([]postgres.StockMovement, error) {
	return u.repo.ListMovements(ctx)
}

func NewPriceChangeUsecase(repo *postgres.PriceChangeRepository) *PriceChangeUsecase {
	return &PriceChangeUsecase{repo: repo}
}

func (u *PriceChangeUsecase) ListAll(ctx context.Context) ([]postgres.PriceChange, error) {
	return u.repo.ListAll(ctx)
}

func (u *PriceChangeUsecase) GetByID(ctx context.Context, id int64) (*postgres.PriceChange, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PriceChangeUsecase) Create(ctx context.Context, p *postgres.PriceChange) error {
	if len(p.Items) == 0 {
		return errors.New("minimal harus ada 1 item")
	}
	for _, item := range p.Items {
		if item.NewPrice < 0 {
			return errors.New("harga tidak boleh negatif: " + item.ProductName)
		}
	}
	number, err := u.repo.GeneratePriceChangeNumber(ctx)
	if err != nil {
		return err
	}
	p.PriceChangeNumber = number
	p.Status = "applied"
	id, err := u.repo.Create(ctx, p)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (u *PriceChangeUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
