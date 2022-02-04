package controller

import (
	"encoding/json"
	"go-musthave-shortener-tpl/internal/entity"
	"go-musthave-shortener-tpl/internal/hellpers"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
)

func MakeShortLinkBatch(service *shortener.Shortener) http.HandlerFunc {
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

		var links []entity.BatchShortenerRequest
		if err := json.NewDecoder(r.Body).Decode(&links); err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}

		var dbLinks []entity.DBBatchShortenerLinks
		for _, link := range links {
			shortURL := hellpers.GenerateRandomString(5)
			dbLinks = append(dbLinks, entity.DBBatchShortenerLinks{ShortURL: shortURL, OriginalURL: link.OriginalURL, UserUID: userUID, CorrelationID: link.CorrelationID})
		}

		allLinks, err := service.Repo.BatchSaveLinks(dbLinks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var respBody []entity.BatchShortenerResponse
		for _, resLink := range allLinks {
			respBody = append(respBody, entity.BatchShortenerResponse{CorrelationID: resLink.CorrelationID, ShortURL: service.Addr + resLink.ShortURL})
		}
		//resultLink, err := service.GenerateShortLink(originalLink.URL, userUID)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		//var responseBody ShortenerResponse
		//responseBody.Result = resultLink
		hellpers.SetUIDCookie(w, userUID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(respBody); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}
