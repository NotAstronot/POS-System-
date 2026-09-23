package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesInvoiceHandler struct {
	usecase *usecase.SalesInvoiceUsecase
}

func NewSalesInvoiceHandler(u *usecase.SalesInvoiceUsecase) *SalesInvoiceHandler {
	return &SalesInvoiceHandler{usecase: u}
}

func (h *SalesInvoiceHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalesInvoiceHandler) GetByID(c *gin.Context) {
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

func (h *SalesInvoiceHandler) Create(c *gin.Context) {
	var req struct {
		SalesOrderID int64   `json:"sales_order_id" binding:"required"`
		CustomerID   int64   `json:"customer_id" binding:"required"`
		InvoiceDate  string  `json:"invoice_date"`
		DueDate      *string `json:"due_date"`
		TotalAmount  float64 `json:"total_amount"`
		Notes        string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	i := &postgres.SalesInvoice{
		SalesOrderID: req.SalesOrderID,
		CustomerID:   req.CustomerID,
		InvoiceDate:  req.InvoiceDate,
		DueDate:      req.DueDate,
		TotalAmount:  req.TotalAmount,
		Notes:        req.Notes,
	}

	if err := h.usecase.Create(c.Request.Context(), i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": i})
}

func (h *SalesInvoiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		SalesOrderID int64   `json:"sales_order_id"`
		CustomerID   int64   `json:"customer_id"`
		InvoiceDate  string  `json:"invoice_date"`
		DueDate      *string `json:"due_date"`
		TotalAmount  float64 `json:"total_amount"`
		Status       string  `json:"status"`
		Notes        string  `json:"notes"`
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
	if req.SalesOrderID != 0 {
		item.SalesOrderID = req.SalesOrderID
	}
	if req.CustomerID != 0 {
		item.CustomerID = req.CustomerID
	}
	if req.InvoiceDate != "" {
		item.InvoiceDate = req.InvoiceDate
	}
	if req.DueDate != nil {
		item.DueDate = req.DueDate
	}
	if req.TotalAmount != 0 {
		item.TotalAmount = req.TotalAmount
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

func (h *SalesInvoiceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "faktur penjualan dihapus"})
}
