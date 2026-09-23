package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesReturnHandler struct {
	usecase *usecase.SalesReturnUsecase
}

func NewSalesReturnHandler(u *usecase.SalesReturnUsecase) *SalesReturnHandler {
	return &SalesReturnHandler{usecase: u}
}

func (h *SalesReturnHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalesReturnHandler) GetByID(c *gin.Context) {
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

func (h *SalesReturnHandler) Create(c *gin.Context) {
	var req struct {
		SalesOrderID int64  `json:"sales_order_id" binding:"required"`
		CustomerID   int64  `json:"customer_id"`
		ReturnDate   string `json:"return_date"`
		Reason       string `json:"reason"`
		Notes        string `json:"notes"`
		Items        []struct {
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			UnitPrice   float64 `json:"unit_price"`
			Reason      string  `json:"reason"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []postgres.SalesReturnItem
	for _, i := range req.Items {
		items = append(items, postgres.SalesReturnItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			UnitPrice:   i.UnitPrice,
		})
	}

	rt := &postgres.SalesReturn{
		SalesOrderID: req.SalesOrderID,
		CustomerID:   req.CustomerID,
		ReturnDate:   req.ReturnDate,
		Notes:        req.Reason, // map Reason to Notes
		Items:        items,
	}

	if err := h.usecase.Create(c.Request.Context(), rt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rt})
}

func (h *SalesReturnHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		ReturnDate string `json:"return_date"`
		Reason     string `json:"reason"`
		Status     string `json:"status"`
		Notes      string `json:"notes"`
		Items      []struct {
			ID          int64   `json:"id"`
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			UnitPrice   float64 `json:"unit_price"`
			Reason      string  `json:"reason"`
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

	if req.ReturnDate != "" {
		item.ReturnDate = req.ReturnDate
	}
	if req.Reason != "" {
		item.Notes = req.Reason
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Items != nil {
		var items []postgres.SalesReturnItem
		for _, i := range req.Items {
			items = append(items, postgres.SalesReturnItem{
				ID:          i.ID,
				ProductName: i.ProductName,
				Quantity:    i.Quantity,
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

func (h *SalesReturnHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "retur penjualan dihapus"})
}
