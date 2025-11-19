package models

type Team struct {
	ID   int
	Name string
}

type TeamUser struct {
	TeamID int
	UserID int
}
