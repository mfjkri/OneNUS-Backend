package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(username string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 8760).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func DecodeJWT(tokenString string) (username string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = claims["sub"].(string)
		return
	}

	username = ""
	return
}

func ValidateJWT(tokenString string, username string) (jwtValid bool, err error) {
	decoded_username, err := DecodeJWT(tokenString)

	if decoded_username == username {
		jwtValid = true
		return
	}

	jwtValid = false
	return
}
