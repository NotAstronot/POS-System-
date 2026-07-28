package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-system/config"
	"pos-system/internal/handler/rest"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := "host=" + cfg.Database.Host +
		" port=" + cfg.Database.Port +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.DBName +
		" sslmode=" + cfg.Database.SSLMode +
		" connect_timeout=10"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	categoryRepo := postgres.NewCategoryRepository(db)
	productRepo := postgres.NewProductRepository(db)
	userRepo := postgres.NewUserRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	shiftRepo := postgres.NewShiftRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	orderItemRepo := postgres.NewOrderItemRepository(db)

	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	shiftUsecase := usecase.NewShiftUsecase(shiftRepo, orderRepo, orderItemRepo, productRepo)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, orderRepo)

	authHandler := rest.NewAuthHandler(authUsecase, cfg.JWT.SecretKey)
	categoryHandler := rest.NewCategoryHandler(categoryUsecase)
	productHandler := rest.NewProductHandler(productUsecase)
	shiftHandler := rest.NewShiftHandler(shiftUsecase)
	transactionHandler := rest.NewTransactionHandler(transactionUsecase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authMw := rest.JWTAuth(cfg.JWT.SecretKey)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		pos := api.Group("")
		pos.Use(authMw)
		{
			pos.GET("/products/search", productHandler.Search)
			pos.GET("/products/barcode", productHandler.Barcode)
			pos.GET("/products/category/:categoryId", productHandler.ListByCategory)
			pos.GET("/categories", categoryHandler.List)
			pos.GET("/availability", productHandler.ListAvailability)
			pos.GET("/availability/:productId", productHandler.GetAvailability)
			pos.POST("/availability", productHandler.SetAvailability)

			pos.GET("/products", productHandler.List)
			pos.POST("/products", productHandler.Create)
			pos.GET("/products/:id", productHandler.GetByID)
			pos.PUT("/products/:id", productHandler.Update)
			pos.DELETE("/products/:id", productHandler.Delete)
			pos.POST("/categories", categoryHandler.Create)
			pos.GET("/categories/:id", categoryHandler.GetByID)
			pos.PUT("/categories/:id", categoryHandler.Update)
			pos.DELETE("/categories/:id", categoryHandler.Delete)
		}

		shift := api.Group("/shifts")
		shift.Use(authMw)
		{
			shift.GET("/active", shiftHandler.GetActiveShift)
		}

		transaction := api.Group("/transactions")
		transaction.Use(authMw)
		{
			transaction.GET("/", transactionHandler.List)
			transaction.POST("/", transactionHandler.Create)
			transaction.GET("/:id", transactionHandler.GetByID)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}