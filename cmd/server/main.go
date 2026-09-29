package main

import (
	"context"
	"log"
	"net/http"

	"BE_GKMI_NTC_GO/internal/database"
	"BE_GKMI_NTC_GO/internal/middleware"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/router"
	"BE_GKMI_NTC_GO/internal/service"
)

// @title BE GKMI NTC API
// @version 1.0
// @description Backend API for GKMI NTC
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the JWT access token with the Bearer prefix. Example: Bearer {token}
func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	log.Println("PostgreSQL connection successful")

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	jemaatRepository := repository.NewJemaatRepository(db)

	authService := service.NewAuthService(userRepository, refreshTokenRepository)
	jemaatService := service.NewJemaatService(jemaatRepository)

	r := router.New(userRepository, authService, jemaatRepository, jemaatService)
	httpHandler := middleware.Logger(middleware.Recovery(r))

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}

	log.Println("BE_GKMI_NTC_GO server running on http://localhost:8080")

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
