package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Encrypt(key []byte, text []byte) (string, []byte, error) {

	aesgcm, err := initEncryption(key)
	if err != nil {
		return "", nil, err
	}
	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		return "", nil, fmt.Errorf("encrypt: %w", err)
	}
	encrypted := aesgcm.Seal(nil, nonce, text, nil)

	return base64.StdEncoding.EncodeToString(encrypted), nonce, nil
}

func initEncryption(key []byte) (cipher.AEAD, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	return aesgcm, nil
}

func Decrypt(key []byte, nonce []byte, encryptedText string) ([]byte, error) {
	aesgcm, err := initEncryption(key)
	if err != nil {
		return nil, err
	}
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, encrypted, nil) // расшифровываем
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
