package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-musthave-shortener-tpl/internal/config"
	"go-musthave-shortener-tpl/internal/controller"
	"go-musthave-shortener-tpl/internal/repository"
	"go-musthave-shortener-tpl/internal/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.LoadConfiguration()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	rep, err := repository.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("Can't create repository: %s", err.Error()))
	}
	service := shortener.New(cfg.Addr, rep)
	log.Println("Starting server at port 8080")

	srv := http.Server{
		Addr:    cfg.Listen,
		Handler: NewRouter(service),
	}
	go func() {
		<-ctx.Done()
		err := service.Repo.FinalSave()
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
	r.Post("/api/shorten/batch", controller.MakeShortLinkBatch(service))
	return r
}
