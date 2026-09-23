package router

import (
	"net/http"

	"BE_GKMI_NTC_GO/internal/handler"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/service"
)

func New(
	userRepository *repository.UserRepository,
	authService *service.AuthService,
	jemaatRepository *repository.JemaatRepository,
	jemaatService *service.JemaatService,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/health",
		handler.Health,
	)

	mux.HandleFunc(
		"/api/auth/login",
		handler.Login(authService),
	)

	mux.HandleFunc(
		"/api/users",
		handler.Users(userRepository),
	)

	mux.HandleFunc(
		"/api/users/",
		handler.GetUserByID(userRepository),
	)

	mux.HandleFunc(
		"/api/jemaat",
		handler.Jemaat(
			jemaatRepository,
			jemaatService,
		),
	)

	mux.HandleFunc(
		"/api/jemaat/",
		handler.Jemaat(jemaatRepository, jemaatService),
	)

	return mux
}
