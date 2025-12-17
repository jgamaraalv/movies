package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

func CreateJWT(user models.User, log logger.Logger) string {
	jwtSecret := GetJWTSecret(log)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Error("Failed to sign JWT", err)
		return ""
	}

	return tokenString
}
