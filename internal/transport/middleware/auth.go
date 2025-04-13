package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserTypeKey contextKey = "user_type"
)

var secret = []byte("super-secret")

type RoleCtxKey struct{}

func GetUserRole(ctx context.Context) string {
	role, _ := ctx.Value(RoleCtxKey{}).(string)
	return role
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		userType, _ := claims["user_type"].(string)
		ctx := context.WithValue(r.Context(), UserTypeKey, userType)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
