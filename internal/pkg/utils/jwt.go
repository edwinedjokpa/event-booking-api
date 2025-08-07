package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// func GenerateToken(userID string, email string, jwtSecret string) (string, error) {
// 	expirationTime := time.Now().Add(1 * time.Hour)
// 	claims := jwt.MapClaims{
// 		"userID": userID,
// 		"email":  email,
// 		"iat":    time.Now().UTC().Unix(),
// 		"exp":    expirationTime.Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	signedToken, err := token.SignedString([]byte(jwtSecret))
// 	if err != nil {
// 		return "", err
// 	}

// 	return signedToken, nil
// }

func GenerateToken(claims jwt.MapClaims, expiresAt time.Duration, jwtSecret string) (string, error) {
	claims["exp"] = time.Now().Add(expiresAt).Unix()
	claims["iat"] = time.Now().UTC().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, secretKey []byte) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("invalid token claims")
	}
	return token, claims, nil

}
