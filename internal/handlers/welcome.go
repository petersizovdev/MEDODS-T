package handlers

import (
	"fmt"
	"net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Server 3000")
}
