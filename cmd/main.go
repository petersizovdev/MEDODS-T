package main

import (
	"fmt"
	"net/http"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Server 3000")
}

func main() {
	http.HandleFunc("/", welcomeHandler)

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
