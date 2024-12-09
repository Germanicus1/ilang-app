package routes

import (
	"backend/handlers"
	"net/http"
)

func RegisterPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", handlers.CreateUserHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)
}
