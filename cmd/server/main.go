// Package main запускает HTTP сервер приложения.
package main

import (
	"log"
	"os"
	"vault-exporter/internal/logger"

	"github.com/kardianos/service"
)

var AppEnv string

func main() {
	logger := logger.NewLogrusLogger()

	// Установка вывода стандартного log в logrus
	log.SetOutput(logger.Writer())
	log.SetFlags(0) // убираем timestamp, так как logrus добавит свой

	svcConfig := &service.Config{
		Name:        "VaultExporterService",
		DisplayName: "Vault Exporter Service",
		Description: "A tool for exporting data from Vault to KS",
	}

	prg := &program{isProd: AppEnv == "production", logger: logger}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Если запустили с параметром install/uninstall/start/stop
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Запуск сервиса
	err = s.Run()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
