package controller

import (
	"context"
	"go-musthave-shortener-tpl/internal/hellpers"
	"net/http"
)

type CtxKey string

const UserCtxKey = CtxKey("UserID")

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUID := hellpers.GetUID(r.Cookies())
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCtxKey, userUID)))
	})
}

func userUIDFromRequest(r *http.Request) string {
	uid := r.Context().Value(UserCtxKey)
	if uid == nil {
		return ""
	}
	if userID, ok := uid.(string); ok {
		return userID
	}
	return ""
}
