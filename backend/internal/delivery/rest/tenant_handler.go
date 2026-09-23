package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	usecase *usecase.TenantUsecase
}

func NewTenantHandler(usecase *usecase.TenantUsecase) *TenantHandler {
	return &TenantHandler{usecase: usecase}
}

func (h *TenantHandler) List(c *gin.Context) {
	list, err := h.usecase.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TenantHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	tenant, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tenant})
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req struct {
		Name             string `json:"name" binding:"required"`
		Slug             string `json:"slug"`
		Domain           string `json:"domain"`
		LogoURL          string `json:"logo_url"`
		SubscriptionPlan string `json:"subscription_plan"`
		MaxUsers         int    `json:"max_users"`
		MaxProducts      int    `json:"max_products"`
		MaxBranches      int    `json:"max_branches"`
		Settings         string `json:"settings"`
		Status           string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant := &postgres.Tenant{
		Name:             req.Name,
		Slug:             req.Slug,
		Domain:           req.Domain,
		LogoURL:          req.LogoURL,
		SubscriptionPlan: req.SubscriptionPlan,
		MaxUsers:         req.MaxUsers,
		MaxProducts:      req.MaxProducts,
		MaxBranches:      req.MaxBranches,
		Settings:         req.Settings,
		Status:           req.Status,
	}

	result, err := h.usecase.Create(tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *TenantHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Name             string `json:"name"`
		Slug             string `json:"slug"`
		Domain           string `json:"domain"`
		LogoURL          string `json:"logo_url"`
		SubscriptionPlan string `json:"subscription_plan"`
		MaxUsers         int    `json:"max_users"`
		MaxProducts      int    `json:"max_products"`
		MaxBranches      int    `json:"max_branches"`
		Settings         string `json:"settings"`
		Status           string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant := &postgres.Tenant{
		ID:               id,
		Name:             req.Name,
		Slug:             req.Slug,
		Domain:           req.Domain,
		LogoURL:          req.LogoURL,
		SubscriptionPlan: req.SubscriptionPlan,
		MaxUsers:         req.MaxUsers,
		MaxProducts:      req.MaxProducts,
		MaxBranches:      req.MaxBranches,
		Settings:         req.Settings,
		Status:           req.Status,
	}

	result, err := h.usecase.Update(tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tenant berhasil dihapus"})
}

func (h *TenantHandler) GetUsage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	usage, err := h.usecase.GetUsage(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": usage})
}
