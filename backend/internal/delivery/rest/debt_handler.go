package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DebtHandler struct {
	usecase *usecase.DebtUsecase
}

func NewDebtHandler(u *usecase.DebtUsecase) *DebtHandler {
	return &DebtHandler{usecase: u}
}

func (h *DebtHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DebtHandler) GetByID(c *gin.Context) {
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

func (h *DebtHandler) Create(c *gin.Context) {
	var req struct {
		CreditorName string  `json:"creditor_name" binding:"required"`
		Amount       float64 `json:"amount" binding:"required"`
		Description  string  `json:"description"`
		DueDate      *string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.Debt{
		CreditorName: req.CreditorName, Amount: req.Amount,
		Description: req.Description, DueDate: req.DueDate,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *DebtHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		CreditorName string  `json:"creditor_name"`
		Amount       float64 `json:"amount"`
		PaidAmount   float64 `json:"paid_amount"`
		Description  string  `json:"description"`
		DueDate      *string `json:"due_date"`
		Status       string  `json:"status"`
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
	if req.CreditorName != "" {
		item.CreditorName = req.CreditorName
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

func (h *DebtHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "hutang dihapus"})
}
