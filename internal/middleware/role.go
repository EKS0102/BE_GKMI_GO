package middleware

import (
	"net/http"

	"BE_GKMI_NTC_GO/internal/response"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := GetRole(r.Context())
			if !ok {
				response.Error(
					w,
					http.StatusUnauthorized,
					"Role not found",
				)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(
				w,
				http.StatusForbidden,
				"Forbidden",
			)
		})
	}
}
