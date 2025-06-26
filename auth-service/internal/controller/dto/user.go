package dto

// CreateUserRequest represents the request body for user registration
type CreateUserRequest struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// TokenResponse represents the response containing the JWT token
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string `json:"username" example:"john_doe"`
	CreatedAt string `json:"created_at" example:"2025-06-26T11:46:00Z"`
}
