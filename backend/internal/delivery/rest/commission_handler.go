package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommissionHandler struct {
	usecase *usecase.CommissionUsecase
}

func NewCommissionHandler(u *usecase.CommissionUsecase) *CommissionHandler {
	return &CommissionHandler{usecase: u}
}

func (h *CommissionHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *CommissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "komisi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CommissionHandler) Create(c *gin.Context) {
	var req struct {
		EmployeeID     int64   `json:"employee_id" binding:"required"`
		CommissionType string  `json:"commission_type"`
		Rate           float64 `json:"rate" binding:"required"`
		Description    string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	commType := req.CommissionType
	if commType == "" {
		commType = "percentage"
	}
	item := &postgres.Commission{
		EmployeeID: req.EmployeeID, CommissionType: commType, Rate: req.Rate, Description: req.Description,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CommissionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		EmployeeID     int64   `json:"employee_id"`
		CommissionType string  `json:"commission_type"`
		Rate           float64 `json:"rate"`
		Description    string  `json:"description"`
		IsActive       *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "komisi tidak ditemukan"})
		return
	}
	if req.EmployeeID != 0 {
		item.EmployeeID = req.EmployeeID
	}
	if req.CommissionType != "" {
		item.CommissionType = req.CommissionType
	}
	if req.Rate > 0 {
		item.Rate = req.Rate
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CommissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "komisi dihapus"})
}
