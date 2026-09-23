package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== Modul Manufaktur: Bill of Materials ====================

type BomHandler struct {
	usecase *usecase.ManufacturingUsecase
}

func NewBomHandler(u *usecase.ManufacturingUsecase) *BomHandler {
	return &BomHandler{usecase: u}
}

func (h *BomHandler) List(c *gin.Context) {
	data, err := h.usecase.ListBoms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *BomHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetBom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BOM tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *BomHandler) Create(c *gin.Context) {
	var req struct {
		ProductID      int64   `json:"product_id" binding:"required"`
		QuantityOutput float64 `json:"quantity_output"`
		Unit           string  `json:"unit"`
		OverheadCost   float64 `json:"overhead_cost"`
		Notes          string  `json:"notes"`
		Items          []struct {
			ProductID        int64   `json:"product_id" binding:"required"`
			QuantityRequired float64 `json:"quantity_required" binding:"required"`
			Unit             string  `json:"unit"`
			CostPerUnit      float64 `json:"cost_per_unit"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bom := &postgres.BillOfMaterials{
		ProductID: req.ProductID, QuantityOutput: req.QuantityOutput,
		Unit: req.Unit, OverheadCost: req.OverheadCost, Notes: req.Notes,
	}
	for _, it := range req.Items {
		bom.Items = append(bom.Items, postgres.BomItem{
			ProductID: it.ProductID, QuantityRequired: it.QuantityRequired,
			Unit: it.Unit, CostPerUnit: it.CostPerUnit,
		})
	}
	userID := jwtUserIDPtr(c)
	bom.CreatedBy = userID
	if err := h.usecase.CreateBom(c.Request.Context(), bom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bom})
}

func (h *BomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		ProductID      int64   `json:"product_id"`
		QuantityOutput float64 `json:"quantity_output"`
		Unit           string  `json:"unit"`
		OverheadCost   float64 `json:"overhead_cost"`
		Notes          string  `json:"notes"`
		Items          []struct {
			ProductID        int64   `json:"product_id"`
			QuantityRequired float64 `json:"quantity_required"`
			Unit             string  `json:"unit"`
			CostPerUnit      float64 `json:"cost_per_unit"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetBom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BOM tidak ditemukan"})
		return
	}
	if req.ProductID != 0 {
		item.ProductID = req.ProductID
		item.ProductName = ""
	}
	if req.QuantityOutput > 0 {
		item.QuantityOutput = req.QuantityOutput
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.OverheadCost >= 0 {
		item.OverheadCost = req.OverheadCost
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Items != nil {
		item.Items = make([]postgres.BomItem, 0, len(req.Items))
		for _, it := range req.Items {
			item.Items = append(item.Items, postgres.BomItem{
				ProductID: it.ProductID, QuantityRequired: it.QuantityRequired,
				Unit: it.Unit, CostPerUnit: it.CostPerUnit,
			})
		}
	}
	if err := h.usecase.UpdateBom(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *BomHandler) Recalculate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.RecalculateBom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *BomHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteBom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "BOM dihapus"})
}

// ==================== Modul Manufaktur: Perintah Kerja (Work Order) ====================

type WorkOrderHandler struct {
	usecase *usecase.ManufacturingUsecase
}

func NewWorkOrderHandler(u *usecase.ManufacturingUsecase) *WorkOrderHandler {
	return &WorkOrderHandler{usecase: u}
}

func (h *WorkOrderHandler) List(c *gin.Context) {
	data, err := h.usecase.ListWorkOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *WorkOrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetWorkOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "perintah kerja tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *WorkOrderHandler) Create(c *gin.Context) {
	var req struct {
		BomID         int64   `json:"bom_id" binding:"required"`
		Quantity      float64 `json:"quantity" binding:"required"`
		WarehouseID   *int64  `json:"warehouse_id"`
		ScheduledDate string  `json:"scheduled_date"`
		Notes         string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wo := &postgres.WorkOrder{
		BomID: req.BomID, Quantity: req.Quantity, WarehouseID: req.WarehouseID,
		ScheduledDate: req.ScheduledDate, Notes: req.Notes,
	}
	userID := jwtUserIDPtr(c)
	wo.CreatedBy = userID
	if err := h.usecase.CreateWorkOrder(c.Request.Context(), wo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wo})
}

func (h *WorkOrderHandler) Start(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.StartWorkOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "perintah kerja dimulai"})
}

func (h *WorkOrderHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := jwtUserIDPtr(c)
	if err := h.usecase.CompleteWorkOrder(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "produksi selesai, stok bahan baku terpakai dan barang jadi masuk gudang"})
}

func (h *WorkOrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.CancelWorkOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "perintah kerja dibatalkan"})
}

func (h *WorkOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteWorkOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "perintah kerja dihapus"})
}
