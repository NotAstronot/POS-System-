package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"strings"
)

type WarehouseUsecase struct {
	repo *postgres.WarehouseRepository
}

func NewWarehouseUsecase(repo *postgres.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

func (u *WarehouseUsecase) List(ctx context.Context) ([]postgres.Warehouse, error) {
	return u.repo.List(ctx)
}

func (u *WarehouseUsecase) GetByID(ctx context.Context, id int64) (*postgres.Warehouse, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *WarehouseUsecase) Create(ctx context.Context, w *postgres.Warehouse) error {
	w.Code = strings.TrimSpace(w.Code)
	w.Name = strings.TrimSpace(w.Name)
	if w.Code == "" {
		return errors.New("kode gudang wajib diisi")
	}
	if w.Name == "" {
		return errors.New("nama gudang wajib diisi")
	}
	return u.repo.Create(ctx, w)
}

func (u *WarehouseUsecase) Update(ctx context.Context, w *postgres.Warehouse) error {
	if w.Name == "" {
		return errors.New("nama gudang wajib diisi")
	}
	return u.repo.Update(ctx, w)
}

func (u *WarehouseUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

type ItemUsecase struct {
	itemRepo  *postgres.ItemRepository
	prodRepo  *postgres.ProductRepository
	stockRepo *postgres.StockRepository
}

func NewItemUsecase(itemRepo *postgres.ItemRepository, prodRepo *postgres.ProductRepository, stockRepo *postgres.StockRepository) *ItemUsecase {
	return &ItemUsecase{itemRepo: itemRepo, prodRepo: prodRepo, stockRepo: stockRepo}
}

func (u *ItemUsecase) List(ctx context.Context) ([]postgres.Item, error) {
	return u.itemRepo.List(ctx)
}

func (u *ItemUsecase) GetByID(ctx context.Context, id int64) (*postgres.Item, error) {
	return u.itemRepo.GetByID(ctx, id)
}

func (u *ItemUsecase) ListWarehouseStock(ctx context.Context, warehouseID int64) ([]postgres.ItemStock, error) {
	return u.itemRepo.ListStockByWarehouseID(ctx, warehouseID)
}

func (u *ItemUsecase) Create(ctx context.Context, it *postgres.Item, variants []postgres.ItemVariant) (*postgres.Product, error) {
	if it.Name == "" {
		return nil, errors.New("nama barang wajib diisi")
	}
	p, err := u.prodRepo.Create(ctx, it.Code, it.Name, it.Description, it.BasePrice, it.Unit, it.ImageURL,
		it.CategoryID, it.HasVariants, it.TaxRate, 0, 0, it.Barcode, it.PurchasePrice, it.MinStock, it.IsService)
	if err != nil {
		return nil, err
	}
	if it.HasVariants && len(variants) > 0 {
		_ = u.itemRepo.ReplaceVariants(ctx, p.ID, variants)
	}
	if it.InitialStock > 0 {
		whID, err := u.stockRepo.GetDefaultWarehouseID(ctx)
		if err != nil {
			return p, nil
		}
		if err := u.stockRepo.AddStock(ctx, p.ID, whID, it.InitialStock, "initial", "item", &p.ID, "Stok awal barang", nil); err != nil {
			return p, nil
		}
		fresh, err := u.prodRepo.GetByID(ctx, p.ID)
		if err == nil {
			return fresh, nil
		}
	}
	return p, nil
}

func (u *ItemUsecase) Update(ctx context.Context, id int64, it *postgres.Item, variants []postgres.ItemVariant) (*postgres.Product, error) {
	if it.Name == "" {
		return nil, errors.New("nama barang wajib diisi")
	}
	p, err := u.prodRepo.Update(ctx, id, it.Code, it.Name, it.Description, it.BasePrice, it.Unit, it.ImageURL,
		it.CategoryID, it.HasVariants, it.TaxRate, 0, 0, it.Barcode, it.PurchasePrice, it.MinStock, it.IsService)
	if err != nil {
		return nil, err
	}
	if it.HasVariants {
		_ = u.itemRepo.ReplaceVariants(ctx, p.ID, variants)
	} else {
		_ = u.itemRepo.ReplaceVariants(ctx, p.ID, nil)
	}
	return p, nil
}

func (u *ItemUsecase) Delete(ctx context.Context, id int64) error {
	return u.prodRepo.Delete(ctx, id)
}
