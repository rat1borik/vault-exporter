// Package config описывает конфигурацию приложения.
package config

import (
	"log"
	"os"
	"vault-exporter/internal/utils"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Server struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		ApiKey   string `yaml:"api_key"`
		TLS      bool   `yaml:"tls"`
		CertPath string `yaml:"cert_path"`
		KeyPath  string `yaml:"key_path"`
	} `yaml:"server"`
	Vault struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"vault"`
	KSFilesPath string `yaml:"ks_files_path"`
	TempPath    string `yaml:"temp_path"`
	KSDatabase  struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
	} `yaml:"ks_database"`
	// Берется из env
	IsProduction bool
}

func LoadConfig(filename string) (*ServerConfig, error) {
	path, err := utils.ExecPath(filename)
	if err != nil {
		log.Fatalf("can't load config: %v", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
