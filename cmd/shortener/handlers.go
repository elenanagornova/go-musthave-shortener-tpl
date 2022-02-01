package main

import (
	"crypto/aes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"math/rand"
	"net/http"
)

type ShortenerRequest struct {
	URL string `json:"url"`
}
type ShortenerResponse struct {
	Result string `json:"result"`
}

type Handlers struct {
	key     []byte
	nounces map[string][]byte
}

var userId = ""

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *Handlers) makeShortLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			URL string `json:"url"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		userId := userIDFromRequest(r)
		resultLink, err := service.GenerateShortLink(req.URL, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resultLink))
	}
}

func (h *Handlers) getLinkByID(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortLink := chi.URLParam(r, "shortLink")
		if len(shortLink) == 0 {
			http.Error(w, "Empty required param shortLink", http.StatusBadRequest)
		}
		userId := userIDFromRequest(r)
		originalLink, err := service.GetLink(shortLink, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Location", originalLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func NewHandlers() *Handlers {
	key, err := generateRandom(aes.BlockSize) // ключ шифрования
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}

	return &Handlers{
		key:     key,
		nounces: map[string][]byte{},
	}
}

func (h *Handlers) makeShortenLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentTtype := r.Header.Get("Content-Type")
		if headerContentTtype != "application/json" {
			http.Error(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		userId := userIDFromRequest(r)
		var originalLink ShortenerRequest
		if err := json.NewDecoder(r.Body).Decode(&originalLink); err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}

		resultLink, err := service.GenerateShortLink(originalLink.URL, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var responseBody ShortenerResponse
		responseBody.Result = resultLink
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}

func (h *Handlers) getUserLinks(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := userIDFromRequest(r)
		links := service.GetLinks(userId)
		if err := json.NewEncoder(w).Encode(links); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}
