package main

import (
	"fmt"
	"net/http"

	"github.com/petersizovdev/MEDODS-T.git/internal/db"
	"github.com/petersizovdev/MEDODS-T.git/internal/handlers"
)


func main() {
	database:=db.Connect()
	if database == nil{
		panic("db is unable!")
	}
	defer database.Close()

	userHandler := &handlers.UserHandler{DB: database}

	http.HandleFunc("/", handlers.WelcomeHandler)
	http.HandleFunc("/users", userHandler.GetUsers)

	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	fmt.Println("Server is running on :3000")
	select {}
}
