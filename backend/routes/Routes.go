package routes

import (
	"backend/handlers"
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux) {
	// Register the Create User endpoint
	mux.HandleFunc("/users", handlers.CreateUserHandler)
	 // Register the dynamic route for /users/{id}
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/users/") {
            handlers.GetUserByIDHandler(w, r)
            return
        }
        http.NotFound(w, r)
    })
}
