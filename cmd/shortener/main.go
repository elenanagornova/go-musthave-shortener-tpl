package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var LinksMap = make(map[string]string)

func main() {
	http.HandleFunc("/", ShortenerHandler)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetLinkByID(w, r)
	case "POST":
		SetShortLink(w, r)
	}
}

func SetShortLink(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	defer request.Body.Close()

	if err != nil {
		http.Error(writer, "Something wrong with request", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(writer, "Request body is empty", http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(GenerateShortLink(string(body))))
}

func GetLinkByID(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimLeft(request.URL.Path, "/")
	if path == "" {
		http.Error(writer, "The path is missing", http.StatusBadRequest)
		return
	}

	originalLink, ok := LinksMap[path]
	if !ok {
		http.Error(writer, "Link not found", http.StatusBadRequest)
		return
	}

	writer.Header().Add("Location", originalLink)
	writer.WriteHeader(http.StatusTemporaryRedirect)
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
	newURL := GenerateRandomString(5)
	LinksMap[newURL] = url
	return newURL
}
