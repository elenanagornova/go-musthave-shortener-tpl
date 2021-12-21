package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var LinksMap map[string]string

func main() {
	LinksMap = make(map[string]string)

	http.HandleFunc("/", multiplexer)
	err := http.ListenAndServe(":8080", nil)
	check(err)

}
func multiplexer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetLinkById(w, r)
	case "POST":
		GetShortLink(w, r)
	}
}

func GetShortLink(writer http.ResponseWriter, request *http.Request) {

	body, err := io.ReadAll(request.Body)
	check(err)
	if body == nil {
		http.Error(writer, "Request body is empty", http.StatusBadRequest)
		return
	}

	var t urlsStruct
	err = json.Unmarshal(body, &t)
	check(err)

	writer.WriteHeader(201)
	fmt.Fprintf(writer, GetTransformUrl(t.Url))
}

func GetLinkById(writer http.ResponseWriter, request *http.Request) {
	q := request.URL.Query().Get("id")
	if q == "" {
		http.Error(writer, "The id parameter is missing", http.StatusBadRequest)
		return
	}
	if isValueExists(LinksMap, q) != true {
		http.Error(writer, "Don't found value", http.StatusBadRequest)
		return
	}
	writer.Header().Add("Location", LinksMap[q])
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func isValueExists(links map[string]string, shortUrl string) bool {
	_, ok := links[shortUrl]
	return ok
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func GetTransformUrl(url string) string {
	newUrl := RandomString(5)
	LinksMap[newUrl] = url
	return newUrl
}

type urlsStruct struct {
	Url string
}
