package repository

import (
	"go-musthave-shortener-tpl/internal/config"
)

func NewRepository(cfg *config.ShortenerConfiguration) (Storager, error) {
	switch {
	case cfg.DatabaseDSN != "":
		return NewDBConnect(cfg.DatabaseDSN)
	case cfg.FileStoragePath != "":
		return NewFileConnect(cfg.FileStoragePath), nil
	default:
		return NewInMemoryConnect(), nil
	}
}
