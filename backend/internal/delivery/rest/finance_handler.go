package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankTransferHandler struct {
	usecase *usecase.BankTransferUsecase
}

func NewBankTransferHandler(u *usecase.BankTransferUsecase) *BankTransferHandler {
	return &BankTransferHandler{usecase: u}
}

func (h *BankTransferHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *BankTransferHandler) GetByID(c *gin.Context) {
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

func (h *BankTransferHandler) Create(c *gin.Context) {
	var req struct {
		FromAccountName string  `json:"from_account_name" binding:"required"`
		FromAccountType string  `json:"from_account_type"`
		ToAccountName   string  `json:"to_account_name" binding:"required"`
		ToAccountType   string  `json:"to_account_type"`
		Amount          float64 `json:"amount" binding:"required"`
		TransferDate    string  `json:"transfer_date"`
		Notes           string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bt := &postgres.BankTransfer{
		FromAccountName: req.FromAccountName, FromAccountType: req.FromAccountType,
		ToAccountName: req.ToAccountName, ToAccountType: req.ToAccountType,
		Amount: req.Amount, TransferDate: req.TransferDate, Notes: req.Notes,
	}
	if err := h.usecase.Create(c.Request.Context(), bt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bt})
}

func (h *BankTransferHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		FromAccountName string  `json:"from_account_name"`
		FromAccountType string  `json:"from_account_type"`
		ToAccountName   string  `json:"to_account_name"`
		ToAccountType   string  `json:"to_account_type"`
		Amount          float64 `json:"amount"`
		TransferDate    string  `json:"transfer_date"`
		Notes           string  `json:"notes"`
		Status          string  `json:"status"`
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
	if req.FromAccountName != "" {
		item.FromAccountName = req.FromAccountName
	}
	if req.FromAccountType != "" {
		item.FromAccountType = req.FromAccountType
	}
	if req.ToAccountName != "" {
		item.ToAccountName = req.ToAccountName
	}
	if req.ToAccountType != "" {
		item.ToAccountType = req.ToAccountType
	}
	if req.Amount != 0 {
		item.Amount = req.Amount
	}
	if req.TransferDate != "" {
		item.TransferDate = req.TransferDate
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *BankTransferHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer bank dihapus"})
}

type BankReconciliationHandler struct {
	usecase *usecase.BankReconciliationUsecase
}

func NewBankReconciliationHandler(u *usecase.BankReconciliationUsecase) *BankReconciliationHandler {
	return &BankReconciliationHandler{usecase: u}
}

func (h *BankReconciliationHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *BankReconciliationHandler) GetByID(c *gin.Context) {
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

func (h *BankReconciliationHandler) Create(c *gin.Context) {
	var req struct {
		AccountName      string  `json:"account_name" binding:"required"`
		AccountType      string  `json:"account_type"`
		PeriodStart      string  `json:"period_start" binding:"required"`
		PeriodEnd        string  `json:"period_end" binding:"required"`
		BookBalance      float64 `json:"book_balance"`
		StatementBalance float64 `json:"statement_balance"`
		Notes            string  `json:"notes"`
		Items            []struct {
			Description     string  `json:"description"`
			TransactionDate string  `json:"transaction_date"`
			Amount          float64 `json:"amount"`
			Type            string  `json:"type"`
			IsMatched       bool    `json:"is_matched"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	br := &postgres.BankReconciliation{
		AccountName: req.AccountName, AccountType: req.AccountType,
		PeriodStart: req.PeriodStart, PeriodEnd: req.PeriodEnd,
		BookBalance: req.BookBalance, StatementBalance: req.StatementBalance,
		Notes: req.Notes,
	}
	for _, i := range req.Items {
		br.Items = append(br.Items, postgres.BankReconciliationItem{
			Description: i.Description, TransactionDate: i.TransactionDate,
			Amount: i.Amount, Type: i.Type, IsMatched: i.IsMatched,
		})
	}
	if err := h.usecase.Create(c.Request.Context(), br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": br})
}

func (h *BankReconciliationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rekonsiliasi bank dihapus"})
}
