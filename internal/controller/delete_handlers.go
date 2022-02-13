package controller

import (
	"encoding/json"
	"go-musthave-shortener-tpl/internal/deleter"
	"go-musthave-shortener-tpl/internal/hellpers"
	"log"
	"net/http"
)

func DeleteLinks(deleteCh chan deleter.DeleteTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := userUIDFromRequest(r)

		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()

		var userLinks []string
		if err := json.NewDecoder(r.Body).Decode(&userLinks); err != nil {
			http.Error(w, "Something wrong with request", http.StatusBadRequest)
			return
		}
		deleteCh <- deleter.DeleteTask{
			UID:   "",
			Links: nil,
		}
		hellpers.SetUIDCookie(w, userUID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
	}
}
