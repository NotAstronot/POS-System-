package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseInvoiceHandler struct {
	usecase *usecase.PurchaseInvoiceUsecase
}

func NewPurchaseInvoiceHandler(u *usecase.PurchaseInvoiceUsecase) *PurchaseInvoiceHandler {
	return &PurchaseInvoiceHandler{usecase: u}
}

func (h *PurchaseInvoiceHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PurchaseInvoiceHandler) GetByID(c *gin.Context) {
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

func (h *PurchaseInvoiceHandler) Create(c *gin.Context) {
	var req struct {
		PurchaseOrderID int64   `json:"purchase_order_id" binding:"required"`
		SupplierID      int64   `json:"supplier_id" binding:"required"`
		InvoiceDate     string  `json:"invoice_date"`
		DueDate         *string `json:"due_date"`
		TotalAmount     float64 `json:"total_amount"`
		Notes           string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.PurchaseInvoice{
		PurchaseOrderID: req.PurchaseOrderID,
		SupplierID:      req.SupplierID,
		InvoiceDate:     req.InvoiceDate,
		DueDate:         req.DueDate,
		TotalAmount:     req.TotalAmount,
		Notes:           req.Notes,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PurchaseInvoiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		InvoiceDate string  `json:"invoice_date"`
		DueDate     *string `json:"due_date"`
		TotalAmount float64 `json:"total_amount"`
		PaidAmount  float64 `json:"paid_amount"`
		Status      string  `json:"status"`
		Notes       string  `json:"notes"`
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
	if req.InvoiceDate != "" {
		item.InvoiceDate = req.InvoiceDate
	}
	if req.DueDate != nil {
		item.DueDate = req.DueDate
	}
	if req.TotalAmount > 0 {
		item.TotalAmount = req.TotalAmount
	}
	item.PaidAmount = req.PaidAmount
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

func (h *PurchaseInvoiceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "faktur pembelian dihapus"})
}
