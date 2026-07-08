package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

const UserKey contextKey = "user"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            header := r.Header.Get("Authorization")
            if !strings.HasPrefix(header, "Bearer ") {
                http.Error(w, "missing token", http.StatusUnauthorized)
                return
            }

            tokenStr := strings.TrimPrefix(header, "Bearer ")
            claims := &Claims{}

            token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
                return []byte(jwtSecret), nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), UserKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func UserFromContext(ctx context.Context) *Claims {
    c, _ := ctx.Value(UserKey).(*Claims)
    return c
}