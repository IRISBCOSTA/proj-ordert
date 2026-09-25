package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/IRISBCOSTA/proj-ordert/internal/httpx"
)

type contextKey string

const userIDKey contextKey = "user_id"

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-TROCAR"
	}
	return []byte(secret)
}

// cria um JWT valido por 24 horas
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			httpx.Error(w, http.StatusUnauthorized, "Token não informado")
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret(), nil
		})
		if err != nil || !token.Valid {
			httpx.Error(w, http.StatusUnauthorized, "Token inválido ou expirado")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httpx.Error(w, http.StatusUnauthorized, "Token inválido")
			return
		}

		userID := int(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func UserID(r *http.Request) int {
	id, _ := r.Context().Value(userIDKey).(int)
	return id
}