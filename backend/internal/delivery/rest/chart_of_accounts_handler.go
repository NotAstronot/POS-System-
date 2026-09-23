package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChartOfAccountsHandler struct {
	usecase *usecase.ChartOfAccountsUsecase
}

func NewChartOfAccountsHandler(u *usecase.ChartOfAccountsUsecase) *ChartOfAccountsHandler {
	return &ChartOfAccountsHandler{usecase: u}
}

func (h *ChartOfAccountsHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *ChartOfAccountsHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "akun tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ChartOfAccountsHandler) Create(c *gin.Context) {
	var req struct {
		Code          string `json:"code" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Category      string `json:"category" binding:"required"`
		NormalBalance string `json:"normal_balance"`
		Description   string `json:"description"`
		IsActive      *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	normalBalance := req.NormalBalance
	if normalBalance == "" {
		if req.Category == "aset" || req.Category == "beban" {
			normalBalance = "debit"
		} else {
			normalBalance = "kredit"
		}
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &postgres.ChartOfAccount{
		Code:          req.Code,
		Name:          req.Name,
		Category:      req.Category,
		NormalBalance: normalBalance,
		Description:   req.Description,
		IsActive:      isActive,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ChartOfAccountsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Code          string `json:"code" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Category      string `json:"category" binding:"required"`
		NormalBalance string `json:"normal_balance"`
		Description   string `json:"description"`
		IsActive      *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	normalBalance := req.NormalBalance
	if normalBalance == "" {
		if req.Category == "aset" || req.Category == "beban" {
			normalBalance = "debit"
		} else {
			normalBalance = "kredit"
		}
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &postgres.ChartOfAccount{
		ID:            id,
		Code:          req.Code,
		Name:          req.Name,
		Category:      req.Category,
		NormalBalance: normalBalance,
		Description:   req.Description,
		IsActive:      isActive,
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ChartOfAccountsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "akun dihapus"})
}
