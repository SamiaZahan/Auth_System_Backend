package service

import (
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type ExtractJwtClaims struct{}

func (e *ExtractJwtClaims) ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := config.Params.JWTSignatureSecret
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
