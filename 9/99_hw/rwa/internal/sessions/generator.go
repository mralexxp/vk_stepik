package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const secret = "SecretKey"

func GenerateSession(username string) (public string, private string, err error) {
	randomBytes := make([]byte, 32)
	if _, err = rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	// публичный ключ
	public = base64.URLEncoding.EncodeToString(randomBytes) + username + time.Now().String()

	// приватный ключ
	// TODO: При проверке сессии работаем также
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(public))
	private = base64.URLEncoding.EncodeToString(h.Sum(nil))

	return public, private, nil
}

func GetPrivateKey(publicKey string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(publicKey))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
