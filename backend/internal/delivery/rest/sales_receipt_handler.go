package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesReceiptHandler struct {
	usecase *usecase.SalesReceiptUsecase
}

func NewSalesReceiptHandler(u *usecase.SalesReceiptUsecase) *SalesReceiptHandler {
	return &SalesReceiptHandler{usecase: u}
}

func (h *SalesReceiptHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalesReceiptHandler) GetByID(c *gin.Context) {
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

func (h *SalesReceiptHandler) Create(c *gin.Context) {
	var req struct {
		SalesInvoiceID int64   `json:"sales_invoice_id" binding:"required"`
		CustomerID     int64   `json:"customer_id"`
		ReceiptDate    string  `json:"receipt_date"`
		Amount         float64 `json:"amount"`
		PaymentMethod  string  `json:"payment_method"`
		Notes          string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash"
	}

	rc := &postgres.SalesReceipt{
		SalesInvoiceID: req.SalesInvoiceID,
		CustomerID:     req.CustomerID,
		ReceiptDate:    req.ReceiptDate,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		Notes:          req.Notes,
	}

	if err := h.usecase.Create(c.Request.Context(), rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rc})
}

func (h *SalesReceiptHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "penerimaan penjualan dihapus"})
}
