package main

import (
	"log"
	"os"
	"path/filepath"
	"vault-exporter/internal/config"

	"gopkg.in/yaml.v3"
)

var ExecDir string

func init() {
	ExecPath, err := os.Executable()
	if err != nil {
		log.Fatalln(err.Error())
	}

	ExecDir, _ = filepath.Split(ExecPath)

}

func LoadConfig(filename string) (*config.ServerConfig, error) {
	data, err := os.ReadFile(filepath.Join(ExecDir, filename))
	if err != nil {
		return nil, err
	}

	var cfg config.ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
