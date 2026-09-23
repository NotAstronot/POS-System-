package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseOrderHandler struct {
	usecase *usecase.PurchaseOrderUsecase
}

func NewPurchaseOrderHandler(u *usecase.PurchaseOrderUsecase) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{usecase: u}
}

func (h *PurchaseOrderHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PurchaseOrderHandler) GetByID(c *gin.Context) {
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

func (h *PurchaseOrderHandler) Create(c *gin.Context) {
	var req struct {
		SupplierID   int64   `json:"supplier_id" binding:"required"`
		OrderDate    string  `json:"order_date"`
		ExpectedDate *string `json:"expected_date"`
		Notes        string  `json:"notes"`
		Items        []struct {
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			Unit        string  `json:"unit"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []postgres.PurchaseItem
	for _, i := range req.Items {
		items = append(items, postgres.PurchaseItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			Unit:        i.Unit,
			UnitPrice:   i.UnitPrice,
		})
	}

	po := &postgres.PurchaseOrder{
		SupplierID:   req.SupplierID,
		OrderDate:    req.OrderDate,
		ExpectedDate: req.ExpectedDate,
		Notes:        req.Notes,
		Items:        items,
	}

	if err := h.usecase.Create(c.Request.Context(), po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": po})
}

func (h *PurchaseOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		SupplierID   int64   `json:"supplier_id"`
		OrderDate    string  `json:"order_date"`
		ExpectedDate *string `json:"expected_date"`
		Status       string  `json:"status"`
		Notes        string  `json:"notes"`
		Items        []struct {
			ID          int64   `json:"id"`
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			Unit        string  `json:"unit"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"items"`
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

	if req.SupplierID != 0 {
		item.SupplierID = req.SupplierID
	}
	if req.OrderDate != "" {
		item.OrderDate = req.OrderDate
	}
	if req.ExpectedDate != nil {
		item.ExpectedDate = req.ExpectedDate
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Items != nil {
		var items []postgres.PurchaseItem
		for _, i := range req.Items {
			items = append(items, postgres.PurchaseItem{
				ID:          i.ID,
				ProductName: i.ProductName,
				Quantity:    i.Quantity,
				Unit:        i.Unit,
				UnitPrice:   i.UnitPrice,
			})
		}
		item.Items = items
	}

	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PurchaseOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "purchase order dihapus"})
}

func (h *PurchaseOrderHandler) ReceiveItems(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		WarehouseID int64 `json:"warehouse_id"`
		Items       []struct {
			ID       int64   `json:"id"`
			Quantity float64 `json:"quantity"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []postgres.PurchaseItem
	for _, i := range req.Items {
		items = append(items, postgres.PurchaseItem{
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	if err := h.usecase.ReceiveItems(c.Request.Context(), id, req.WarehouseID, items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item berhasil diterima"})
}
