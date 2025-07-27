package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKeyHeader := headers.Get("Authorization")
	if apiKeyHeader == "" {
		return "", errors.New("No authorization header included")
	}

	substrings := strings.Fields(apiKeyHeader)
	if len(substrings) != 2 || substrings[0] != "ApiKey" {
		return "", errors.New("malformed ApiKey")
	}

	return substrings[1], nil
}
