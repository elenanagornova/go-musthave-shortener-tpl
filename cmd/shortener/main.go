package main

import (
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-musthave-shortener-tpl/internal/shortener"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	r.Use(middleware.Logger)
	r.Use(MiddlewareFunc)

	r.Post("/api/shorten", makeShortenLink(service))
	r.Post("/", makeShortLink(service))
	r.Get("/{shortLink}", getLinkByID(service))
	return r
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func MiddlewareFunc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func LengthHandle(w http.ResponseWriter, r *http.Request) {
	// переменная reader будет равна r.Body или *gzip.Reader
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Length: %d", len(body))
}
