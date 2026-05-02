package middleware

import (
	"context"
	"net/http"
	"strings"

	"warehouse-app/utils"
)

type key string

const UserContext key = "user"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}

		ctx := context.WithValue(r.Context(), UserContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleAllowed(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserContext).(*utils.Claims)

			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", 403)
		})
	}
}
