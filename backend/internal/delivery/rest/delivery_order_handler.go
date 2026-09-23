package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeliveryOrderHandler struct {
	usecase *usecase.DeliveryOrderUsecase
}

func NewDeliveryOrderHandler(u *usecase.DeliveryOrderUsecase) *DeliveryOrderHandler {
	return &DeliveryOrderHandler{usecase: u}
}

func (h *DeliveryOrderHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DeliveryOrderHandler) GetByID(c *gin.Context) {
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

func (h *DeliveryOrderHandler) Create(c *gin.Context) {
	var req struct {
		SalesOrderID int64  `json:"sales_order_id" binding:"required"`
		CustomerID   int64  `json:"customer_id"`
		DeliveryDate string `json:"delivery_date"`
		Notes        string `json:"notes"`
		Items        []struct {
			ProductName string  `json:"product_name"`
			Quantity    float64 `json:"quantity"`
			Unit        string  `json:"unit"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []postgres.DeliveryOrderItem
	for _, i := range req.Items {
		items = append(items, postgres.DeliveryOrderItem{
			ProductName: i.ProductName,
			Quantity:    i.Quantity,
			Unit:        i.Unit,
		})
	}

	d := &postgres.DeliveryOrder{
		SalesOrderID: req.SalesOrderID,
		CustomerID:   req.CustomerID,
		DeliveryDate: req.DeliveryDate,
		Notes:        req.Notes,
		Items:        items,
	}

	if err := h.usecase.Create(c.Request.Context(), d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (h *DeliveryOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		DeliveryDate string `json:"delivery_date"`
		Status       string `json:"status"`
		Notes        string `json:"notes"`
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
	if req.DeliveryDate != "" {
		item.DeliveryDate = req.DeliveryDate
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *DeliveryOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "delivery order dihapus"})
}
