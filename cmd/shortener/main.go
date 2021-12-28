package main

import (
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
)

const addr string = "localhost:8080"

func main() {
	service := shortener.New(addr)

	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(addr, NewRouter(service)))
}
func NewRouter(service *shortener.Shortener) chi.Router {
	r := chi.NewRouter()
	r.Post("/", makeShortLink(service))
	r.Get("/{shortLink}", getLinkByID(service))
	return r
}
