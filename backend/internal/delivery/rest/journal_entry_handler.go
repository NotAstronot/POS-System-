package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JournalEntryHandler struct {
	usecase *usecase.JournalEntryUsecase
}

func NewJournalEntryHandler(u *usecase.JournalEntryUsecase) *JournalEntryHandler {
	return &JournalEntryHandler{usecase: u}
}

func (h *JournalEntryHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *JournalEntryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jurnal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *JournalEntryHandler) Create(c *gin.Context) {
	var req struct {
		EntryDate   string `json:"entry_date"`
		Description string `json:"description"`
		Reference   string `json:"reference"`
		Items       []struct {
			AccountID int64   `json:"account_id"`
			Debit     float64 `json:"debit"`
			Credit    float64 `json:"credit"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items := make([]postgres.JournalEntryItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, postgres.JournalEntryItem{
			AccountID: it.AccountID,
			Debit:     it.Debit,
			Credit:    it.Credit,
		})
	}
	entry := &postgres.JournalEntry{
		EntryDate:   req.EntryDate,
		Description: req.Description,
		Reference:   req.Reference,
	}
	if err := h.usecase.Create(c.Request.Context(), entry, items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entry})
}

func (h *JournalEntryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "jurnal dihapus"})
}
