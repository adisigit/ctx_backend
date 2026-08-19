package main

import (
	"context"
	"ctx_backend/internal/server"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	log.Println("Server exiting")
	done <- true
}

func main() {
	server := server.NewServer()
	done := make(chan bool, 1)
	go gracefulShutdown(server, done)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
	<-done
	log.Println("Graceful shutdown complete.")
}
