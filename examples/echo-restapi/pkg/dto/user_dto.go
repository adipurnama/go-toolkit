package dto

// CreateUserRequest dto for create single user's request.
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// CreateUserResponse is dto for create single user's response.
type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
