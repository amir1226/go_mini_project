package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example:
// Authorization: ApiKey {insert API key here}
func GetAPIKey(h http.Header) (string, error) {
	auth_val := h.Get("Authorization")
	if auth_val == "" {
		return "", errors.New("authorization header is missing")
	}

	vals := strings.Split(auth_val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid Authorization type. Should be 'ApiKey'")
	}

	return vals[1], nil
}
