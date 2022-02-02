package main

import (
	"context"
	"go-musthave-shortener-tpl/internal/crypto"
	"go-musthave-shortener-tpl/internal/shortener"
	"net/http"
	"time"
)

type CtxKey string

const UserCtxKey = CtxKey("UserID")

func (s *Handlers) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		userIDCookie, err := r.Cookie("userID")
		if err != nil {
			if err != http.ErrNoCookie {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Setting cookie
			userID = shortener.GenerateRandomString(5)
			sign, nonce, err := crypto.Encrypt(s.key, []byte(userID))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			s.nounces[userID] = nonce
			expiration := time.Now().Add(365 * 24 * time.Hour)
			uidCookie := http.Cookie{Name: "userID", Value: userID, Expires: expiration}
			signCookie := http.Cookie{Name: "sign", Value: sign, Expires: expiration}
			http.SetCookie(w, &uidCookie)
			http.SetCookie(w, &signCookie)
			next.ServeHTTP(w, userIDtoRequest(r, userID))
			return
		}
		// Validating sign
		userID = userIDCookie.Value
		signCookie, err := r.Cookie("sign")
		if signCookie == nil || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sign := signCookie.Value
		nounce, hasNounce := s.nounces[userID]
		if !hasNounce {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		testUserID, err := crypto.Decrypt(s.key, nounce, sign)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userID != string(testUserID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, userIDtoRequest(r, userID))
		return
	})
}

func userIDtoRequest(r *http.Request, userID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), UserCtxKey, userID))
}

func userIDFromRequest(r *http.Request) string {
	uid := r.Context().Value(UserCtxKey)
	if uid == nil {
		return ""
	}
	if userID, ok := uid.(string); ok {
		return userID
	}
	return ""
}
