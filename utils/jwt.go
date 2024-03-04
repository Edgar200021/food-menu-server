package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(secretKey []byte, claims jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return s, nil
}
