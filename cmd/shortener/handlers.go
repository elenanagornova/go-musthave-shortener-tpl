package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/controller"
	"go-musthave-shortener-tpl/internal/hellpers"
	"go-musthave-shortener-tpl/internal/shortener"
	"io"
	"log"
	"net/http"
)

type ShortenerRequest struct {
	URL string `json:"url"`
}
type ShortenerResponse struct {
	Result string `json:"result"`
}

func makeShortLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			URL string `json:"url"`
		}{}
		userUId := controller.UserUIDFromRequest(r)
		hellpers.SetUIDCookie(w, userUId)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Something wrong with request", http.StatusBadRequest)
				return
			}
			req.URL = string(body)
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()

		resultLink, err := service.GenerateShortLink(req.URL, userUId)
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

		userUId := controller.UserUIDFromRequest(r)
		hellpers.SetUIDCookie(w, userUId)

		shortLink := chi.URLParam(r, "shortLink")
		if len(shortLink) == 0 {
			http.Error(w, "Empty required param shortLink", http.StatusBadRequest)
		}

		originalLink, err := service.GetLink(shortLink, userUId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Location", originalLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func makeShortenLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUId := controller.UserUIDFromRequest(r)
		hellpers.SetUIDCookie(w, userUId)
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

		var originalLink ShortenerRequest
		if err := json.NewDecoder(r.Body).Decode(&originalLink); err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}

		resultLink, err := service.GenerateShortLink(originalLink.URL, userUId)
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

func getUserLinks(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := controller.UserUIDFromRequest(r)
		links := service.GetLinks(userId)
		if err := json.NewEncoder(w).Encode(links); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}
