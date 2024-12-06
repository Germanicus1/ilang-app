package routes

import (
	"backend/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", handlers.CreateUserHandler)
    mux.HandleFunc("GET /users/{id}", handlers.GetUserByIDHandler)
	mux.HandleFunc("PATCH /users/{id}", handlers.UpdateUserByIDHandler)
	mux.HandleFunc("DELETE /users/{id}", handlers. DeleteUserByIDHandler)
}
