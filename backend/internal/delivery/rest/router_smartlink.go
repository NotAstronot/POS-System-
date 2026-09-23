package rest

import (
	"github.com/gin-gonic/gin"
)

func registerSmartlinkRoutes(api *gin.RouterGroup, h Handlers, authMw, adminMw, smartlinkMw gin.HandlerFunc) {
	smartlink := api.Group("/admin/smartlink")
	smartlink.Use(authMw, adminMw, smartlinkMw)
	{
		smartlink.GET("/marketplace/connections", h.Marketplace.ListConnections)
		smartlink.POST("/marketplace/connections", h.Marketplace.CreateConnection)
		smartlink.PUT("/marketplace/connections/:id", h.Marketplace.UpdateConnection)
		smartlink.DELETE("/marketplace/connections/:id", h.Marketplace.DeleteConnection)
		smartlink.GET("/marketplace/orders", h.Marketplace.ListOrders)
		smartlink.POST("/marketplace/orders/import", h.Marketplace.ImportOrders)
		smartlink.POST("/marketplace/orders/:id/post", h.Marketplace.PostOrder)
		smartlink.DELETE("/marketplace/orders/:id", h.Marketplace.DeleteOrder)

		smartlink.GET("/bank/imports", h.BankStatement.ListImports)
		smartlink.POST("/bank/imports", h.BankStatement.ImportLines)
		smartlink.POST("/bank/imports/:id/reconcile", h.BankStatement.Reconcile)
		smartlink.DELETE("/bank/imports/:id", h.BankStatement.DeleteImport)

		smartlink.GET("/tax/exports", h.EFaktur.ListExports)
		smartlink.POST("/tax/exports/generate", h.EFaktur.Generate)
		smartlink.DELETE("/tax/exports/:id", h.EFaktur.DeleteExport)
	}
}

func registerManufacturingRoutes(api *gin.RouterGroup, h Handlers, authMw, adminMw, manufacturingMw gin.HandlerFunc) {
	manufacturing := api.Group("/admin/manufacturing")
	manufacturing.Use(authMw, adminMw, manufacturingMw)
	{
		manufacturing.GET("/boms", h.Bom.List)
		manufacturing.GET("/boms/:id", h.Bom.GetByID)
		manufacturing.POST("/boms", h.Bom.Create)
		manufacturing.PUT("/boms/:id", h.Bom.Update)
		manufacturing.POST("/boms/:id/recalculate", h.Bom.Recalculate)
		manufacturing.DELETE("/boms/:id", h.Bom.Delete)

		manufacturing.GET("/work-orders", h.WorkOrder.List)
		manufacturing.GET("/work-orders/:id", h.WorkOrder.GetByID)
		manufacturing.POST("/work-orders", h.WorkOrder.Create)
		manufacturing.PUT("/work-orders/:id/start", h.WorkOrder.Start)
		manufacturing.PUT("/work-orders/:id/complete", h.WorkOrder.Complete)
		manufacturing.PUT("/work-orders/:id/cancel", h.WorkOrder.Cancel)
		manufacturing.DELETE("/work-orders/:id", h.WorkOrder.Delete)
	}
}

func registerSuperAdminRoutes(api *gin.RouterGroup, h Handlers, authMw gin.HandlerFunc) {
	superAdmin := api.Group("/super-admin")
	superAdmin.Use(authMw, RequireSuperAdmin())
	{
		superAdmin.GET("/tenants", h.Tenant.List)
		superAdmin.GET("/tenants/:id", h.Tenant.GetByID)
		superAdmin.POST("/tenants", h.Tenant.Create)
		superAdmin.PUT("/tenants/:id", h.Tenant.Update)
		superAdmin.DELETE("/tenants/:id", h.Tenant.Delete)
		superAdmin.GET("/tenants/:id/usage", h.Tenant.GetUsage)
	}
}
