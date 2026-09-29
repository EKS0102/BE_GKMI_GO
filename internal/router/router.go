package router

import (
	"net/http"

	_ "BE_GKMI_NTC_GO/docs"
	"BE_GKMI_NTC_GO/internal/handler"
	"BE_GKMI_NTC_GO/internal/middleware"
	"BE_GKMI_NTC_GO/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func New(userRepository *repository.UserRepository, authService handler.AuthServiceInterface, jemaatRepository *repository.JemaatRepository, jemaatService handler.JemaatServiceInterface) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", handler.Health)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.Handle("/api/auth/login", handler.Login(authService))
	mux.Handle("/api/auth/refresh", handler.Refresh(authService))
	mux.Handle("/api/auth/logout", handler.Logout(authService))

	mux.Handle("/api/auth/me", middleware.JWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		username, ok := middleware.GetUsername(r.Context())
		if !ok {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}

		role, ok := middleware.GetRole(r.Context())
		if !ok {
			http.Error(w, "Invalid role", http.StatusUnauthorized)
			return
		}

		handler.Me(userID, username, role).ServeHTTP(w, r)
	})))

	mux.Handle("/api/users", middleware.JWT(handler.Users(userRepository)))
	mux.Handle("/api/users/", middleware.JWT(handler.GetUserByID(userRepository)))

	jemaatHandler := handler.Jemaat(jemaatRepository, jemaatService)

	mux.Handle("GET /api/jemaat", middleware.JWT(middleware.RequireRole("viewer", "staff", "admin")(jemaatHandler)))
	mux.Handle("GET /api/jemaat/{id}", middleware.JWT(middleware.RequireRole("viewer", "staff", "admin")(handler.GetJemaatByID(jemaatService))))
	mux.Handle("POST /api/jemaat", middleware.JWT(middleware.RequireRole("staff", "admin")(jemaatHandler)))
	mux.Handle("POST /api/jemaat/bulk", middleware.JWT(middleware.RequireRole("staff", "admin")(jemaatHandler)))
	mux.Handle("PUT /api/jemaat/{id}", middleware.JWT(middleware.RequireRole("staff", "admin")(jemaatHandler)))
	mux.Handle("DELETE /api/jemaat/{id}", middleware.JWT(middleware.RequireRole("admin")(jemaatHandler)))

	return mux
}
