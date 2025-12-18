package token

import (
	"os"

	"github.com/jgamaraalv/movies.git/pkg/logger"
)

func GetJWTSecret(log logger.Logger) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-dev"
		log.Info("JWT_SECRET not set, using default development secret")
	} else {
		log.Info("Using JWT_SECRET from environment")
	}
	return jwtSecret
}
