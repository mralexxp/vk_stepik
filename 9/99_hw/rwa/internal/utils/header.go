package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func GetHeaderToken(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	if !strings.HasPrefix(authHeader, "Token ") {
		return "", fmt.Errorf("invalid authorization header")
	}

	token = strings.TrimPrefix(authHeader, "Token ")
	if token == "" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return token, nil
}
