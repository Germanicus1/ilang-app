package routes

import (
	"backend/handlers"
	"fmt"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request caught:", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not Found"))
})

	// Register the Create User endpoint
	mux.HandleFunc("/users", handlers.CreateUserHandler)
}
