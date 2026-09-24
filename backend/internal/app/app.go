package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"pos-system/internal/config"
	"pos-system/internal/delivery/rest"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
)

// App is the composition root: it wires repositories, use cases, handlers,
// and the HTTP router according to Clean Architecture dependency rules.
type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *gin.Engine
}

// New builds the application graph from configuration and opens the database.
func New(cfgPath string) (*App, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}

	handlers := wire(db, cfg)

	r := gin.Default()
	rest.RegisterRoutes(r, handlers, rest.RouteConfig{
		JWTSecret: cfg.JWT.SecretKey,
		DB:        db,
	})

	return &App{Config: cfg, DB: db, Router: r}, nil
}

// Close releases the database connection pool.
func (a *App) Close() {
	if a.DB != nil {
		_ = a.DB.Close()
	}
}

func connectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := "host=" + cfg.Database.Host +
		" port=" + cfg.Database.Port +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.DBName +
		" sslmode=" + cfg.Database.SSLMode +
		" connect_timeout=10"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	var rlsTest string
	if err := db.QueryRow("SHOW row_security").Scan(&rlsTest); err != nil {
		log.Printf("warning: could not check RLS status: %v", err)
	} else {
		log.Printf("RLS status: %s", rlsTest)
	}

	return db, nil
}

// wire constructs repositories → use cases → handlers (dependency direction:
// delivery → usecase → repository → database).
func wire(db *sql.DB, cfg *config.Config) rest.Handlers {
	// Repositories (framework & drivers layer)
	categoryRepo := postgres.NewCategoryRepository(db)
	productRepo := postgres.NewProductRepository(db)
	userRepo := postgres.NewUserRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	shiftRepo := postgres.NewShiftRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	orderItemRepo := postgres.NewOrderItemRepository(db)
	branchRepo := postgres.NewBranchRepository(db)
	departmentRepo := postgres.NewDepartmentRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)
	commissionRepo := postgres.NewCommissionRepository(db)
	salaryRepo := postgres.NewSalaryRepository(db)
	cashRepo := postgres.NewCashTransactionRepository(db)
	debtRepo := postgres.NewDebtRepository(db)
	receivableRepo := postgres.NewReceivableRepository(db)
	withdrawalRepo := postgres.NewWithdrawalRepository(db)
	supplierRepo := postgres.NewSupplierRepository(db)
	purchaseOrderRepo := postgres.NewPurchaseOrderRepository(db)
	purchaseInvoiceRepo := postgres.NewPurchaseInvoiceRepository(db)
	purchasePaymentRepo := postgres.NewPurchasePaymentRepository(db)
	purchaseReturnRepo := postgres.NewPurchaseReturnRepository(db)
	customerCategoryRepo := postgres.NewCustomerCategoryRepository(db)
	salesCategoryRepo := postgres.NewSalesCategoryRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	salesQuotationRepo := postgres.NewSalesQuotationRepository(db)
	salesOrderRepo := postgres.NewSalesOrderRepository(db)
	deliveryOrderRepo := postgres.NewDeliveryOrderRepository(db)
	salesInvoiceRepo := postgres.NewSalesInvoiceRepository(db)
	salesReceiptRepo := postgres.NewSalesReceiptRepository(db)
	salesReturnRepo := postgres.NewSalesReturnRepository(db)
	warehouseRepo := postgres.NewWarehouseRepository(db)
	itemRepo := postgres.NewItemRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	transferRepo := postgres.NewStockTransferRepository(db)
	adjustmentRepo := postgres.NewStockAdjustmentRepository(db)
	priceChangeRepo := postgres.NewPriceChangeRepository(db)
	bankTransferRepo := postgres.NewBankTransferRepository(db)
	bankReconciliationRepo := postgres.NewBankReconciliationRepository(db)
	chartOfAccountsRepo := postgres.NewChartOfAccountRepository(db)
	journalEntryRepo := postgres.NewJournalEntryRepository(db)
	fixedAssetRepo := postgres.NewFixedAssetRepository(db)
	periodEndRepo := postgres.NewPeriodEndRepository(db)
	revenueRepo := postgres.NewRevenueRepository(db)
	marketplaceRepo := postgres.NewMarketplaceRepository(db)
	bankStatementRepo := postgres.NewBankStatementRepository(db)
	efakturRepo := postgres.NewEfakturRepository(db)
	bomRepo := postgres.NewBomRepository(db)
	workOrderRepo := postgres.NewWorkOrderRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	tenantRepo := postgres.NewTenantRepository(db)

	// Use cases (application business rules layer)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, tenantRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	shiftUsecase := usecase.NewShiftUsecase(shiftRepo, orderRepo, orderItemRepo, productRepo)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, orderRepo)
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	departmentUsecase := usecase.NewDepartmentUsecase(departmentRepo)
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo)
	commissionUsecase := usecase.NewCommissionUsecase(commissionRepo)
	salaryUsecase := usecase.NewSalaryUsecase(salaryRepo)
	cashUsecase := usecase.NewCashUsecase(cashRepo)
	debtUsecase := usecase.NewDebtUsecase(debtRepo)
	receivableUsecase := usecase.NewReceivableUsecase(receivableRepo)
	withdrawalUsecase := usecase.NewWithdrawalUsecase(withdrawalRepo)
	supplierUsecase := usecase.NewSupplierUsecase(supplierRepo)
	purchaseOrderUsecase := usecase.NewPurchaseOrderUsecase(purchaseOrderRepo, stockRepo)
	purchaseInvoiceUsecase := usecase.NewPurchaseInvoiceUsecase(purchaseInvoiceRepo)
	purchasePaymentUsecase := usecase.NewPurchasePaymentUsecase(purchasePaymentRepo, purchaseInvoiceRepo)
	purchaseReturnUsecase := usecase.NewPurchaseReturnUsecase(purchaseReturnRepo, stockRepo)
	customerCategoryUsecase := usecase.NewCustomerCategoryUsecase(customerCategoryRepo)
	salesCategoryUsecase := usecase.NewSalesCategoryUsecase(salesCategoryRepo)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	salesQuotationUsecase := usecase.NewSalesQuotationUsecase(salesQuotationRepo)
	salesOrderUsecase := usecase.NewSalesOrderUsecase(salesOrderRepo)
	deliveryOrderUsecase := usecase.NewDeliveryOrderUsecase(deliveryOrderRepo, salesOrderRepo, stockRepo)
	salesInvoiceUsecase := usecase.NewSalesInvoiceUsecase(salesInvoiceRepo)
	salesReceiptUsecase := usecase.NewSalesReceiptUsecase(salesReceiptRepo, salesInvoiceRepo)
	salesReturnUsecase := usecase.NewSalesReturnUsecase(salesReturnRepo, salesOrderRepo, stockRepo)
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo)
	itemUsecase := usecase.NewItemUsecase(itemRepo, productRepo, stockRepo)
	transferUsecase := usecase.NewStockTransferUsecase(transferRepo, stockRepo)
	adjustmentUsecase := usecase.NewStockAdjustmentUsecase(adjustmentRepo, stockRepo)
	priceChangeUsecase := usecase.NewPriceChangeUsecase(priceChangeRepo)
	stockUsecase := usecase.NewStockUsecase(stockRepo)
	bankTransferUsecase := usecase.NewBankTransferUsecase(bankTransferRepo, cashRepo, stockRepo)
	bankReconciliationUsecase := usecase.NewBankReconciliationUsecase(bankReconciliationRepo)
	chartOfAccountsUsecase := usecase.NewChartOfAccountsUsecase(chartOfAccountsRepo)
	journalEntryUsecase := usecase.NewJournalEntryUsecase(journalEntryRepo)
	fixedAssetUsecase := usecase.NewFixedAssetUsecase(fixedAssetRepo)
	periodEndUsecase := usecase.NewPeriodEndUsecase(periodEndRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, revenueRepo, shiftRepo)
	marketplaceUsecase := usecase.NewMarketplaceUsecase(marketplaceRepo, salesOrderRepo)
	bankStatementUsecase := usecase.NewBankStatementUsecase(bankStatementRepo, bankReconciliationRepo)
	efakturUsecase := usecase.NewEFakturUsecase(efakturRepo)
	manufacturingUsecase := usecase.NewManufacturingUsecase(bomRepo, workOrderRepo, stockRepo, itemRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo)
	tenantUsecase := usecase.NewTenantUsecase(tenantRepo)

	// Handlers (interface adapters / delivery layer)
	return rest.Handlers{
		Auth:               rest.NewAuthHandler(authUsecase, cfg.JWT.SecretKey),
		Category:           rest.NewCategoryHandler(categoryUsecase),
		Product:            rest.NewProductHandler(productUsecase),
		User:               rest.NewUserHandler(userUsecase),
		Upload:             rest.NewUploadHandler(),
		Shift:              rest.NewShiftHandler(shiftUsecase),
		QRIS:               rest.NewQRISHandler(),
		Transaction:        rest.NewTransactionHandler(transactionUsecase, orderUsecase),
		Branch:             rest.NewBranchHandler(branchUsecase),
		Department:         rest.NewDepartmentHandler(departmentUsecase),
		Employee:           rest.NewEmployeeHandler(employeeUsecase),
		Commission:         rest.NewCommissionHandler(commissionUsecase),
		Salary:             rest.NewSalaryHandler(salaryUsecase),
		Cash:               rest.NewCashHandler(cashUsecase),
		Debt:               rest.NewDebtHandler(debtUsecase),
		Receivable:         rest.NewReceivableHandler(receivableUsecase),
		Withdrawal:         rest.NewWithdrawalHandler(withdrawalUsecase),
		Supplier:           rest.NewSupplierHandler(supplierUsecase),
		PurchaseOrder:      rest.NewPurchaseOrderHandler(purchaseOrderUsecase),
		PurchaseInvoice:    rest.NewPurchaseInvoiceHandler(purchaseInvoiceUsecase),
		PurchasePayment:    rest.NewPurchasePaymentHandler(purchasePaymentUsecase),
		PurchaseReturn:     rest.NewPurchaseReturnHandler(purchaseReturnUsecase),
		CustomerCategory:   rest.NewCustomerCategoryHandler(customerCategoryUsecase),
		SalesCategory:      rest.NewSalesCategoryHandler(salesCategoryUsecase),
		Customer:           rest.NewCustomerHandler(customerUsecase),
		SalesQuotation:     rest.NewSalesQuotationHandler(salesQuotationUsecase),
		SalesOrder:         rest.NewSalesOrderHandler(salesOrderUsecase),
		DeliveryOrder:      rest.NewDeliveryOrderHandler(deliveryOrderUsecase),
		SalesInvoice:       rest.NewSalesInvoiceHandler(salesInvoiceUsecase),
		SalesReceipt:       rest.NewSalesReceiptHandler(salesReceiptUsecase),
		SalesReturn:        rest.NewSalesReturnHandler(salesReturnUsecase),
		Warehouse:          rest.NewWarehouseHandler(warehouseUsecase),
		Item:               rest.NewItemHandler(itemUsecase),
		StockTransfer:      rest.NewStockTransferHandler(transferUsecase),
		StockAdjustment:    rest.NewStockAdjustmentHandler(adjustmentUsecase),
		PriceChange:        rest.NewPriceChangeHandler(priceChangeUsecase),
		Stock:              rest.NewStockHandler(stockUsecase),
		BankTransfer:       rest.NewBankTransferHandler(bankTransferUsecase),
		BankReconciliation: rest.NewBankReconciliationHandler(bankReconciliationUsecase),
		ChartOfAccounts:    rest.NewChartOfAccountsHandler(chartOfAccountsUsecase),
		JournalEntry:       rest.NewJournalEntryHandler(journalEntryUsecase),
		FixedAsset:         rest.NewFixedAssetHandler(fixedAssetUsecase),
		PeriodEnd:          rest.NewPeriodEndHandler(periodEndUsecase),
		Order:              rest.NewOrderHandler(orderUsecase),
		Marketplace:        rest.NewMarketplaceHandler(marketplaceUsecase),
		BankStatement:      rest.NewBankStatementHandler(bankStatementUsecase),
		EFaktur:            rest.NewEFakturHandler(efakturUsecase),
		Bom:                rest.NewBomHandler(manufacturingUsecase),
		WorkOrder:          rest.NewWorkOrderHandler(manufacturingUsecase),
		Report:             rest.NewReportHandler(reportUsecase),
		Tenant:             rest.NewTenantHandler(tenantUsecase),
	}
}
