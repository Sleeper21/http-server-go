package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	headerValue := headers.Get("Authorization")
	if headerValue == "" {
		return "", errors.New("no authorization header found")
	}

	// the header value will be something like: Bearer TOKEN_STRING
	// stripping off the Bearer prefix and whitespace
	if !strings.HasPrefix(headerValue, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	bearerToken := strings.TrimPrefix(headerValue, "Bearer ")
	if bearerToken == "" {
		return "", errors.New("no token found in authorization header")
	}

	return bearerToken, nil
}
