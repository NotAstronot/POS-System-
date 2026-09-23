package rest

import (
	"net/http"
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
