package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	usecase *usecase.CustomerUsecase
}

func NewCustomerHandler(u *usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
}

func (h *CustomerHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
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

func (h *CustomerHandler) Create(c *gin.Context) {
	var req struct {
		Name               string `json:"name" binding:"required"`
		Phone              string `json:"phone"`
		Email              string `json:"email"`
		NPWP               string `json:"npwp"`
		Address            string `json:"address"`
		CustomerCategoryID *int64 `json:"customer_category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.Customer{
		Name:               req.Name,
		Phone:              req.Phone,
		Email:              req.Email,
		NPWP:               req.NPWP,
		Address:            req.Address,
		CustomerCategoryID: req.CustomerCategoryID,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Name               string `json:"name"`
		Phone              string `json:"phone"`
		Email              string `json:"email"`
		NPWP               string `json:"npwp"`
		Address            string `json:"address"`
		CustomerCategoryID *int64 `json:"customer_category_id"`
		IsActive           *bool  `json:"is_active"`
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
	if req.CustomerCategoryID != nil {
		item.CustomerCategoryID = req.CustomerCategoryID
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

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pelanggan dihapus"})
}
