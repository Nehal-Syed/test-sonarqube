package main

import (
	"fmt"
	"log"
	"net/http"

	"test-sonarqube/handler"
	"test-sonarqube/repository"
	"test-sonarqube/service"
)

func main() {
	// Initialize dependencies
	userRepo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Setup routes
	http.HandleFunc("POST /users", userHandler.CreateUser)
	http.HandleFunc("GET /users", userHandler.ListUsers)
	http.HandleFunc("GET /users/", userHandler.GetUser)
	http.HandleFunc("PUT /users/", userHandler.UpdateUser)
	http.HandleFunc("DELETE /users/", userHandler.DeleteUser)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
