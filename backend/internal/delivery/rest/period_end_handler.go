package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PeriodEndHandler struct {
	usecase *usecase.PeriodEndUsecase
}

func NewPeriodEndHandler(u *usecase.PeriodEndUsecase) *PeriodEndHandler {
	return &PeriodEndHandler{usecase: u}
}

func (h *PeriodEndHandler) List(c *gin.Context) {
	data, err := h.usecase.ListPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PeriodEndHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "periode tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PeriodEndHandler) OpenPeriod(c *gin.Context) {
	var req struct {
		Period string `json:"period" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.OpenPeriod(c.Request.Context(), req.Period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "periode dibuka"})
}

func (h *PeriodEndHandler) ListFxDifferences(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.usecase.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "periode tidak ditemukan"})
		return
	}
	data, err := h.usecase.ListFx(c.Request.Context(), p.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PeriodEndHandler) CreateFxDifference(c *gin.Context) {
	var req struct {
		Period      string  `json:"period" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
		Description string  `json:"description"`
		GainAmount  float64 `json:"gain_amount"`
		LossAmount  float64 `json:"loss_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.FxDifference{
		Period:      req.Period,
		Currency:    req.Currency,
		Description: req.Description,
		GainAmount:  req.GainAmount,
		LossAmount:  req.LossAmount,
	}
	if err := h.usecase.CreateFx(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PeriodEndHandler) RunDepreciation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	result, err := h.usecase.RunDepreciation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PeriodEndHandler) ClosePeriod(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.ClosePeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "buku ditutup"})
}
