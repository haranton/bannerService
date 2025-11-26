package middleware

import (
	"bannerService/internals/handlers"
	"context"
	"net/http"
)

type roleKey struct{}

const (
	RoleUser   = "user"
	RoleAdmin  = "admin"
	UserToken  = "user_token"
	AdminToken = "admin_token"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			handlers.WriteJSONError(w, http.StatusUnauthorized, "missing token")
			return
		}

		var role string

		switch token {
		case UserToken:
			role = RoleUser
		case AdminToken:
			role = AdminToken
		default:
			handlers.WriteJSONError(w, http.StatusForbidden, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), roleKey{}, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role, ok := r.Context().Value(roleKey{}).(string)
		if !ok || role != RoleAdmin {
			handlers.WriteJSONError(w, http.StatusForbidden, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
