package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*

It should extract the api key from the Authorization header, which is expected to be in this format:

Authorization: ApiKey THE_KEY_HERE

You'll need to strip out the ApiKey part and the whitespace and return just the key.

*/

func GetApiKey(headers http.Header) (string, error) {
	headerValue := headers.Get("Authorization")
	if headerValue == "" {
		return "", errors.New("no Authorization header found")
	}

	if !strings.HasPrefix(headerValue, "ApiKey ") {
		return "", errors.New("invalid Authorization value format")
	}

	apiKey := strings.TrimPrefix(headerValue, "ApiKey ")
	if apiKey == "" {
		return "", errors.New("no api key found")
	}

	return strings.TrimSpace(apiKey), nil
}
