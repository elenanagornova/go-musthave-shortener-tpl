package shortener

import (
	"fmt"
	"math/rand"
	"time"
)

type Shortener struct {
	linksMap map[string]string
	addr     string
}

func New(addr string) *Shortener {
	return &Shortener{
		addr:     addr,
		linksMap: make(map[string]string),
	}
}

func (s *Shortener) generateRandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)
	for i := 0; i < n; i++ {
		result += string(letters[rnd.Intn(len(letters))])
	}
	return result
}

func (s *Shortener) GenerateShortLink(url string) string {
	id := s.generateRandomString(5)
	s.linksMap[id] = url
	return fmt.Sprintf("http://%s/%s", s.addr, id)
}

func (s *Shortener) GetLink(url string) (string, error) {
	link, ok := s.linksMap[url]
	if !ok {
		return "", fmt.Errorf("link not found")
	}
	return link, nil
}
