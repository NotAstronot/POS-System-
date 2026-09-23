package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesQuotationHandler struct {
	usecase *usecase.SalesQuotationUsecase
}

func NewSalesQuotationHandler(u *usecase.SalesQuotationUsecase) *SalesQuotationHandler {
	return &SalesQuotationHandler{usecase: u}
}

func (h *SalesQuotationHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalesQuotationHandler) GetByID(c *gin.Context) {
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

func (h *SalesQuotationHandler) Create(c *gin.Context) {
	var req struct {
		CustomerID      int64   `json:"customer_id" binding:"required"`
		SalesCategoryID *int64  `json:"sales_category_id"`
		QuotationDate   string  `json:"quotation_date"`
		ValidUntil      *string `json:"valid_until"`
		Notes           string  `json:"notes"`
		Items           []struct {
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			Unit        string  `json:"unit"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []postgres.SalesQuotationItem
	for _, i := range req.Items {
		items = append(items, postgres.SalesQuotationItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			Unit:        i.Unit,
			UnitPrice:   i.UnitPrice,
		})
	}

	q := &postgres.SalesQuotation{
		CustomerID:    req.CustomerID,
		QuotationDate: req.QuotationDate,
		ValidUntil:    derefString(req.ValidUntil),
		Notes:         req.Notes,
		Items:         items,
	}

	if err := h.usecase.Create(c.Request.Context(), q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": q})
}

func (h *SalesQuotationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		CustomerID    int64  `json:"customer_id"`
		QuotationDate string `json:"quotation_date"`
		ValidUntil    string `json:"valid_until"`
		Status        string `json:"status"`
		Notes         string `json:"notes"`
		Items         []struct {
			ID          int64   `json:"id"`
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			Unit        string  `json:"unit"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"items"`
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

	if req.CustomerID != 0 {
		item.CustomerID = req.CustomerID
	}
	if req.QuotationDate != "" {
		item.QuotationDate = req.QuotationDate
	}
	if req.ValidUntil != "" {
		item.ValidUntil = req.ValidUntil
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Items != nil {
		var items []postgres.SalesQuotationItem
		for _, i := range req.Items {
			items = append(items, postgres.SalesQuotationItem{
				ID:          i.ID,
				ProductName: i.ProductName,
				Quantity:    i.Quantity,
				Unit:        i.Unit,
				UnitPrice:   i.UnitPrice,
			})
		}
		item.Items = items
	}

	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SalesQuotationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "penawaran penjualan dihapus"})
}
