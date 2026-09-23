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

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	log.Println("PostgreSQL connection successful")

	userRepository := repository.NewUserRepository(db)
	jemaatRepository := repository.NewJemaatRepository(db)

	authService := service.NewAuthService(userRepository)
	jemaatService := service.NewJemaatService(jemaatRepository)

	r := router.New(
		userRepository,
		authService,
		jemaatRepository,
		jemaatService,
	)

	handler := middleware.Logger(
		middleware.Recovery(r),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("BE_GKMI_NTC_GO server running on http://localhost:8080")

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
