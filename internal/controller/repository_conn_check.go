package controller

import (
	"go-musthave-shortener-tpl/internal/shortener"
	"net/http"
)

func CheckConn(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.Repo.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
