package config

import (
	"flag"
	"os"
)

type ShortenerConfiguration struct {
	Listen          string
	Addr            string
	FileStoragePath string
	DatabaseDSN     string
}

func LoadConfiguration() *ShortenerConfiguration {
	cfg := &ShortenerConfiguration{}
	if cfg.Listen = os.Getenv("SERVER_ADDRESS"); cfg.Listen == "" {
		flag.StringVar(&cfg.Listen, "a", ":8080", "Server address")
	}

	if cfg.Addr = os.Getenv("BASE_URL"); cfg.Addr == "" {
		flag.StringVar(&cfg.Addr, "b", "http://localhost:8080/", "Server base URL")
	}

	if cfg.FileStoragePath = os.Getenv("FILE_STORAGE_PATH"); cfg.FileStoragePath == "" {
		flag.StringVar(&cfg.FileStoragePath, "f", "links.log", "File storage path")
	}

	if cfg.DatabaseDSN = os.Getenv("DATABASE_DSN"); cfg.DatabaseDSN == "" {
		flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://shorteneruser:pgpwd4@localhost:5432/shortenerdb?sslmode=disable", "")
	}
	flag.Parse()

	return cfg
}
