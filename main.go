package main

import (
	service "crud-app/app/Service"
	"crud-app/app/repository"
	"crud-app/database"
	"crud-app/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	mongoClient := database.ConnectDB()
	db := database.GetDatabase(mongoClient)

	// repositories
	userRepo := repository.NewUserRepository(db)
	alumniRepo := repository.NewAlumniRepository(db)
	pekerjaanRepo := repository.NewPekerjaanRepository(db)

	// Service
	authService := service.NewAuthService(userRepo)
	alumniService := service.NewAlumniService(alumniRepo)
	PekerjaanService := service.NewPekerjaanService(pekerjaanRepo)
	userService := service.NewUserHandler(userRepo)
	r := mux.NewRouter()

	routes.UserRoutes(r, PekerjaanService, alumniService, authService, &userRepo, userService)
	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("🚀 Server running at http://localhost:" + port)
	http.ListenAndServe(":"+port, r)
}
