package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"net/http"
	"strings"
	"time"
)

const TokenExp = time.Hour * 3

type AuthMiddleware struct {
	config *config.Config
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{config: config}
}

func (m *AuthMiddleware) Auth(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := m.config.AuthSecret

		userID := ""
		authCookie, err := r.Cookie("Authorization")
		if errors.Is(err, http.ErrNoCookie) || strings.TrimSpace(authCookie.Value) == "" ||
			!verifySignUserID(key, strings.Replace(authCookie.Value, "Authorization", "", 1)) {
			userID = generateUserID()
			userIDSign := signUserID(key, userID)
			w.Header().Set("Authorization", userIDSign)

			http.SetCookie(w, &http.Cookie{Name: "Authorization", Value: userIDSign, HttpOnly: true, Secure: true, Path: "/", Expires: time.Now().Add(TokenExp)})
		} else if strings.TrimSpace(r.Header.Get("Authorization")) != "" {
			userID = getUserID(key, r.Header.Get("Authorization"))
		} else {
			userID = getUserID(key, strings.Replace(authCookie.Value, "Authorization", "", 1))
		}

		r = r.WithContext(context.WithValue(r.Context(), models.ContextUserID, userID))

		nextHandler.ServeHTTP(w, r)
	})
}

func generateUserID() string {
	userID := make([]byte, 16)
	if _, err := rand.Read(userID); err != nil {
		panic(err)
	}

	return hex.EncodeToString(userID)
}

func signUserID(key, userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return ""
	}

	return tokenString
}

func verifySignUserID(key, token string) bool {
	return strings.TrimSpace(getUserID(key, token)) != ""
}

func getUserID(key, tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.UserID
}
