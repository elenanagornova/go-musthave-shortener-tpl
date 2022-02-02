package hellpers

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)
	for i := 0; i < n; i++ {
		result += string(letters[rnd.Intn(len(letters))])
	}
	return result
}
