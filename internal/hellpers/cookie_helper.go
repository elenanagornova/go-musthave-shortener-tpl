package hellpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const CookieName = "UID"

var uidLen = 5
var secret = "mysecret" // Прочитать из env/конфига

func CalculateHash(uid string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(uid))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkHash(uid string, hash string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(uid))
	sign, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return hmac.Equal(sign, h.Sum(nil))
}

func GetUID(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		// пробуем получить значение uid и hash из куки
		if cookie.Name == CookieName {
			parts := strings.Split(cookie.Value, ":")
			if len(parts) != 2 {
				// если в куки нет обоих параметров, то генерируем новый uid
				return GenerateRandomString(uidLen)
			}
			uid, hash := parts[0], parts[1]
			if checkHash(uid, hash) {
				return uid
			}
		}
	}
	return GenerateRandomString(uidLen)
}

// SetUIDCookie сохраняет в куку uid пользователя вместе с его hmac
func SetUIDCookie(w http.ResponseWriter, uid string) {
	uuidSigned := fmt.Sprintf("%s:%s", uid, CalculateHash(uid))

	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  uuidSigned,
		MaxAge: 3000000,
	})
}
