package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// TODO mutex
var LinksMap = make(map[string]string)

const addr string = "localhost:8080"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", ShortenerHandler)
	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(addr, nil))
}
func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetLinkByID(w, r)
	case "POST":
		MakeShortLink(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}

	// TODO default
}

func MakeShortLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Something wrong with request", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(GenerateShortLink(string(body))))
}

func GetLinkByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "The path is missing", http.StatusBadRequest)
		return
	}

	originalLink, ok := LinksMap[path]
	if !ok {
		http.Error(w, "Link not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func GenerateShortLink(url string) string {
	id := GenerateRandomString(5)
	LinksMap[id] = url
	return "http://" + addr + "/" + id
}
