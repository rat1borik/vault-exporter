package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"vault-exporter/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

type program struct {
	server *http.Server
	isProd bool
}

func (p *program) Start(s service.Service) error {
	if p.isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Starting in development mode")
	}
	r := router.SetupRouter()

	p.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Запуск в отдельной горутине
	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (p *program) Stop(s service.Service) error {
	log.Println("Stopping HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}
