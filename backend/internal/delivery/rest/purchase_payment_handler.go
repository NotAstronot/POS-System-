package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchasePaymentHandler struct {
	usecase *usecase.PurchasePaymentUsecase
}

func NewPurchasePaymentHandler(u *usecase.PurchasePaymentUsecase) *PurchasePaymentHandler {
	return &PurchasePaymentHandler{usecase: u}
}

func (h *PurchasePaymentHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PurchasePaymentHandler) GetByID(c *gin.Context) {
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

func (h *PurchasePaymentHandler) Create(c *gin.Context) {
	var req struct {
		PurchaseInvoiceID int64   `json:"purchase_invoice_id" binding:"required"`
		SupplierID        int64   `json:"supplier_id" binding:"required"`
		PaymentDate       string  `json:"payment_date"`
		Amount            float64 `json:"amount" binding:"required"`
		PaymentMethod     string  `json:"payment_method"`
		Notes             string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.PurchasePayment{
		PurchaseInvoiceID: req.PurchaseInvoiceID,
		SupplierID:        req.SupplierID,
		PaymentDate:       req.PaymentDate,
		Amount:            req.Amount,
		PaymentMethod:     req.PaymentMethod,
		Notes:             req.Notes,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PurchasePaymentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pembayaran dihapus"})
}
