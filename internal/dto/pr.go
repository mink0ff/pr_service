package dto

type CreatePRRequest struct {
	Title    string `json:"title" binding:"required"`
	AuthorID int    `json:"authorId" binding:"required"`
}

type PRResponse struct {
	ID        int            `json:"id"`
	Title     string         `json:"title"`
	AuthorID  int            `json:"authorId"`
	TeamID    int            `json:"teamId"`
	Status    string         `json:"status"`
	Reviewers []UserResponse `json:"reviewers"`
}

type ReassignReviewerRequest struct {
	ReviewerID int `json:"reviewerId" binding:"required"`
}
