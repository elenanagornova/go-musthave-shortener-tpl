package controller

import (
	"context"
	"go-musthave-shortener-tpl/internal/shortener"
	"net/http"
)

func CheckConn(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// обратиться к бд
		if err := service.DBConn.Ping(context.Background()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
