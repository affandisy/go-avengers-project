package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set")
		secret = "default-secret"
	}
	return []byte(secret)
}

type JWTClaim struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, role string) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
