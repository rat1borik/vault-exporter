package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"vault-exporter/internal/config"
	"vault-exporter/internal/infrastructure"
	"vault-exporter/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	svc "github.com/kardianos/service"
)

type program struct {
	server *http.Server
	isProd bool
	db     *pgxpool.Pool
}

func (p *program) Start(s svc.Service) error {

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return err
	}

	if pool, err := runDb(cfg); err != nil {
		log.Fatal(err)
		return err
	} else {
		p.db = pool
	}

	if p.isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Starting in development mode")

	}
	r := router.SetupServer(cfg, p.db)

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

	// Сначала корректно останавливаем HTTP сервер
	err := p.server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	p.db.Close()

	return err
}

func runDb(cfg *config.ServerConfig) (*pgxpool.Pool, error) {
	pgcfg, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.KSDatabase.User,
		cfg.KSDatabase.Password,
		cfg.KSDatabase.Host,
		cfg.KSDatabase.Port,
		cfg.KSDatabase.Name))
	if err != nil {
		return nil, err
	}

	pgcfg.ConnConfig.Tracer = infrastructure.PgTracer{IsProd: false}
	pgcfg.MaxConnLifetime = time.Hour
	pgcfg.MaxConns = 10
	pgcfg.MinIdleConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), pgcfg)
	if err != nil {
		return nil, err
	}

	return pool, err
}
