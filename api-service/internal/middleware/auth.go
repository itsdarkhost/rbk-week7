package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const currentUserKey contextKey = "current_user"

type CurrentUser struct {
	Id    int
	Email string
	Role  string
}

func Auth(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(w, http.StatusUnauthorized, errors.New("authorization header is required"))
				return
			}

			tokenValue := strings.TrimPrefix(header, "Bearer ")
			if tokenValue == header || tokenValue == "" {
				writeError(w, http.StatusUnauthorized, errors.New("bearer token is required"))
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return jwtSecret, nil
			}, jwt.WithExpirationRequired())
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, errors.New("invalid token"))
				return
			}

			user, err := userFromClaims(claims)
			if err != nil {
				writeError(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), currentUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, errors.New("user is required"))
				return
			}
			if user.Role != role {
				writeError(w, http.StatusForbidden, errors.New("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) (CurrentUser, bool) {
	user, ok := ctx.Value(currentUserKey).(CurrentUser)
	return user, ok
}

func userFromClaims(claims jwt.MapClaims) (CurrentUser, error) {
	userId, ok := claims["user_id"].(float64)
	if !ok || userId <= 0 {
		return CurrentUser{}, errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return CurrentUser{}, errors.New("invalid token claims")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return CurrentUser{}, errors.New("invalid token claims")
	}

	return CurrentUser{
		Id:    int(userId),
		Email: email,
		Role:  role,
	}, nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
