package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/helpers"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	listen          = helpers.GetEnvOrDefault("SERVER_ADDRESS", ":8080")
	addr            = helpers.GetEnvOrDefault("BASE_URL", "http://localhost:8080/")
	fileStoragePath = helpers.GetEnvOrDefault("FILE_STORAGE_PATH", "links.log")
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	service := shortener.New(addr)
	err := service.Restore(fileStoragePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("Can't restore data from file: %s", err.Error())
	}
	log.Println("Starting server at port 8080")
	srv := http.Server{
		Addr:    listen,
		Handler: NewRouter(service),
	}
	go func() {
		<-ctx.Done()
		err := service.Save(fileStoragePath)
		if err != nil {
			log.Printf("Can't save data in file")
		}
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
