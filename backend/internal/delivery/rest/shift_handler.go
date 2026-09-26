package rest

import (
	"net/http"
	"strconv"
	"pos-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	usecase *usecase.ShiftUsecase
}

func NewShiftHandler(usecase *usecase.ShiftUsecase) *ShiftHandler {
	return &ShiftHandler{usecase: usecase}
}

func (h *ShiftHandler) GetActiveShift(c *gin.Context) {
	shift, err := h.usecase.GetActiveShift(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active shift"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *ShiftHandler) History(c *gin.Context) {
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	list, err := h.usecase.History(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}
