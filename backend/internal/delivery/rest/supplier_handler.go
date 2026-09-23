package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	usecase *usecase.SupplierUsecase
}

func NewSupplierHandler(u *usecase.SupplierUsecase) *SupplierHandler {
	return &SupplierHandler{usecase: u}
}

func (h *SupplierHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SupplierHandler) GetByID(c *gin.Context) {
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

func (h *SupplierHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		NPWP        string `json:"npwp"`
		Address     string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.Supplier{
		Name: req.Name, ContactName: req.ContactName,
		Phone: req.Phone, Email: req.Email, NPWP: req.NPWP, Address: req.Address,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Name        string `json:"name"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		NPWP        string `json:"npwp"`
		Address     string `json:"address"`
		IsActive    *bool  `json:"is_active"`
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
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.ContactName != "" {
		item.ContactName = req.ContactName
	}
	if req.Phone != "" {
		item.Phone = req.Phone
	}
	if req.Email != "" {
		item.Email = req.Email
	}
	if req.NPWP != "" {
		item.NPWP = req.NPWP
	}
	if req.Address != "" {
		item.Address = req.Address
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

func (h *SupplierHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "supplier dihapus"})
}
