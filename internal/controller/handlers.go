package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go-musthave-shortener-tpl/internal/entity"
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

func MakeShortLink(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := hellpers.GetUID(r.Cookies())

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
		resultLink, err := service.GenerateShortLink(string(body), userUID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				hellpers.SetUIDCookie(w, userUID)
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(resultLink))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
		hellpers.SetUIDCookie(w, userUID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resultLink))
	}
}

func GetLinkByID(service *shortener.Shortener) http.HandlerFunc {
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

func MakeShortLinkJSON(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := userUIDFromRequest(r)
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

		resultLink, err := service.GenerateShortLink(originalLink.URL, userUID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(resultLink))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var responseBody ShortenerResponse
		responseBody.Result = resultLink

		hellpers.SetUIDCookie(w, userUID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}

func GetUserLinks(service *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := userUIDFromRequest(r)
		links := service.GetLinks(userUID)
		var fullLinks []entity.UserLinks
		for _, link := range links {
			fullLinks = append(fullLinks, entity.UserLinks{ShortURL: service.Addr + link.ShortURL, OriginalURL: link.OriginalURL})
		}
		if len(links) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		hellpers.SetUIDCookie(w, userUID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fullLinks); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}
