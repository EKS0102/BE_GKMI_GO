package middleware

import (
	"log"
	"net/http"

	"BE_GKMI_NTC_GO/internal/response"
)

func Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				log.Printf("PANIC: %v", err)

				response.Error(
					w,
					http.StatusInternalServerError,
					"Internal Server Error",
				)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
