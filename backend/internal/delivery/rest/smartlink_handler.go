package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== SmartLink e-Commerce ====================

type MarketplaceHandler struct {
	usecase *usecase.MarketplaceUsecase
}

func NewMarketplaceHandler(u *usecase.MarketplaceUsecase) *MarketplaceHandler {
	return &MarketplaceHandler{usecase: u}
}

func (h *MarketplaceHandler) ListConnections(c *gin.Context) {
	data, err := h.usecase.ListConnections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *MarketplaceHandler) CreateConnection(c *gin.Context) {
	var req struct {
		Platform           string `json:"platform" binding:"required"`
		ShopName           string `json:"shop_name" binding:"required"`
		APIToken           string `json:"api_token"`
		AccountName        string `json:"account_name"`
		CustomerID         *int64 `json:"customer_id"`
		SalesCategoryID    *int64 `json:"sales_category_id"`
		ShippingFeeAccount string `json:"shipping_fee_account"`
		CommissionAccount  string `json:"commission_account"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.MarketplaceConnection{
		Platform: req.Platform, ShopName: req.ShopName, APIToken: req.APIToken,
		AccountName: req.AccountName, CustomerID: req.CustomerID,
		SalesCategoryID: req.SalesCategoryID, ShippingFeeAccount: req.ShippingFeeAccount,
		CommissionAccount: req.CommissionAccount, IsActive: true,
	}
	if err := h.usecase.CreateConnection(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *MarketplaceHandler) UpdateConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Platform           string `json:"platform"`
		ShopName           string `json:"shop_name"`
		APIToken           string `json:"api_token"`
		AccountName        string `json:"account_name"`
		CustomerID         *int64 `json:"customer_id"`
		SalesCategoryID    *int64 `json:"sales_category_id"`
		ShippingFeeAccount string `json:"shipping_fee_account"`
		CommissionAccount  string `json:"commission_account"`
		IsActive           *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetConnection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "koneksi tidak ditemukan"})
		return
	}
	if req.Platform != "" {
		item.Platform = req.Platform
	}
	if req.ShopName != "" {
		item.ShopName = req.ShopName
	}
	if req.APIToken != "" {
		item.APIToken = req.APIToken
	}
	if req.AccountName != "" {
		item.AccountName = req.AccountName
	}
	if req.CustomerID != nil {
		item.CustomerID = req.CustomerID
	}
	if req.SalesCategoryID != nil {
		item.SalesCategoryID = req.SalesCategoryID
	}
	if req.ShippingFeeAccount != "" {
		item.ShippingFeeAccount = req.ShippingFeeAccount
	}
	if req.CommissionAccount != "" {
		item.CommissionAccount = req.CommissionAccount
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := h.usecase.UpdateConnection(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *MarketplaceHandler) DeleteConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteConnection(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "koneksi marketplace dihapus"})
}

func (h *MarketplaceHandler) ListOrders(c *gin.Context) {
	data, err := h.usecase.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *MarketplaceHandler) ImportOrders(c *gin.Context) {
	var req struct {
		ConnectionID int64 `json:"connection_id" binding:"required"`
		Orders       []struct {
			MarketplaceOrderID string  `json:"marketplace_order_id" binding:"required"`
			OrderDate          string  `json:"order_date"`
			CustomerName       string  `json:"customer_name"`
			ProductID          *int64  `json:"product_id"`
			ProductName        string  `json:"product_name" binding:"required"`
			Quantity           float64 `json:"quantity" binding:"required"`
			UnitPrice          float64 `json:"unit_price" binding:"required"`
			ShippingFee        float64 `json:"shipping_fee"`
			PlatformFee        float64 `json:"platform_fee"`
		} `json:"orders" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders := make([]postgres.MarketplaceOrder, 0, len(req.Orders))
	for _, o := range req.Orders {
		orders = append(orders, postgres.MarketplaceOrder{
			MarketplaceOrderID: o.MarketplaceOrderID,
			OrderDate:          o.OrderDate,
			CustomerName:       o.CustomerName,
			ProductID:          o.ProductID,
			ProductName:        o.ProductName,
			Quantity:           o.Quantity,
			UnitPrice:          o.UnitPrice,
			ShippingFee:        o.ShippingFee,
			PlatformFee:        o.PlatformFee,
		})
	}
	userID := jwtUserIDPtr(c)
	imported, err := h.usecase.ImportOrders(c.Request.Context(), req.ConnectionID, orders, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"imported": imported}})
}

func (h *MarketplaceHandler) PostOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := jwtUserIDPtr(c)
	if err := h.usecase.PostToSalesOrder(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pesanan diposting ke Pesanan Penjualan"})
}

func (h *MarketplaceHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pesanan marketplace dihapus"})
}

// ==================== SmartLink e-Banking ====================

type BankStatementHandler struct {
	usecase *usecase.BankStatementUsecase
}

func NewBankStatementHandler(u *usecase.BankStatementUsecase) *BankStatementHandler {
	return &BankStatementHandler{usecase: u}
}

func (h *BankStatementHandler) ListImports(c *gin.Context) {
	data, err := h.usecase.ListImports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *BankStatementHandler) ImportLines(c *gin.Context) {
	var req struct {
		AccountName   string `json:"account_name" binding:"required"`
		AccountType   string `json:"account_type"`
		StatementDate string `json:"statement_date"`
		Source        string `json:"source"`
		FileName      string `json:"file_name"`
		Lines         []struct {
			TransactionDate string  `json:"transaction_date"`
			Description     string  `json:"description"`
			Reference       string  `json:"reference"`
			Amount          float64 `json:"amount" binding:"required"`
			Type            string  `json:"type"`
		} `json:"lines" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	imp := &postgres.BankStatementImport{
		AccountName:   req.AccountName,
		AccountType:   req.AccountType,
		StatementDate: req.StatementDate,
		Source:        req.Source,
	}
	for _, l := range req.Lines {
		imp.Lines = append(imp.Lines, postgres.BankStatementImportLine{
			TransactionDate: l.TransactionDate,
			Description:     l.Description,
			Amount:          l.Amount,
			Type:            l.Type,
		})
	}
	userID := jwtUserIDPtr(c)
	imp.CreatedBy = userID
	if err := h.usecase.ImportLines(c.Request.Context(), imp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": imp})
}

func (h *BankStatementHandler) Reconcile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		BookBalance float64 `json:"book_balance"`
	}
	_ = c.ShouldBindJSON(&req)
	userID := jwtUserIDPtr(c)
	if err := h.usecase.Reconcile(c.Request.Context(), id, req.BookBalance, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rekonsiliasi bank berhasil dibuat dari mutasi"})
}

func (h *BankStatementHandler) DeleteImport(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteImport(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "impor mutasi bank dihapus"})
}

// ==================== SmartLink Tax / e-Faktur ====================

type EFakturHandler struct {
	usecase *usecase.EFakturUsecase
}

func NewEFakturHandler(u *usecase.EFakturUsecase) *EFakturHandler {
	return &EFakturHandler{usecase: u}
}

func (h *EFakturHandler) ListExports(c *gin.Context) {
	data, err := h.usecase.ListExports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *EFakturHandler) Generate(c *gin.Context) {
	var req struct {
		ExportType string `json:"export_type" binding:"required"`
		Period     string `json:"period" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := jwtUserIDPtr(c)
	item, err := h.usecase.Generate(c.Request.Context(), req.ExportType, req.Period, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *EFakturHandler) DeleteExport(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteExport(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ekspor e-Faktur dihapus"})
}
