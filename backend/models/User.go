package models

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
