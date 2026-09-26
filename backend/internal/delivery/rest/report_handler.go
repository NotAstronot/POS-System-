package rest

import (
	"net/http"
	"strconv"

	"pos-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	usecase *usecase.ReportUsecase
}

func NewReportHandler(u *usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{usecase: u}
}

func (h *ReportHandler) Outlets(c *gin.Context) {
	list, err := h.usecase.Outlets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

func (h *ReportHandler) Summary(c *gin.Context) {
	period := c.DefaultQuery("period", "this_month")
	var outletID *int64
	if v := c.Query("outlet_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "outlet_id tidak valid"})
			return
		}
		outletID = &id
	}

	res, err := h.usecase.Summary(c.Request.Context(), period, outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

func (h *ReportHandler) PaymentMethods(c *gin.Context) {
	period := c.DefaultQuery("period", "this_month")
	var outletID *int64
	if v := c.Query("outlet_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "outlet_id tidak valid"})
			return
		}
		outletID = &id
	}

	res, err := h.usecase.PaymentMethods(c.Request.Context(), period, outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
