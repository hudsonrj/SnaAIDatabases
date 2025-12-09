package task

import "time"

type Task struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // pending, in_progress, completed
	Priority    string    `json:"priority"` // low, medium, high
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTask(projectID int, title, description, priority string) *Task {
	now := time.Now()
	return &Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

