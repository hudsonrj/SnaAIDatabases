package checklist

import "time"

// Checklist representa um checklist para projetos/tarefas
type Checklist struct {
	ID          int
	TaskID      *int
	ProjectID   *int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ChecklistItem representa um item de checklist
type ChecklistItem struct {
	ID          int
	ChecklistID int
	Title       string
	Description string
	Completed   bool
	Order       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewChecklist cria um novo checklist
func NewChecklist(title, description string) *Checklist {
	now := time.Now()
	return &Checklist{
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewChecklistItem cria um novo item de checklist
func NewChecklistItem(checklistID int, title, description string, order int) *ChecklistItem {
	now := time.Now()
	return &ChecklistItem{
		ChecklistID: checklistID,
		Title:       title,
		Description: description,
		Completed:   false,
		Order:       order,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

