package project

import "time"

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // active, completed, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProject(name, description string) *Project {
	now := time.Now()
	return &Project{
		Name:        name,
		Description: description,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

