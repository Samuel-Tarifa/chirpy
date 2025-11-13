package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token_string := headers.Get("Authorization")
	if token_string == "" {
		return "", fmt.Errorf("no Authorization provided")
	}
	token_string, _ = strings.CutPrefix(token_string, "Bearer ")

	return token_string, nil

}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no apiKey provided")
	}
	apiKey, _ = strings.CutPrefix(apiKey, "ApiKey ")
	return apiKey, nil
}
