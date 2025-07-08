package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

var AppEnv string

func main() {
	svcConfig := &service.Config{
		Name:        "VaultExporterService",
		DisplayName: "Vault Exporter Service",
		Description: "A tool for exporting data from Vault to KS",
	}

	prg := &program{isProd: AppEnv == "production"}
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
	}
}
