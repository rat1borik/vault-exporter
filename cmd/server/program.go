package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"vault-exporter/internal/config"
	"vault-exporter/internal/router"

	"github.com/gin-gonic/gin"
	svc "github.com/kardianos/service"
)

type program struct {
	server *http.Server
	isProd bool
}

func (p *program) Start(s svc.Service) error {

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return err
	}

	// Запускаем коннект с базой
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.KSDatabase.User,
		cfg.KSDatabase.Password,
		cfg.KSDatabase.Host,
		cfg.KSDatabase.Port,
		cfg.KSDatabase.Name))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Определяем environment
	if p.isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Starting in development mode")

	}
	r := router.SetupServer(cfg, db)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	p.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("HTTP server on %s", addr)

	// Запуск в отдельной горутине
	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (p *program) Stop(s svc.Service) error {
	log.Println("Stopping HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}
