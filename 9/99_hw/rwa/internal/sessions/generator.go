package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const secret = "SecretKey"

// Создает новую сессию, возвращая
func GenerateSession() (public string, err error) {
	randomBytes := make([]byte, 32)
	if _, err = rand.Read(randomBytes); err != nil {
		return "", err
	}

	// публичный ключ
	public = base64.URLEncoding.EncodeToString([]byte(string(randomBytes)))

	return public, nil
}

// Получить приватный ключ на основе публичного токена
func GetPrivateKey(publicKey string) string {
	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(publicKey))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
