// auth middleware

package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtpkg "github.com/suthar345Piyush/financedashboard/pkg/jwt"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextRole   contextKey = "role"
	ContextEmail  contextKey = "email"
)

// authentication (using JWT to authenticate user)

func Authenticate(jwtManager *jwtpkg.Manager) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				response.Unauthorized(w)
				return
			}

			// spliting auth header  and getting middle substring b/w those seperators

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Unauthorized(w)
				return
			}

			//claims

			claims, err := jwtManager.ParseToken(parts[1])
			if err != nil {
				response.Unauthorized(w)
				return
			}

			// claims into request context (userid, role, email id)

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextRole, claims.Role)
			ctx = context.WithValue(ctx, ContextEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// function to get userId

func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(ContextUserID).(string); ok {
		return id
	}
	return ""
}

// getting role from context

func GetRole(r *http.Request) string {
	if role, ok := r.Context().Value(ContextRole).(string); ok {
		return role
	}

	return ""
}
