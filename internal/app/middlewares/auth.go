package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	config *config.Config
}

func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{config: config}
}

func (m *AuthMiddleware) Auth(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := m.config.AuthSecret

		userID := ""
		authCookie, err := r.Cookie("USER_ID")
		authSignCookie, _ := r.Cookie("USER_SIGN")
		if errors.Is(err, http.ErrNoCookie) || strings.TrimSpace(authCookie.Value) == "" ||
			!verifySignUserID(key, strings.Replace(authCookie.Value, "USER_ID", "", 1), strings.Replace(authSignCookie.Value, "USER_SIGN", "", 1)) {
			userID = generateUserID()
			userIDSign := signUserID(key, []byte(userID))

			http.SetCookie(w, &http.Cookie{Name: "USER_ID", Value: userID, HttpOnly: true, Secure: true, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "USER_SIGN", Value: userIDSign, HttpOnly: true, Secure: true, Path: "/"})
		} else {
			userID = strings.Replace(authCookie.Value, "USER_ID", "", 1)
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

func signUserID(key string, userID []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(userID)

	return hex.EncodeToString(h.Sum(nil))
}

func verifySignUserID(key, userID, sign string) bool {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(userID))

	return hmac.Equal([]byte(hex.EncodeToString(h.Sum(nil))), []byte(sign))
}
