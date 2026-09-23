package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReceivableHandler struct {
	usecase *usecase.ReceivableUsecase
}

func NewReceivableHandler(u *usecase.ReceivableUsecase) *ReceivableHandler {
	return &ReceivableHandler{usecase: u}
}

func (h *ReceivableHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *ReceivableHandler) GetByID(c *gin.Context) {
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

func (h *ReceivableHandler) Create(c *gin.Context) {
	var req struct {
		DebtorName  string  `json:"debtor_name" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.Receivable{
		DebtorName: req.DebtorName, Amount: req.Amount,
		Description: req.Description, DueDate: req.DueDate,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ReceivableHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		DebtorName  string  `json:"debtor_name"`
		Amount      float64 `json:"amount"`
		PaidAmount  float64 `json:"paid_amount"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
		Status      string  `json:"status"`
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
	if req.DebtorName != "" {
		item.DebtorName = req.DebtorName
	}
	if req.Amount > 0 {
		item.Amount = req.Amount
	}
	item.PaidAmount = req.PaidAmount
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.DueDate != nil {
		item.DueDate = req.DueDate
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ReceivableHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "piutang dihapus"})
}
