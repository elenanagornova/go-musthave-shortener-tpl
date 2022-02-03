package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-musthave-shortener-tpl/internal/controller"
	"go-musthave-shortener-tpl/internal/repository"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	listen          string
	Addr            string
	fileStoragePath string
	DatabaseDSN     string
)

func main() {
	if listen = os.Getenv("SERVER_ADDRESS"); listen == "" {
		flag.StringVar(&listen, "a", ":8080", "Server address")
	}

	if Addr = os.Getenv("BASE_URL"); Addr == "" {
		flag.StringVar(&Addr, "b", "http://localhost:8080/", "Server base URL")
	}

	if fileStoragePath = os.Getenv("FILE_STORAGE_PATH"); fileStoragePath == "" {
		flag.StringVar(&fileStoragePath, "f", "links.log", "File storage path")
	}

	if DatabaseDSN = os.Getenv("DATABASE_DSN"); DatabaseDSN == "" {
		flag.StringVar(&DatabaseDSN, "d", "postgres://shorteneruser:pgpwd4@localhost:5432/shortenerdb", "")
	}
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	conn, err := repository.CreateDBConnect(DatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	service := shortener.New(Addr, conn)
	err = service.Restore(fileStoragePath)
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
	r.Use(middleware.Logger)

	r.Use(middleware.Compress(5))
	r.Use(controller.GzipDecompressor)
	r.Use(controller.UserMiddleware)

	r.Post("/api/shorten", controller.MakeShortLinkJSON(service))
	r.Post("/", controller.MakeShortLink(service))
	r.Get("/{shortLink}", controller.GetLinkByID(service))
	r.Get("/user/urls", controller.GetUserLinks(service))
	r.Get("/ping", controller.CheckConn(service))
	return r
}
