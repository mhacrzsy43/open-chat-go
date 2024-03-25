package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // Token有效期1小时
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// ValidateToken validates the jwt token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

// GetUserIDFromToken extracts userID from the jwt token if valid
func GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err // Return 0 for userID in case of error (0 is not a valid userID)
	}

	return claims.UserID, nil
}
