package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/petersizovdev/MEDODS-T.git/internal/db"
	"github.com/petersizovdev/MEDODS-T.git/internal/handlers"
	"github.com/petersizovdev/MEDODS-T.git/pkg/env"
)

func main() {
	err := env.LoadEnv(".env")
	if err != nil {
		fmt.Println("Err to load .env", err)
	}
	port := os.Getenv("PORT")

	database := db.Connect()
	if database == nil {
		panic("db is unable!")
	}
	defer database.Close()

	userHandler := &handlers.UserHandler{DB: database}
	authHandler := &handlers.AuthHandler{DB: database} 

	http.HandleFunc("/", handlers.WelcomeHandler)
	http.HandleFunc("/users", userHandler.GetUsers)
	http.HandleFunc("/token", authHandler.GenerateTokens)
	http.HandleFunc("/refresh", authHandler.RefreshTokens)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	fmt.Printf("Server is running on :%s \n", port)
	<-make(chan struct{})
}
