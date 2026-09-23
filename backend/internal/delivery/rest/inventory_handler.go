package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	usecase *usecase.WarehouseUsecase
}

func NewWarehouseHandler(u *usecase.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{usecase: u}
}

func (h *WarehouseHandler) List(c *gin.Context) {
	data, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *WarehouseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		Name     string `json:"name" binding:"required"`
		BranchID *int64 `json:"branch_id"`
		Address  string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w := &postgres.Warehouse{Code: req.Code, Name: req.Name, BranchID: req.BranchID, Address: req.Address}
	if err := h.usecase.Create(c.Request.Context(), w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": w})
}

func (h *WarehouseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		BranchID *int64 `json:"branch_id"`
		Address  string `json:"address"`
		IsActive *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	if req.Code != "" {
		item.Code = req.Code
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.BranchID != nil {
		item.BranchID = req.BranchID
	}
	if req.Address != "" {
		item.Address = req.Address
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *WarehouseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "gudang dihapus"})
}

type ItemHandler struct {
	usecase *usecase.ItemUsecase
}

func NewItemHandler(u *usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{usecase: u}
}

func (h *ItemHandler) List(c *gin.Context) {
	data, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *ItemHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ItemHandler) WarehouseStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	data, err := h.usecase.ListWarehouseStock(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

type variantReq struct {
	Name            string  `json:"name"`
	Barcode         string  `json:"barcode"`
	AdditionalPrice float64 `json:"additional_price"`
	Stock           int     `json:"stock"`
	IsActive        *bool   `json:"is_active"`
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req struct {
		Code          string       `json:"code"`
		Name          string       `json:"name" binding:"required"`
		Barcode       string       `json:"barcode"`
		Description   string       `json:"description"`
		CategoryID    *int64       `json:"category_id"`
		PurchasePrice float64      `json:"purchase_price"`
		BasePrice     float64      `json:"base_price"`
		Unit          string       `json:"unit"`
		MinStock      int          `json:"min_stock"`
		IsService     bool         `json:"is_service"`
		HasVariants   bool         `json:"has_variants"`
		TaxRate       float64      `json:"tax_rate"`
		InitialStock  float64      `json:"initial_stock"`
		Variants      []variantReq `json:"variants"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it := &postgres.Item{
		Code: req.Code, Name: req.Name, Barcode: req.Barcode, Description: req.Description,
		CategoryID: req.CategoryID, PurchasePrice: req.PurchasePrice, BasePrice: req.BasePrice,
		Unit: req.Unit, MinStock: req.MinStock, IsService: req.IsService,
		HasVariants: req.HasVariants, TaxRate: req.TaxRate, InitialStock: req.InitialStock,
	}
	var variants []postgres.ItemVariant
	for _, v := range req.Variants {
		isActive := true
		if v.IsActive != nil {
			isActive = *v.IsActive
		}
		variants = append(variants, postgres.ItemVariant{Name: v.Name, Barcode: v.Barcode, AdditionalPrice: v.AdditionalPrice, Stock: v.Stock, IsActive: isActive})
	}
	p, err := h.usecase.Create(c.Request.Context(), it, variants)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ItemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Code          string       `json:"code"`
		Name          string       `json:"name" binding:"required"`
		Barcode       string       `json:"barcode"`
		Description   string       `json:"description"`
		CategoryID    *int64       `json:"category_id"`
		PurchasePrice float64      `json:"purchase_price"`
		BasePrice     float64      `json:"base_price"`
		Unit          string       `json:"unit"`
		MinStock      int          `json:"min_stock"`
		IsService     bool         `json:"is_service"`
		HasVariants   bool         `json:"has_variants"`
		TaxRate       float64      `json:"tax_rate"`
		Variants      []variantReq `json:"variants"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it := &postgres.Item{
		Code: req.Code, Name: req.Name, Barcode: req.Barcode, Description: req.Description,
		CategoryID: req.CategoryID, PurchasePrice: req.PurchasePrice, BasePrice: req.BasePrice,
		Unit: req.Unit, MinStock: req.MinStock, IsService: req.IsService,
		HasVariants: req.HasVariants, TaxRate: req.TaxRate,
	}
	var variants []postgres.ItemVariant
	for _, v := range req.Variants {
		isActive := true
		if v.IsActive != nil {
			isActive = *v.IsActive
		}
		variants = append(variants, postgres.ItemVariant{Name: v.Name, Barcode: v.Barcode, AdditionalPrice: v.AdditionalPrice, Stock: v.Stock, IsActive: isActive})
	}
	p, err := h.usecase.Update(c.Request.Context(), id, it, variants)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "barang dihapus"})
}

type StockTransferHandler struct {
	usecase *usecase.StockTransferUsecase
}

func NewStockTransferHandler(u *usecase.StockTransferUsecase) *StockTransferHandler {
	return &StockTransferHandler{usecase: u}
}

func (h *StockTransferHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StockTransferHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

type transferItemReq struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

func (h *StockTransferHandler) Create(c *gin.Context) {
	var req struct {
		FromWarehouseID int64             `json:"from_warehouse_id" binding:"required"`
		ToWarehouseID   int64             `json:"to_warehouse_id" binding:"required"`
		TransferDate    string            `json:"transfer_date"`
		Notes           string            `json:"notes"`
		Items           []transferItemReq `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var items []postgres.StockTransferItem
	for _, i := range req.Items {
		items = append(items, postgres.StockTransferItem{ProductID: i.ProductID, ProductName: i.ProductName, Quantity: i.Quantity, Unit: i.Unit})
	}
	t := &postgres.StockTransfer{
		FromWarehouseID: req.FromWarehouseID, ToWarehouseID: req.ToWarehouseID,
		TransferDate: req.TransferDate, Notes: req.Notes, Items: items,
	}
	if err := h.usecase.Create(c.Request.Context(), t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func (h *StockTransferHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		FromWarehouseID int64             `json:"from_warehouse_id"`
		ToWarehouseID   int64             `json:"to_warehouse_id"`
		TransferDate    string            `json:"transfer_date"`
		Notes           string            `json:"notes"`
		Items           []transferItemReq `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	if item.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hanya transfer draft yang bisa diedit"})
		return
	}
	if req.FromWarehouseID != 0 {
		item.FromWarehouseID = req.FromWarehouseID
	}
	if req.ToWarehouseID != 0 {
		item.ToWarehouseID = req.ToWarehouseID
	}
	if req.TransferDate != "" {
		item.TransferDate = req.TransferDate
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if len(req.Items) > 0 {
		item.Items = nil
		for _, i := range req.Items {
			item.Items = append(item.Items, postgres.StockTransferItem{ProductID: i.ProductID, ProductName: i.ProductName, Quantity: i.Quantity, Unit: i.Unit})
		}
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *StockTransferHandler) Send(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Send(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer dikirim"})
}

func (h *StockTransferHandler) Receive(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Receive(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer diterima"})
}

func (h *StockTransferHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Cancel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer dibatalkan"})
}

func (h *StockTransferHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer dihapus"})
}

type StockAdjustmentHandler struct {
	usecase *usecase.StockAdjustmentUsecase
}

func NewStockAdjustmentHandler(u *usecase.StockAdjustmentUsecase) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{usecase: u}
}

func (h *StockAdjustmentHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StockAdjustmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *StockAdjustmentHandler) Create(c *gin.Context) {
	var req struct {
		WarehouseID    int64  `json:"warehouse_id" binding:"required"`
		AdjustmentDate string `json:"adjustment_date"`
		Notes          string `json:"notes"`
		Items          []struct {
			ProductID   int64   `json:"product_id"`
			ProductName string  `json:"product_name"`
			SystemQty   float64 `json:"system_qty"`
			ActualQty   float64 `json:"actual_qty"`
			Reason      string  `json:"reason"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := &postgres.StockAdjustment{
		WarehouseID: req.WarehouseID, AdjustmentDate: req.AdjustmentDate, Notes: req.Notes,
	}
	for _, i := range req.Items {
		a.Items = append(a.Items, postgres.StockAdjustmentItem{
			ProductID: i.ProductID, ProductName: i.ProductName,
			SystemQty: i.SystemQty, ActualQty: i.ActualQty,
			Difference: i.ActualQty - i.SystemQty, Reason: i.Reason,
		})
	}
	if err := h.usecase.Create(c.Request.Context(), a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

func (h *StockAdjustmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "penyesuaian dihapus"})
}

type PriceChangeHandler struct {
	usecase *usecase.PriceChangeUsecase
}

func NewPriceChangeHandler(u *usecase.PriceChangeUsecase) *PriceChangeHandler {
	return &PriceChangeHandler{usecase: u}
}

func (h *PriceChangeHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PriceChangeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PriceChangeHandler) Create(c *gin.Context) {
	var req struct {
		ChangeDate string `json:"change_date"`
		Notes      string `json:"notes"`
		Items      []struct {
			ProductID   int64   `json:"product_id"`
			ProductName string  `json:"product_name"`
			OldPrice    float64 `json:"old_price"`
			NewPrice    float64 `json:"new_price"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &postgres.PriceChange{ChangeDate: req.ChangeDate, Notes: req.Notes}
	for _, i := range req.Items {
		p.Items = append(p.Items, postgres.PriceChangeItem{
			ProductID: i.ProductID, ProductName: i.ProductName, OldPrice: i.OldPrice, NewPrice: i.NewPrice,
		})
	}
	if err := h.usecase.Create(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *PriceChangeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "penyesuaian harga dihapus"})
}

type StockHandler struct {
	usecase *usecase.StockUsecase
}

func NewStockHandler(u *usecase.StockUsecase) *StockHandler {
	return &StockHandler{usecase: u}
}

func (h *StockHandler) ListMovements(c *gin.Context) {
	data, err := h.usecase.ListMovements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
