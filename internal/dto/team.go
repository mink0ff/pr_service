package dto

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddUserToTeamRequest struct {
	UserID int `json:"userId" binding:"required"`
}

type TeamResponse struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Users []UserResponse `json:"users"`
}
