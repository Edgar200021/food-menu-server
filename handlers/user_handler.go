package handlers

import (
	"encoding/json"
	"fmt"
	"food-menu/storage"
	"food-menu/types"
	"food-menu/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"time"
)

type UserHandler struct {
	UserStorage storage.UserPgStorage
}

func (u *UserHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return

	}

	var user types.CreateUser

	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.UserStorage.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
func (u *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data := make(map[string]string, 2)
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if data["password"] == "" || data["email"] == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	user, dbErr := u.UserStorage.GetByEmail(data["email"])
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	var (
		accessToken  string
		refreshToken string
	)

	accessToken, err = utils.SignJWT([]byte(os.Getenv("JWT_SECRET")), jwt.MapClaims{
		"id":      user.ID,
		"expires": time.Now().Add(time.Minute * 30),
	})
	refreshToken, err = utils.SignJWT([]byte(os.Getenv("JWT_SECRET")), jwt.MapClaims{
		"id":      user.ID,
		"expires": time.Now().Add(time.Minute * 60 * 24 * 30),
	})

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   os.Getenv("GO_ENV") == "production",
		Expires:  time.Now().Add(time.Minute * 30),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   os.Getenv("GO_ENV") == "production",
		Expires:  time.Now().Add(time.Minute * 60 * 24 * 30),
		Path:     "/",
	})

	jsonResponse, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
func (u *UserHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	refreshToken, _ := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := refreshToken.Claims.(jwt.MapClaims), refreshToken.Valid; ok {

		layout := "2006-01-02T15:04:05.9999999-07:00"
		expiresDate, _ := time.Parse(layout, claims["expires"].(string))

		if time.Now().After(expiresDate) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		user, err := u.UserStorage.GetById(int(claims["id"].(float64)))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var (
			accessToken  string
			refreshToken string
		)

		accessToken, err = utils.SignJWT([]byte(os.Getenv("JWT_SECRET")), jwt.MapClaims{
			"id":      user.ID,
			"expires": time.Now().Add(time.Minute * 30),
		})
		refreshToken, err = utils.SignJWT([]byte(os.Getenv("JWT_SECRET")), jwt.MapClaims{
			"id":      user.ID,
			"expires": time.Now().Add(time.Minute * 60 * 24 * 30),
		})

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   os.Getenv("GO_ENV") == "production",
			Expires:  time.Now().Add(time.Minute * 30),
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   os.Getenv("GO_ENV") == "production",
			Expires:  time.Now().Add(time.Minute * 60 * 24 * 30),
			Path:     "/",
		})
	}
}
