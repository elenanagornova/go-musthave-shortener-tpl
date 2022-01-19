package main

import (
	"context"
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	listen          *string
	addr            *string
	fileStoragePath *string
)

func init() {
	if listenFromEnv := os.Getenv("SERVER_ADDRESS"); listenFromEnv != "" {
		listen = &listenFromEnv
	} else {
		listen = flag.String("a", ":8080", "Server address")
	}

	if addrFromEnv := os.Getenv("BASE_URL"); addrFromEnv != "" {
		addr = &addrFromEnv
	} else {
		addr = flag.String("b", "http://localhost:8080/", "Server base URL")
	}

	if fileStoragePathFromEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathFromEnv != "" {
		fileStoragePath = &fileStoragePathFromEnv
	} else {
		fileStoragePath = flag.String("f", "links.log", "File storage path")
	}

}

func main() {
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	service := shortener.New(*addr)
	err := service.Restore(*fileStoragePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("Can't restore data from file: %s", err.Error())
	}
	log.Println("Starting server at port 8080")
	srv := http.Server{
		Addr:    *listen,
		Handler: NewRouter(service),
	}
	go func() {
		<-ctx.Done()
		err := service.Save(*fileStoragePath)
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
