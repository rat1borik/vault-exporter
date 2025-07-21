package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("authorization header does not start with Bearer")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	return token, nil
}
