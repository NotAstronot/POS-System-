package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	usecase      *usecase.TransactionUsecase
	orderUsecase *usecase.OrderUsecase
}

func NewTransactionHandler(usecase *usecase.TransactionUsecase, orderUsecase *usecase.OrderUsecase) *TransactionHandler {
	return &TransactionHandler{usecase: usecase, orderUsecase: orderUsecase}
}

// Sync accepts a batch of offline cashier transactions (orders + payments) and
// persists them to PostgreSQL. Idempotent per client_order_id, so re-sending
// the queue after a network failure never duplicates data.
func (h *TransactionHandler) Sync(c *gin.Context) {
	var req struct {
		Transactions []createOrderReq `json:"transactions"`
		Orders       []createOrderReq `json:"orders"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	list := req.Transactions
	if len(list) == 0 {
		list = req.Orders
	}
	if len(list) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tidak ada transaksi untuk disinkronkan"})
		return
	}

	userID, err := jwtUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	batch := make([]postgres.CreateOrderInput, 0, len(list))
	for _, o := range list {
		batch = append(batch, o.toInput(userID))
	}

	res := h.orderUsecase.SyncOrders(c.Request.Context(), batch)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

func (h *TransactionHandler) List(c *gin.Context) {
	list, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req struct {
		OrderID       int64   `json:"order_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		PaymentMethod string  `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash"
	}
	t, err := h.usecase.Create(c.Request.Context(), req.OrderID, req.Amount, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TransactionHandler) GetLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	log, err := h.usecase.GetLogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": log})
}

func (h *TransactionHandler) GetLogsByUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	logs, err := h.usecase.GetLogsByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}
