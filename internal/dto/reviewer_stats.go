package dto

type ReviewerStatsItem struct {
	UserID string `json:"user_id"`
	Count  int64  `json:"count"`
}

type ReviewerStatsResponse struct {
	Items []ReviewerStatsItem `json:"items"`
}
