package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-system/internal/app"
)

func main() {
	application, err := app.New("config.yaml")
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Close()

	srv := &http.Server{
		Addr:    ":" + application.Config.Server.Port,
		Handler: application.Router,
	}

	go func() {
		log.Printf("server starting on port %s", application.Config.Server.Port)
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
