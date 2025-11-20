package dto

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"isActive"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}
