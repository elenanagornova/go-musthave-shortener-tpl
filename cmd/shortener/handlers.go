package main

import (
	"encoding/json"
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

		resultLink, err := service.GenerateShortLink(string(body))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Location", originalLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

type ShortenerRequest struct {
	URL string `json:"url"`
}
type ShortenerResponse struct {
	Result string `json:"result"`
}

func makeShortenLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentTtype := r.Header.Get("Content-Type")
		if headerContentTtype != "application/json" {
			http.Error(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
			return
		}

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

		var originalLink ShortenerRequest
		err = json.Unmarshal(body, &originalLink)
		if err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}

		resultLink, err := service.GenerateShortLink(originalLink.URL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var responseBody ShortenerResponse
		responseBody.Result = resultLink
		respBody, err := json.Marshal(&responseBody)
		if err != nil {
			http.Error(w, "Marshalling error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(respBody)
	}
}
