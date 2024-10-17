package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte = []byte("mybigsecretkey")

// hashes password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePasswords(plainPassword string, hashedPassword string) error {
	plainpwdBytes := []byte(plainPassword)

	hashedpwdBytes := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedpwdBytes, plainpwdBytes)

	if err != nil {
		return err
	}

	return nil
}

func GenerateJWT(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func DecodeJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["sub"].(string)

		return username, nil
	} else {
		return "", err
	}
}
