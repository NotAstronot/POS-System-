package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FixedAssetHandler struct {
	usecase *usecase.FixedAssetUsecase
}

func NewFixedAssetHandler(u *usecase.FixedAssetUsecase) *FixedAssetHandler {
	return &FixedAssetHandler{usecase: u}
}

func (h *FixedAssetHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *FixedAssetHandler) Schedule(c *gin.Context) {
	data, err := h.usecase.Schedule(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *FixedAssetHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "aset tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *FixedAssetHandler) Create(c *gin.Context) {
	var req struct {
		AssetCode               string  `json:"asset_code" binding:"required"`
		Name                    string  `json:"name" binding:"required"`
		AssetType               string  `json:"asset_type"`
		PurchaseDate            string  `json:"purchase_date"`
		Cost                    float64 `json:"cost"`
		SalvageValue            float64 `json:"salvage_value"`
		UsefulLifeMonths        int     `json:"useful_life_months"`
		DepreciationMethod      string  `json:"depreciation_method"`
		AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
		DepreciationAccountID   int64   `json:"depreciation_account_id"`
		AccumulatedAccountID    int64   `json:"accumulated_account_id"`
		IsActive                *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetType := req.AssetType
	if assetType == "" {
		assetType = "peralatan"
	}
	method := req.DepreciationMethod
	if method == "" {
		method = "garis_lurus"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &postgres.FixedAsset{
		AssetCode:               req.AssetCode,
		Name:                    req.Name,
		AssetType:               assetType,
		PurchaseDate:            req.PurchaseDate,
		Cost:                    req.Cost,
		SalvageValue:            req.SalvageValue,
		UsefulLifeMonths:        req.UsefulLifeMonths,
		DepreciationMethod:      method,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		DepreciationAccountID:   req.DepreciationAccountID,
		AccumulatedAccountID:    req.AccumulatedAccountID,
		IsActive:                isActive,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *FixedAssetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		AssetCode               string  `json:"asset_code" binding:"required"`
		Name                    string  `json:"name" binding:"required"`
		AssetType               string  `json:"asset_type"`
		PurchaseDate            string  `json:"purchase_date"`
		Cost                    float64 `json:"cost"`
		SalvageValue            float64 `json:"salvage_value"`
		UsefulLifeMonths        int     `json:"useful_life_months"`
		DepreciationMethod      string  `json:"depreciation_method"`
		AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
		DepreciationAccountID   int64   `json:"depreciation_account_id"`
		AccumulatedAccountID    int64   `json:"accumulated_account_id"`
		IsActive                *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetType := req.AssetType
	if assetType == "" {
		assetType = "peralatan"
	}
	method := req.DepreciationMethod
	if method == "" {
		method = "garis_lurus"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &postgres.FixedAsset{
		ID:                      id,
		AssetCode:               req.AssetCode,
		Name:                    req.Name,
		AssetType:               assetType,
		PurchaseDate:            req.PurchaseDate,
		Cost:                    req.Cost,
		SalvageValue:            req.SalvageValue,
		UsefulLifeMonths:        req.UsefulLifeMonths,
		DepreciationMethod:      method,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		DepreciationAccountID:   req.DepreciationAccountID,
		AccumulatedAccountID:    req.AccumulatedAccountID,
		IsActive:                isActive,
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *FixedAssetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aset dihapus"})
}
