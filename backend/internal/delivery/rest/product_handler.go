package rest

import (
	"net/http"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(usecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

func (h *ProductHandler) List(c *gin.Context) {
	list, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProductHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if len(q) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query too long"})
		return
	}
	list, err := h.usecase.Search(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ProductHandler) Barcode(c *gin.Context) {
	code := c.Query("code")
	p, err := h.usecase.Barcode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProductHandler) ListByCategory(c *gin.Context) {
	catID, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	list, err := h.usecase.ListByCategory(c.Request.Context(), catID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req struct {
		Code          string  `json:"code"`
		Name          string  `json:"name" binding:"required"`
		Description   string  `json:"description"`
		BasePrice     float64 `json:"base_price"`
		Unit          string  `json:"unit"`
		ImageURL      string  `json:"image_url"`
		CategoryID    *int64  `json:"category_id"`
		HasVariants   bool    `json:"has_variants"`
		TaxRate       float64 `json:"tax_rate"`
		SortOrder     int     `json:"sort_order"`
		Stock         int     `json:"stock"`
		Barcode       string  `json:"barcode"`
		PurchasePrice float64 `json:"purchase_price"`
		MinStock      int     `json:"min_stock"`
		IsService     bool    `json:"is_service"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.usecase.Create(c.Request.Context(), req.Code, req.Name, req.Description, req.BasePrice, req.Unit, req.ImageURL, req.CategoryID, req.HasVariants, req.TaxRate, req.SortOrder, req.Stock, req.Barcode, req.PurchasePrice, req.MinStock, req.IsService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Code          string  `json:"code"`
		Name          string  `json:"name" binding:"required"`
		Description   string  `json:"description"`
		BasePrice     float64 `json:"base_price"`
		Unit          string  `json:"unit"`
		ImageURL      string  `json:"image_url"`
		CategoryID    *int64  `json:"category_id"`
		HasVariants   bool    `json:"has_variants"`
		TaxRate       float64 `json:"tax_rate"`
		SortOrder     int     `json:"sort_order"`
		Stock         int     `json:"stock"`
		Barcode       string  `json:"barcode"`
		PurchasePrice float64 `json:"purchase_price"`
		MinStock      int     `json:"min_stock"`
		IsService     bool    `json:"is_service"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.usecase.Update(c.Request.Context(), id, req.Code, req.Name, req.Description, req.BasePrice, req.Unit, req.ImageURL, req.CategoryID, req.HasVariants, req.TaxRate, req.SortOrder, req.Stock, req.Barcode, req.PurchasePrice, req.MinStock, req.IsService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) SetAvailability(c *gin.Context) {
	var req struct {
		ProductID   int64  `json:"product_id" binding:"required"`
		Date        string `json:"date" binding:"required"`
		IsAvailable bool   `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.usecase.SetAvailability(c.Request.Context(), req.ProductID, req.Date, req.IsAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

func (h *ProductHandler) GetAvailability(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	date := c.Query("date")
	a, err := h.usecase.GetAvailability(c.Request.Context(), productID, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

func (h *ProductHandler) ListAvailability(c *gin.Context) {
	date := c.Query("date")
	list, err := h.usecase.ListAvailability(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ProductHandler) ListLowStock(c *gin.Context) {
	threshold := 5
	if t := c.Query("threshold"); t != "" {
		if v, err := strconv.Atoi(t); err == nil && v > 0 {
			threshold = v
		}
	}
	list, err := h.usecase.ListLowStock(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
