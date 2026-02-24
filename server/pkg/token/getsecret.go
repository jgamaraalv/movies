package token

import (
	"os"

	"github.com/jgamaraalv/movies.git/pkg/logger"
)

func GetJWTSecret(log logger.Logger) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Error("JWT_SECRET environment variable is not set", nil)
		os.Exit(1)
	}
	return jwtSecret
}
