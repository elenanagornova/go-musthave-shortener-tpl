package main

import (
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"io"
	"log"
	"net/http"
)

func makeShortLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		if len(body) == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		resultLink := service.GenerateShortLink(string(body))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resultLink))
	}
}

func getLinkByID(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortLink := chi.URLParam(r, "shortLink")
		if len(shortLink) == 0 {
			http.Error(w, "Empty required param shortLink", http.StatusBadRequest)
		}

		originalLink, err := service.GetLink(shortLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
