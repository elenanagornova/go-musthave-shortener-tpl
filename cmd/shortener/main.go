package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	listen, exists := os.LookupEnv("SERVER_ADDRESS")
	if !exists {
		log.Fatal("Отсутствует переменная окружения SERVER_ADDRESS")
	}
	if listen == "" {
		log.Fatal("Отсутствует значение переменной окружения SERVER_ADDRESS")
	}

	addr, exists := os.LookupEnv("BASE_URL")
	if !exists {
		log.Fatal("Отсутствует переменная окружения BASE_URL")
	}
	if listen == "" {
		log.Fatal("Отсутствует значение переменной окружения BASE_URL")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	service := shortener.New(addr)

	log.Println("Starting server at port 8080")
	srv := http.Server{
		Addr:    listen,
		Handler: NewRouter(service),
	}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
func NewRouter(service *shortener.Shortener) chi.Router {
	r := chi.NewRouter()
	r.Post("/api/shorten", makeShortenLink(service))
	r.Post("/", makeShortLink(service))
	r.Get("/{shortLink}", getLinkByID(service))
	return r
}
