package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	listen          string
	addr            string
	fileStoragePath string
)

func main() {
	if listen = os.Getenv("SERVER_ADDRESS"); listen == "" {
		flag.StringVar(&listen, "a", ":8080", "Server address")
	}

	if addr = os.Getenv("BASE_URL"); addr == "" {
		flag.StringVar(&addr, "b", "http://localhost:8080/", "Server base URL")
	}

	if fileStoragePath = os.Getenv("FILE_STORAGE_PATH"); fileStoragePath == "" {
		flag.StringVar(&fileStoragePath, "f", "links.log", "File storage path")
	}

	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	service := shortener.New(addr)
	err := service.Restore(fileStoragePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("Can't restore data from file: %s", err.Error()))
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
