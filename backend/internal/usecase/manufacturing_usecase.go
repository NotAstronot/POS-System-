package usecase

import (
	"context"
	"database/sql"
	"errors"
	"pos-system/internal/repository/postgres"
)

type ManufacturingUsecase struct {
	bomRepo       *postgres.BomRepository
	workOrderRepo *postgres.WorkOrderRepository
	stockRepo     *postgres.StockRepository
	itemRepo      *postgres.ItemRepository
}

func NewManufacturingUsecase(
	bomRepo *postgres.BomRepository,
	workOrderRepo *postgres.WorkOrderRepository,
	stockRepo *postgres.StockRepository,
	itemRepo *postgres.ItemRepository,
) *ManufacturingUsecase {
	return &ManufacturingUsecase{
		bomRepo: bomRepo, workOrderRepo: workOrderRepo,
		stockRepo: stockRepo, itemRepo: itemRepo,
	}
}

// ==================== Bill of Materials ====================

func (u *ManufacturingUsecase) ListBoms(ctx context.Context) ([]postgres.BillOfMaterials, error) {
	return u.bomRepo.List(ctx)
}

func (u *ManufacturingUsecase) GetBom(ctx context.Context, id int64) (*postgres.BillOfMaterials, error) {
	return u.bomRepo.GetByID(ctx, id)
}

func (u *ManufacturingUsecase) CreateBom(ctx context.Context, b *postgres.BillOfMaterials) error {
	if b.ProductID == 0 {
		return errors.New("barang jadi wajib dipilih")
	}
	if len(b.Items) == 0 {
		return errors.New("minimal harus ada 1 bahan baku")
	}
	if b.QuantityOutput <= 0 {
		b.QuantityOutput = 1
	}
	if b.Unit == "" {
		b.Unit = "pcs"
	}
	if err := u.enrichBom(ctx, b); err != nil {
		return err
	}
	number, err := u.bomRepo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	b.BomNumber = number
	b.Status = "active"
	b.CostPerUnit = u.calculateCostPerUnit(b)
	id, err := u.bomRepo.Create(ctx, b)
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (u *ManufacturingUsecase) UpdateBom(ctx context.Context, b *postgres.BillOfMaterials) error {
	if b.ProductID == 0 {
		return errors.New("barang jadi wajib dipilih")
	}
	if len(b.Items) == 0 {
		return errors.New("minimal harus ada 1 bahan baku")
	}
	if err := u.enrichBom(ctx, b); err != nil {
		return err
	}
	b.CostPerUnit = u.calculateCostPerUnit(b)
	return u.bomRepo.Update(ctx, b)
}

func (u *ManufacturingUsecase) DeleteBom(ctx context.Context, id int64) error {
	return u.bomRepo.Delete(ctx, id)
}

// RecalculateBom menghitung ulang HPP (cost_per_unit) berdasarkan harga
// bahan baku terbaru dari modul Persediaan & Gudang.
func (u *ManufacturingUsecase) RecalculateBom(ctx context.Context, id int64) (*postgres.BillOfMaterials, error) {
	b, err := u.bomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("BOM tidak ditemukan")
	}
	for i := range b.Items {
		item, err := u.itemRepo.GetByID(ctx, b.Items[i].ProductID)
		if err == nil {
			b.Items[i].ProductName = item.Name
			b.Items[i].CostPerUnit = item.PurchasePrice
		}
		b.Items[i].EstimatedCost = b.Items[i].QuantityRequired * b.Items[i].CostPerUnit
	}
	b.CostPerUnit = u.calculateCostPerUnit(b)
	if err := u.bomRepo.Update(ctx, b); err != nil {
		return nil, err
	}
	return u.bomRepo.GetByID(ctx, id)
}

// enrichBom melengkapi nama bahan baku & biaya dari modul Persediaan.
func (u *ManufacturingUsecase) enrichBom(ctx context.Context, b *postgres.BillOfMaterials) error {
	if b.ProductName == "" {
		item, err := u.itemRepo.GetByID(ctx, b.ProductID)
		if err != nil {
			return errors.New("barang jadi tidak ditemukan")
		}
		b.ProductName = item.Name
	}
	for i := range b.Items {
		it := b.Items[i]
		if it.ProductID == 0 {
			return errors.New("bahan baku wajib dipilih")
		}
		item, err := u.itemRepo.GetByID(ctx, it.ProductID)
		if err != nil {
			return errors.New("bahan baku tidak ditemukan: " + it.ProductName)
		}
		it.ProductName = item.Name
		if it.Unit == "" {
			it.Unit = "pcs"
		}
		if it.CostPerUnit <= 0 {
			it.CostPerUnit = item.PurchasePrice
		}
		if it.QuantityRequired <= 0 {
			return errors.New("jumlah kebutuhan bahan baku harus lebih dari 0")
		}
		it.EstimatedCost = it.QuantityRequired * it.CostPerUnit
		b.Items[i] = it
	}
	return nil
}

// calculateCostPerUnit = (total biaya bahan baku + overhead) / jumlah output.
func (u *ManufacturingUsecase) calculateCostPerUnit(b *postgres.BillOfMaterials) float64 {
	total := float64(0)
	for _, it := range b.Items {
		total += it.EstimatedCost
	}
	total += b.OverheadCost
	if b.QuantityOutput <= 0 {
		return 0
	}
	return total / b.QuantityOutput
}

// ==================== Work Order / Perintah Kerja ====================

func (u *ManufacturingUsecase) ListWorkOrders(ctx context.Context) ([]postgres.WorkOrder, error) {
	return u.workOrderRepo.List(ctx)
}

func (u *ManufacturingUsecase) GetWorkOrder(ctx context.Context, id int64) (*postgres.WorkOrder, error) {
	return u.workOrderRepo.GetByID(ctx, id)
}

// CreateWorkOrder membuat Perintah Kerja dari BOM terpilih.
// Kebutuhan bahan baku diambil otomatis dari BOM dan diskalakan
// sesuai jumlah barang jadi yang diproduksi.
func (u *ManufacturingUsecase) CreateWorkOrder(ctx context.Context, w *postgres.WorkOrder) error {
	if w.BomID == 0 {
		return errors.New("BOM wajib dipilih")
	}
	if w.Quantity <= 0 {
		return errors.New("jumlah produksi harus lebih dari 0")
	}
	bom, err := u.bomRepo.GetByID(ctx, w.BomID)
	if err != nil {
		return errors.New("BOM tidak ditemukan")
	}
	if bom.Status != "active" {
		return errors.New("BOM tidak aktif")
	}

	scale := w.Quantity / bom.QuantityOutput
	w.BomNumber = bom.BomNumber
	w.ProductID = bom.ProductID
	w.ProductName = bom.ProductName
	w.Status = "planned"
	w.Items = make([]postgres.WorkOrderItem, 0, len(bom.Items))
	for _, bi := range bom.Items {
		w.Items = append(w.Items, postgres.WorkOrderItem{
			ProductID:        bi.ProductID,
			ProductName:      bi.ProductName,
			QuantityRequired: bi.QuantityRequired * scale,
			Unit:             bi.Unit,
			CostPerUnit:      bi.CostPerUnit,
			Subtotal:         bi.EstimatedCost * scale,
		})
	}
	number, err := u.workOrderRepo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	w.WorkOrderNumber = number
	id, err := u.workOrderRepo.Create(ctx, w)
	if err != nil {
		return err
	}
	w.ID = id
	return nil
}

// CompleteWorkOrder menyelesaikan produksi:
// 1. Kurangi stok bahan baku di gudang (modul Persediaan).
// 2. Tambah stok barang jadi sesuai hasil produksi.
// 3. Hitung biaya aktual (HPP) dan simpan ke perintah kerja.
func (u *ManufacturingUsecase) CompleteWorkOrder(ctx context.Context, id int64, createdBy *int64) error {
	w, err := u.workOrderRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("perintah kerja tidak ditemukan")
	}
	if w.Status != "planned" && w.Status != "in_progress" {
		return errors.New("hanya perintah kerja planned/in_progress yang dapat diselesaikan")
	}
	if w.WarehouseID == nil {
		return errors.New("gudang produksi wajib dipilih")
	}

	actualCost := float64(0)
	for _, it := range w.Items {
		qty := it.QuantityRequired
		cur, err := u.stockRepo.GetStock(ctx, it.ProductID, *w.WarehouseID)
		if err == sql.ErrNoRows && qty > 0 {
			return errors.New("stok bahan baku tidak tersedia: " + it.ProductName)
		}
		if err == nil && cur < qty {
			return errors.New("stok bahan baku tidak cukup: " + it.ProductName)
		}
		if err := u.stockRepo.AddStock(ctx, it.ProductID, *w.WarehouseID, -qty,
			"out", "work_order", &id, "Bahan baku untuk WO "+w.WorkOrderNumber, createdBy); err != nil {
			return err
		}
		actualCost += it.Subtotal
	}

	bom, _ := u.bomRepo.GetByID(ctx, w.BomID)
	if bom != nil {
		actualCost += bom.OverheadCost * (w.Quantity / bom.QuantityOutput)
	}

	refID := id
	if err := u.stockRepo.AddStock(ctx, w.ProductID, *w.WarehouseID, w.Quantity,
		"in", "work_order", &refID, "Hasil produksi WO "+w.WorkOrderNumber, createdBy); err != nil {
		return err
	}

	return u.workOrderRepo.UpdateStatus(ctx, id, "completed", actualCost)
}

func (u *ManufacturingUsecase) StartWorkOrder(ctx context.Context, id int64) error {
	w, err := u.workOrderRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("perintah kerja tidak ditemukan")
	}
	if w.Status != "planned" {
		return errors.New("hanya perintah kerja planned yang dapat dimulai")
	}
	return u.workOrderRepo.UpdateStatus(ctx, id, "in_progress", w.ActualCost)
}

func (u *ManufacturingUsecase) CancelWorkOrder(ctx context.Context, id int64) error {
	w, err := u.workOrderRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("perintah kerja tidak ditemukan")
	}
	if w.Status == "completed" {
		return errors.New("perintah kerja yang selesai tidak dapat dibatalkan")
	}
	return u.workOrderRepo.UpdateStatus(ctx, id, "cancelled", w.ActualCost)
}

func (u *ManufacturingUsecase) DeleteWorkOrder(ctx context.Context, id int64) error {
	return u.workOrderRepo.Delete(ctx, id)
}
