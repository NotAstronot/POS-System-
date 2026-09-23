package rest

import (
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(api *gin.RouterGroup, h Handlers, authMw, adminMw gin.HandlerFunc) {
	admin := api.Group("/admin")
	admin.Use(authMw, adminMw)
	{
		admin.GET("/users", h.User.List)
		admin.GET("/users/:id", h.User.GetByID)
		admin.POST("/users", h.User.Create)
		admin.PUT("/users/:id", h.User.Update)
		admin.PUT("/users/:id/password", h.User.UpdatePassword)
		admin.DELETE("/users/:id", h.User.Delete)
		admin.GET("/users/:id/permissions", h.User.GetPermissions)
		admin.PUT("/users/:id/permissions", h.User.SetPermissions)
		admin.GET("/permissions", h.User.GetAllPermissions)
		admin.POST("/upload", h.Upload.Upload)
		admin.GET("/transactions/log", h.Transaction.GetLogsByUsername)
		admin.GET("/transactions/:id/log", h.Transaction.GetLog)

		registerAdminMasterRoutes(admin, h)
		registerAdminPurchasingRoutes(admin, h)
		registerAdminSalesRoutes(admin, h)
		registerAdminInventoryRoutes(admin, h)
		registerAdminFinanceRoutes(admin, h)
	}
}

func registerAdminMasterRoutes(admin *gin.RouterGroup, h Handlers) {
	admin.GET("/branches", h.Branch.List)
	admin.GET("/branches/:id", h.Branch.GetByID)
	admin.POST("/branches", h.Branch.Create)
	admin.PUT("/branches/:id", h.Branch.Update)
	admin.DELETE("/branches/:id", h.Branch.Delete)

	admin.GET("/departments", h.Department.List)
	admin.GET("/departments/:id", h.Department.GetByID)
	admin.POST("/departments", h.Department.Create)
	admin.PUT("/departments/:id", h.Department.Update)
	admin.DELETE("/departments/:id", h.Department.Delete)

	admin.GET("/employees", h.Employee.List)
	admin.GET("/employees/:id", h.Employee.GetByID)
	admin.POST("/employees", h.Employee.Create)
	admin.PUT("/employees/:id", h.Employee.Update)
	admin.DELETE("/employees/:id", h.Employee.Delete)

	admin.GET("/commissions", h.Commission.List)
	admin.GET("/commissions/:id", h.Commission.GetByID)
	admin.POST("/commissions", h.Commission.Create)
	admin.PUT("/commissions/:id", h.Commission.Update)
	admin.DELETE("/commissions/:id", h.Commission.Delete)

	admin.GET("/salaries", h.Salary.List)
	admin.GET("/salaries/:id", h.Salary.GetByID)
	admin.POST("/salaries", h.Salary.Create)
	admin.PUT("/salaries/:id", h.Salary.Update)
	admin.DELETE("/salaries/:id", h.Salary.Delete)
	admin.PUT("/salaries/:id/pay", h.Salary.MarkPaid)

	admin.GET("/cash", h.Cash.List)
	admin.GET("/cash/:id", h.Cash.GetByID)
	admin.POST("/cash", h.Cash.Create)
	admin.DELETE("/cash/:id", h.Cash.Delete)

	admin.GET("/debts", h.Debt.List)
	admin.GET("/debts/:id", h.Debt.GetByID)
	admin.POST("/debts", h.Debt.Create)
	admin.PUT("/debts/:id", h.Debt.Update)
	admin.DELETE("/debts/:id", h.Debt.Delete)

	admin.GET("/receivables", h.Receivable.List)
	admin.GET("/receivables/:id", h.Receivable.GetByID)
	admin.POST("/receivables", h.Receivable.Create)
	admin.PUT("/receivables/:id", h.Receivable.Update)
	admin.DELETE("/receivables/:id", h.Receivable.Delete)

	admin.GET("/withdrawals", h.Withdrawal.List)
	admin.GET("/withdrawals/:id", h.Withdrawal.GetByID)
	admin.POST("/withdrawals", h.Withdrawal.Create)
	admin.PUT("/withdrawals/:id", h.Withdrawal.Update)
	admin.DELETE("/withdrawals/:id", h.Withdrawal.Delete)
	admin.PUT("/withdrawals/:id/approve", h.Withdrawal.Approve)
	admin.PUT("/withdrawals/:id/reject", h.Withdrawal.Reject)
}

func registerAdminPurchasingRoutes(admin *gin.RouterGroup, h Handlers) {
	admin.GET("/suppliers", h.Supplier.List)
	admin.GET("/suppliers/:id", h.Supplier.GetByID)
	admin.POST("/suppliers", h.Supplier.Create)
	admin.PUT("/suppliers/:id", h.Supplier.Update)
	admin.DELETE("/suppliers/:id", h.Supplier.Delete)

	admin.GET("/purchase-orders", h.PurchaseOrder.List)
	admin.GET("/purchase-orders/:id", h.PurchaseOrder.GetByID)
	admin.POST("/purchase-orders", h.PurchaseOrder.Create)
	admin.PUT("/purchase-orders/:id", h.PurchaseOrder.Update)
	admin.DELETE("/purchase-orders/:id", h.PurchaseOrder.Delete)
	admin.PUT("/purchase-orders/:id/receive", h.PurchaseOrder.ReceiveItems)

	admin.GET("/purchase-invoices", h.PurchaseInvoice.List)
	admin.GET("/purchase-invoices/:id", h.PurchaseInvoice.GetByID)
	admin.POST("/purchase-invoices", h.PurchaseInvoice.Create)
	admin.PUT("/purchase-invoices/:id", h.PurchaseInvoice.Update)
	admin.DELETE("/purchase-invoices/:id", h.PurchaseInvoice.Delete)

	admin.GET("/purchase-payments", h.PurchasePayment.List)
	admin.GET("/purchase-payments/:id", h.PurchasePayment.GetByID)
	admin.POST("/purchase-payments", h.PurchasePayment.Create)
	admin.DELETE("/purchase-payments/:id", h.PurchasePayment.Delete)

	admin.GET("/purchase-returns", h.PurchaseReturn.List)
	admin.GET("/purchase-returns/:id", h.PurchaseReturn.GetByID)
	admin.POST("/purchase-returns", h.PurchaseReturn.Create)
	admin.PUT("/purchase-returns/:id", h.PurchaseReturn.Update)
	admin.DELETE("/purchase-returns/:id", h.PurchaseReturn.Delete)
}

func registerAdminSalesRoutes(admin *gin.RouterGroup, h Handlers) {
	admin.GET("/customer-categories", h.CustomerCategory.List)
	admin.GET("/customer-categories/:id", h.CustomerCategory.GetByID)
	admin.POST("/customer-categories", h.CustomerCategory.Create)
	admin.PUT("/customer-categories/:id", h.CustomerCategory.Update)
	admin.DELETE("/customer-categories/:id", h.CustomerCategory.Delete)

	admin.GET("/sales-categories", h.SalesCategory.List)
	admin.GET("/sales-categories/:id", h.SalesCategory.GetByID)
	admin.POST("/sales-categories", h.SalesCategory.Create)
	admin.PUT("/sales-categories/:id", h.SalesCategory.Update)
	admin.DELETE("/sales-categories/:id", h.SalesCategory.Delete)

	admin.GET("/customers", h.Customer.List)
	admin.GET("/customers/:id", h.Customer.GetByID)
	admin.POST("/customers", h.Customer.Create)
	admin.PUT("/customers/:id", h.Customer.Update)
	admin.DELETE("/customers/:id", h.Customer.Delete)

	admin.GET("/sales-quotations", h.SalesQuotation.List)
	admin.GET("/sales-quotations/:id", h.SalesQuotation.GetByID)
	admin.POST("/sales-quotations", h.SalesQuotation.Create)
	admin.PUT("/sales-quotations/:id", h.SalesQuotation.Update)
	admin.DELETE("/sales-quotations/:id", h.SalesQuotation.Delete)

	admin.GET("/sales-orders", h.SalesOrder.List)
	admin.GET("/sales-orders/:id", h.SalesOrder.GetByID)
	admin.POST("/sales-orders", h.SalesOrder.Create)
	admin.PUT("/sales-orders/:id", h.SalesOrder.Update)
	admin.DELETE("/sales-orders/:id", h.SalesOrder.Delete)

	admin.GET("/delivery-orders", h.DeliveryOrder.List)
	admin.GET("/delivery-orders/:id", h.DeliveryOrder.GetByID)
	admin.POST("/delivery-orders", h.DeliveryOrder.Create)
	admin.PUT("/delivery-orders/:id", h.DeliveryOrder.Update)
	admin.DELETE("/delivery-orders/:id", h.DeliveryOrder.Delete)

	admin.GET("/sales-invoices", h.SalesInvoice.List)
	admin.GET("/sales-invoices/:id", h.SalesInvoice.GetByID)
	admin.POST("/sales-invoices", h.SalesInvoice.Create)
	admin.PUT("/sales-invoices/:id", h.SalesInvoice.Update)
	admin.DELETE("/sales-invoices/:id", h.SalesInvoice.Delete)

	admin.GET("/sales-receipts", h.SalesReceipt.List)
	admin.GET("/sales-receipts/:id", h.SalesReceipt.GetByID)
	admin.POST("/sales-receipts", h.SalesReceipt.Create)
	admin.DELETE("/sales-receipts/:id", h.SalesReceipt.Delete)

	admin.GET("/sales-returns", h.SalesReturn.List)
	admin.GET("/sales-returns/:id", h.SalesReturn.GetByID)
	admin.POST("/sales-returns", h.SalesReturn.Create)
	admin.PUT("/sales-returns/:id", h.SalesReturn.Update)
	admin.DELETE("/sales-returns/:id", h.SalesReturn.Delete)
}

func registerAdminInventoryRoutes(admin *gin.RouterGroup, h Handlers) {
	admin.GET("/warehouses", h.Warehouse.List)
	admin.GET("/warehouses/:id", h.Warehouse.GetByID)
	admin.POST("/warehouses", h.Warehouse.Create)
	admin.PUT("/warehouses/:id", h.Warehouse.Update)
	admin.DELETE("/warehouses/:id", h.Warehouse.Delete)
	admin.GET("/warehouses/:id/stock", h.Item.WarehouseStock)

	admin.GET("/items", h.Item.List)
	admin.GET("/items/:id", h.Item.GetByID)
	admin.POST("/items", h.Item.Create)
	admin.PUT("/items/:id", h.Item.Update)
	admin.DELETE("/items/:id", h.Item.Delete)

	admin.GET("/stock-movements", h.Stock.ListMovements)

	admin.GET("/stock-transfers", h.StockTransfer.List)
	admin.GET("/stock-transfers/:id", h.StockTransfer.GetByID)
	admin.POST("/stock-transfers", h.StockTransfer.Create)
	admin.PUT("/stock-transfers/:id", h.StockTransfer.Update)
	admin.PUT("/stock-transfers/:id/send", h.StockTransfer.Send)
	admin.PUT("/stock-transfers/:id/receive", h.StockTransfer.Receive)
	admin.PUT("/stock-transfers/:id/cancel", h.StockTransfer.Cancel)
	admin.DELETE("/stock-transfers/:id", h.StockTransfer.Delete)

	admin.GET("/stock-adjustments", h.StockAdjustment.List)
	admin.GET("/stock-adjustments/:id", h.StockAdjustment.GetByID)
	admin.POST("/stock-adjustments", h.StockAdjustment.Create)
	admin.DELETE("/stock-adjustments/:id", h.StockAdjustment.Delete)

	admin.GET("/price-changes", h.PriceChange.List)
	admin.GET("/price-changes/:id", h.PriceChange.GetByID)
	admin.POST("/price-changes", h.PriceChange.Create)
	admin.DELETE("/price-changes/:id", h.PriceChange.Delete)
}

func registerAdminFinanceRoutes(admin *gin.RouterGroup, h Handlers) {
	admin.GET("/bank-transfers", h.BankTransfer.List)
	admin.GET("/bank-transfers/:id", h.BankTransfer.GetByID)
	admin.POST("/bank-transfers", h.BankTransfer.Create)
	admin.PUT("/bank-transfers/:id", h.BankTransfer.Update)
	admin.DELETE("/bank-transfers/:id", h.BankTransfer.Delete)

	admin.GET("/bank-reconciliations", h.BankReconciliation.List)
	admin.GET("/bank-reconciliations/:id", h.BankReconciliation.GetByID)
	admin.POST("/bank-reconciliations", h.BankReconciliation.Create)
	admin.DELETE("/bank-reconciliations/:id", h.BankReconciliation.Delete)

	admin.GET("/chart-of-accounts", h.ChartOfAccounts.List)
	admin.GET("/chart-of-accounts/:id", h.ChartOfAccounts.GetByID)
	admin.POST("/chart-of-accounts", h.ChartOfAccounts.Create)
	admin.PUT("/chart-of-accounts/:id", h.ChartOfAccounts.Update)
	admin.DELETE("/chart-of-accounts/:id", h.ChartOfAccounts.Delete)

	admin.GET("/journal-entries", h.JournalEntry.List)
	admin.GET("/journal-entries/:id", h.JournalEntry.GetByID)
	admin.POST("/journal-entries", h.JournalEntry.Create)
	admin.DELETE("/journal-entries/:id", h.JournalEntry.Delete)

	admin.GET("/fixed-assets", h.FixedAsset.List)
	admin.GET("/fixed-assets/schedule", h.FixedAsset.Schedule)
	admin.GET("/fixed-assets/:id", h.FixedAsset.GetByID)
	admin.POST("/fixed-assets", h.FixedAsset.Create)
	admin.PUT("/fixed-assets/:id", h.FixedAsset.Update)
	admin.DELETE("/fixed-assets/:id", h.FixedAsset.Delete)

	admin.GET("/period-end", h.PeriodEnd.List)
	admin.GET("/period-end/:id", h.PeriodEnd.GetByID)
	admin.POST("/period-end/open", h.PeriodEnd.OpenPeriod)
	admin.GET("/period-end/:id/fx-differences", h.PeriodEnd.ListFxDifferences)
	admin.POST("/period-end/fx-differences", h.PeriodEnd.CreateFxDifference)
	admin.POST("/period-end/:id/depreciation", h.PeriodEnd.RunDepreciation)
	admin.POST("/period-end/:id/close", h.PeriodEnd.ClosePeriod)
}
