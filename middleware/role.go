// Middleware for Role of the User (Viewer, Analyst, Admin)
// Role based access control ( middleware | guard)

// Roles are on basis of hierarchial : Admin > Analyst > Viewer

package middleware

import (
	"net/http"

	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
)

// require role function will allow users with sufficient role level
// getting from our Role model (schema)

func RequireRole(roles ...models.Role) func(http.Handler) http.Handler {

	allowed := make(map[models.Role]bool)

	for _, r := range roles {
		allowed[r] = true
	}

	// returning the http handler at the end

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userRole := models.Role(GetRole(r))

			if userRole == "" {
				response.Unauthorized(w)
				return
			}

			if !allowed[userRole] {
				response.Forbidden(w)
				return
			}

			next.ServeHTTP(w, r)

		})
	}

}

// require minimum role means it allows the given role OR any higher role

func RequireMinRole(minimum models.Role) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := models.Role(GetRole(r))

			if userRole == "" {
				response.Unauthorized(w)
				return
			}

			if !userRole.HasPermission(minimum) {
				response.Forbidden(w)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
