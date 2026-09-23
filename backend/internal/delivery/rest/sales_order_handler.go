package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesOrderHandler struct {
	usecase *usecase.SalesOrderUsecase
}

func NewSalesOrderHandler(u *usecase.SalesOrderUsecase) *SalesOrderHandler {
	return &SalesOrderHandler{usecase: u}
}

func (h *SalesOrderHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalesOrderHandler) GetByID(c *gin.Context) {
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

func (h *SalesOrderHandler) Create(c *gin.Context) {
	var req struct {
		CustomerID      int64   `json:"customer_id" binding:"required"`
		QuotationID     *int64  `json:"quotation_id"`
		SalesCategoryID *int64  `json:"sales_category_id"`
		OrderDate       string  `json:"order_date"`
		ExpectedDate    *string `json:"expected_date"`
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

	var items []postgres.SalesOrderItem
	for _, i := range req.Items {
		items = append(items, postgres.SalesOrderItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			Unit:        i.Unit,
			UnitPrice:   i.UnitPrice,
		})
	}

	o := &postgres.SalesOrder{
		CustomerID:      req.CustomerID,
		QuotationID:     req.QuotationID,
		SalesCategoryID: req.SalesCategoryID,
		OrderDate:       req.OrderDate,
		DueDate:         derefString(req.ExpectedDate),
		Notes:           req.Notes,
		Items:           items,
	}

	if err := h.usecase.Create(c.Request.Context(), o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": o})
}

func (h *SalesOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		CustomerID      int64   `json:"customer_id"`
		QuotationID     *int64  `json:"quotation_id"`
		SalesCategoryID *int64  `json:"sales_category_id"`
		OrderDate       string  `json:"order_date"`
		ExpectedDate    *string `json:"expected_date"`
		Status          string  `json:"status"`
		Notes           string  `json:"notes"`
		Items           []struct {
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
	if req.QuotationID != nil {
		item.QuotationID = req.QuotationID
	}
	if req.SalesCategoryID != nil {
		item.SalesCategoryID = req.SalesCategoryID
	}
	if req.OrderDate != "" {
		item.OrderDate = req.OrderDate
	}
	if req.ExpectedDate != nil {
		item.DueDate = *req.ExpectedDate
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Items != nil {
		var items []postgres.SalesOrderItem
		for _, i := range req.Items {
			items = append(items, postgres.SalesOrderItem{
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

func (h *SalesOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pesanan penjualan dihapus"})
}
