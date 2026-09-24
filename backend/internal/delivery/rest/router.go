package rest

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handlers is the collection of all HTTP handlers used by the router.
type Handlers struct {
	Auth               *AuthHandler
	Category           *CategoryHandler
	Product            *ProductHandler
	User               *UserHandler
	Upload             *UploadHandler
	Shift              *ShiftHandler
	QRIS               *QRISHandler
	Transaction        *TransactionHandler
	Branch             *BranchHandler
	Department         *DepartmentHandler
	Employee           *EmployeeHandler
	Commission         *CommissionHandler
	Salary             *SalaryHandler
	Cash               *CashHandler
	Debt               *DebtHandler
	Receivable         *ReceivableHandler
	Withdrawal         *WithdrawalHandler
	Supplier           *SupplierHandler
	PurchaseOrder      *PurchaseOrderHandler
	PurchaseInvoice    *PurchaseInvoiceHandler
	PurchasePayment    *PurchasePaymentHandler
	PurchaseReturn     *PurchaseReturnHandler
	CustomerCategory   *CustomerCategoryHandler
	SalesCategory      *SalesCategoryHandler
	Customer           *CustomerHandler
	SalesQuotation     *SalesQuotationHandler
	SalesOrder         *SalesOrderHandler
	DeliveryOrder      *DeliveryOrderHandler
	SalesInvoice       *SalesInvoiceHandler
	SalesReceipt       *SalesReceiptHandler
	SalesReturn        *SalesReturnHandler
	Warehouse          *WarehouseHandler
	Item               *ItemHandler
	StockTransfer      *StockTransferHandler
	StockAdjustment    *StockAdjustmentHandler
	PriceChange        *PriceChangeHandler
	Stock              *StockHandler
	BankTransfer       *BankTransferHandler
	BankReconciliation *BankReconciliationHandler
	ChartOfAccounts    *ChartOfAccountsHandler
	JournalEntry       *JournalEntryHandler
	FixedAsset         *FixedAssetHandler
	PeriodEnd          *PeriodEndHandler
	Order              *OrderHandler
	Marketplace        *MarketplaceHandler
	BankStatement      *BankStatementHandler
	EFaktur            *EFakturHandler
	Bom                *BomHandler
	WorkOrder          *WorkOrderHandler
	Report             *ReportHandler
	Tenant             *TenantHandler
}

// RouteConfig carries cross-cutting dependencies needed to register routes.
type RouteConfig struct {
	JWTSecret string
	DB        *sql.DB
}

// RegisterRoutes mounts every HTTP route on the given engine.
func RegisterRoutes(r *gin.Engine, h Handlers, cfg RouteConfig) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Static("/uploads", "./uploads")
	r.Static("/exports", "./exports")

	// JWTAuth validates the token; TenantContext extracts/validates tenant_id
	// and injects it into context.Context + the PostgreSQL RLS session.
	authMw := Chain(JWTAuth(cfg.JWTSecret), TenantContext(cfg.JWTSecret, cfg.DB))
	adminMw := RequireAdmin()
	posAccessMw := RequirePermission(cfg.DB, "pos_access")
	menuManageMw := RequirePermission(cfg.DB, "menu_manage")
	smartlinkMw := RequirePermission(cfg.DB, "smartlink_manage")
	manufacturingMw := RequirePermission(cfg.DB, "manufacturing_manage")

	api := r.Group("/api/v1")
	{
		registerAuthRoutes(api, h)
		registerPOSRoutes(api, h, authMw, posAccessMw, menuManageMw, cfg.DB)
		registerPaymentRoutes(api, h, authMw)
		registerShiftRoutes(api, h, authMw)
		registerTransactionRoutes(api, h, authMw, posAccessMw)
		registerAdminRoutes(api, h, authMw, adminMw)
		registerSmartlinkRoutes(api, h, authMw, adminMw, smartlinkMw)
		registerManufacturingRoutes(api, h, authMw, adminMw, manufacturingMw)
		registerSuperAdminRoutes(api, h, authMw)
	}
}

func registerAuthRoutes(api *gin.RouterGroup, h Handlers) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Auth.Login)
		auth.POST("/register-merchant", h.Auth.RegisterMerchant)
	}
}

func registerPOSRoutes(api *gin.RouterGroup, h Handlers, authMw, posAccessMw, menuManageMw gin.HandlerFunc, db *sql.DB) {
	pos := api.Group("")
	pos.Use(authMw)
	{
		pos.GET("/products/search", posAccessMw, h.Product.Search)
		pos.GET("/products/barcode", posAccessMw, h.Product.Barcode)
		pos.GET("/products/category/:categoryId", posAccessMw, h.Product.ListByCategory)
		pos.GET("/categories", posAccessMw, h.Category.List)
		pos.GET("/availability", posAccessMw, h.Product.ListAvailability)
		pos.GET("/availability/:productId", posAccessMw, h.Product.GetAvailability)
		pos.POST("/availability", posAccessMw, h.Product.SetAvailability)
		pos.GET("/products/low-stock", posAccessMw, h.Product.ListLowStock)

		pos.GET("/products", posAccessMw, h.Product.List)
		pos.POST("/products", menuManageMw, h.Product.Create)
		pos.GET("/products/:id", posAccessMw, h.Product.GetByID)
		pos.PUT("/products/:id", menuManageMw, h.Product.Update)
		pos.DELETE("/products/:id", menuManageMw, h.Product.Delete)
		pos.POST("/categories", menuManageMw, h.Category.Create)
		pos.GET("/categories/:id", posAccessMw, h.Category.GetByID)
		pos.PUT("/categories/:id", menuManageMw, h.Category.Update)
		pos.DELETE("/categories/:id", menuManageMw, h.Category.Delete)

		pos.POST("/orders", posAccessMw, h.Order.Create)
		pos.POST("/sync/orders", posAccessMw, h.Order.Sync)
		pos.GET("/orders/recent", posAccessMw, h.Order.Recent)
		pos.GET("/revenue/summary", posAccessMw, h.Order.RevenueSummary)
	}

	reports := api.Group("/reports")
	reports.Use(authMw, RequirePermission(db, "report_view"))
	{
		reports.GET("/outlets", h.Report.Outlets)
		reports.GET("/summary", h.Report.Summary)
	}
}

func registerPaymentRoutes(api *gin.RouterGroup, h Handlers, authMw gin.HandlerFunc) {
	payments := api.Group("/payments")
	payments.Use(authMw)
	{
		payments.POST("/qris", h.QRIS.Generate)
	}
}

func registerShiftRoutes(api *gin.RouterGroup, h Handlers, authMw gin.HandlerFunc) {
	shift := api.Group("/shifts")
	shift.Use(authMw)
	{
		shift.GET("/active", h.Shift.GetActiveShift)
		shift.POST("/open", h.Order.OpenShift)
		shift.POST("/:id/close", h.Order.CloseShift)
	}
}

func registerTransactionRoutes(api *gin.RouterGroup, h Handlers, authMw, posAccessMw gin.HandlerFunc) {
	transaction := api.Group("/transactions")
	transaction.Use(authMw)
	{
		transaction.GET("", h.Transaction.List)
		transaction.POST("", h.Transaction.Create)
		transaction.GET("/:id", h.Transaction.GetByID)
		transaction.POST("/sync", posAccessMw, h.Transaction.Sync)
	}
}
