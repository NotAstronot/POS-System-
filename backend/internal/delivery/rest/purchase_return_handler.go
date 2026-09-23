package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseReturnHandler struct {
	usecase *usecase.PurchaseReturnUsecase
}

func NewPurchaseReturnHandler(u *usecase.PurchaseReturnUsecase) *PurchaseReturnHandler {
	return &PurchaseReturnHandler{usecase: u}
}

func (h *PurchaseReturnHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PurchaseReturnHandler) GetByID(c *gin.Context) {
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

func (h *PurchaseReturnHandler) Create(c *gin.Context) {
	var req struct {
		PurchaseOrderID int64  `json:"purchase_order_id" binding:"required"`
		SupplierID      int64  `json:"supplier_id" binding:"required"`
		ReturnDate      string `json:"return_date"`
		Reason          string `json:"reason"`
		Notes           string `json:"notes"`
		Items           []struct {
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			UnitPrice   float64 `json:"unit_price"`
			Reason      string  `json:"reason"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &postgres.PurchaseReturn{
		PurchaseOrderID: req.PurchaseOrderID,
		SupplierID:      req.SupplierID,
		ReturnDate:      req.ReturnDate,
		Reason:          req.Reason,
		Notes:           req.Notes,
	}

	var returnItems []postgres.PurchaseReturnItem
	for _, i := range req.Items {
		returnItems = append(returnItems, postgres.PurchaseReturnItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			UnitPrice:   i.UnitPrice,
			Reason:      i.Reason,
		})
	}
	item.Items = returnItems

	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PurchaseReturnHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
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
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PurchaseReturnHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pengembalian pembelian dihapus"})
}
