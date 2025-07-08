package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vault-exporter/internal/router"
)

func main() {
	r := router.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go runServer(srv)

	// Создаем канал для перехвата сигналов OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Попытка корректно завершить работу сервера
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func runServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
