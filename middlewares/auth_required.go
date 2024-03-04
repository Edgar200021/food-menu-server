package middlewares

import (
	"context"
	"fmt"
	"food-menu/storage"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"reflect"
	"time"
)

func AuthRequired(handler http.HandlerFunc, userStorage storage.UserStorage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("accessToken")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, _ := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected string method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil

		})

		if claims, ok := accessToken.Claims.(jwt.MapClaims), accessToken.Valid; ok {

			layout := "2006-01-02T15:04:05.9999999-07:00"
			expiresDate, _ := time.Parse(layout, claims["expires"].(string))

			if time.Now().After(expiresDate) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			fmt.Println(reflect.TypeOf(claims["id"]))

			user, err := userStorage.GetById(int(claims["id"].(float64)))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(context.Background(), "user", user)
			handler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

	}
}
