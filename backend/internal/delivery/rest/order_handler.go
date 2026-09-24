package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

type OrderHandler struct {
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(u *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func jwtUserID(c *gin.Context) (int64, error) {
	raw, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("unauthorized")
	}
	switch v := raw.(type) {
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("invalid user id")
	}
}

func jwtUserIDPtr(c *gin.Context) *int64 {
	id, err := jwtUserID(c)
	if err != nil {
		return nil
	}
	return &id
}

type createOrderReq struct {
	ClientOrderID  string  `json:"client_order_id"`
	OutletID       string  `json:"outlet_id"`
	ShiftID        *int64  `json:"shift_id"`
	OrderType      string  `json:"order_type"`
	TableNumber    string  `json:"table_number"`
	CustomerName   string  `json:"customer_name"`
	CustomerPhone  string  `json:"customer_phone"`
	DiscountCode   string  `json:"discount_code"`
	DiscountAmount float64 `json:"discount_amount"`
	Items          []struct {
		ProductID    int64  `json:"product_id"`
		Quantity     int    `json:"quantity"`
		VariantLabel string `json:"variant_label"`
		Note         string `json:"note"`
	} `json:"items"`
	Payments []struct {
		Method string  `json:"method"`
		Amount float64 `json:"amount"`
		RefNo  string  `json:"reference_no"`
	} `json:"payments"`
}

func (req *createOrderReq) toInput(userID int64) postgres.CreateOrderInput {
	in := postgres.CreateOrderInput{
		ClientOrderID:  req.ClientOrderID,
		OutletID:       req.OutletID,
		UserID:         userID,
		ShiftID:        req.ShiftID,
		OrderType:      req.OrderType,
		TableNumber:    req.TableNumber,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		DiscountCode:   req.DiscountCode,
		DiscountAmount: req.DiscountAmount,
	}
	for _, it := range req.Items {
		in.Items = append(in.Items, postgres.OrderItemInput{
			ProductID:    it.ProductID,
			Quantity:     it.Quantity,
			VariantLabel: it.VariantLabel,
			Note:         it.Note,
		})
	}
	for _, p := range req.Payments {
		in.Payments = append(in.Payments, postgres.PaymentInput{
			Method: p.Method,
			Amount: p.Amount,
			RefNo:  p.RefNo,
		})
	}
	return in
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := jwtUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.CreateOrder(c.Request.Context(), req.toInput(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.writeOrderResponse(c, result)
}

func (h *OrderHandler) Sync(c *gin.Context) {
	var req struct {
		Orders []createOrderReq `json:"orders"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tidak ada order untuk disinkronkan"})
		return
	}

	userID, err := jwtUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	batch := make([]postgres.CreateOrderInput, 0, len(req.Orders))
	for _, o := range req.Orders {
		batch = append(batch, o.toInput(userID))
	}

	res := h.usecase.SyncOrders(c.Request.Context(), batch)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

func (h *OrderHandler) writeOrderResponse(c *gin.Context, result *postgres.OrderResult) {
	paymentMethod := result.Order.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "cash"
	}

	receiptItems := make([]gin.H, 0, len(result.Items))
	for _, it := range result.Items {
		receiptItems = append(receiptItems, gin.H{
			"name":      it.ProductName,
			"variant":   it.VariantLabel,
			"qty":       it.Quantity,
			"price":     it.Price,
			"sub_total": it.Subtotal,
		})
	}

	var amountPaid float64
	for _, p := range result.Payments {
		amountPaid += p.Amount
	}
	if amountPaid == 0 {
		amountPaid = result.Total
	}

	receipt := gin.H{
		"store_name":     result.StoreName,
		"store_address":  result.StoreAddr,
		"store_phone":    result.StorePhone,
		"order_number":   result.Order.OrderNumber,
		"order_type":     result.Order.OrderType,
		"table_number":   result.Order.TableNumber,
		"cashier_name":   result.CashierName,
		"items":          receiptItems,
		"sub_total":      result.Subtotal,
		"discount":       result.Discount,
		"tax":            result.Tax,
		"grand_total":    result.Total,
		"amount_paid":    amountPaid,
		"change_amount":  result.Change,
		"payment_method": paymentMethod,
		"created_at":     result.CreatedAt.Format("02/01/2006 15:04"),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"order":   result.Order,
			"receipt": receipt,
		},
	})
}

func (h *OrderHandler) Recent(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orders, err := h.usecase.RecentOrders(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) RevenueSummary(c *gin.Context) {
	period := c.DefaultQuery("period", "daily")
	summary, err := h.usecase.RevenueSummary(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	recent, err := h.usecase.RecentOrders(c.Request.Context(), 10)
	if err == nil {
		summary.RecentOrders = recent
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *OrderHandler) OpenShift(c *gin.Context) {
	var req struct {
		UserID    interface{} `json:"user_id"`
		CashStart float64     `json:"cash_start"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var uid int64
	if req.UserID != nil {
		switch v := req.UserID.(type) {
		case float64:
			uid = int64(v)
		case string:
			uid, _ = strconv.ParseInt(v, 10, 64)
		case int64:
			uid = v
		}
	}
	if uid == 0 {
		uid, _ = jwtUserID(c)
	}
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan"})
		return
	}
	log.Printf("OpenShift request: user_id=%d cash_start=%f", uid, req.CashStart)
	shift, err := h.usecase.OpenShift(c.Request.Context(), uid, req.CashStart)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":           shift.ID,
		"shift_number": shift.ID,
		"user_id":      shift.UserID,
		"cash_start":   shift.OpeningBalance,
		"status":       shift.Status,
		"opened_at":    shift.OpenedAt.Format(time.RFC3339),
	}})
}

func (h *OrderHandler) CloseShift(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("CloseShift error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		CashEnd float64 `json:"cash_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.CloseShift(c.Request.Context(), id, req.CashEnd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shift ditutup"})
}
