package handler

import (
	"fmt"
	"time"

	"github.com/snip/internal/repository"
	"github.com/snip/internal/task"
)

type TaskHandler interface {
	CreateTask(projectID int, title, description, priority string, dueDate *time.Time) error
	ListTasks(projectID int, status string) error
	ShowTask(id int) error
	UpdateTask(id int, title, description, status, priority string, dueDate *time.Time) error
	DeleteTask(id int) error
	ToggleTaskComplete(id int) error
}

type taskHandler struct {
	taskRepo    repository.TaskRepository
	projectRepo repository.ProjectRepository
}

func NewTaskHandler(taskRepo repository.TaskRepository, projectRepo repository.ProjectRepository) TaskHandler {
	return &taskHandler{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (h *taskHandler) CreateTask(projectID int, title, description, priority string, dueDate *time.Time) error {
	if priority == "" {
		priority = "medium"
	}

	t := task.NewTask(projectID, title, description, priority)
	if dueDate != nil {
		t.DueDate = dueDate
	}

	if err := h.taskRepo.Create(t); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Printf("Tarefa criada com sucesso!\n")
	fmt.Printf("● #%d  %s [%s]\n", t.ID, t.Title, t.Priority)
	return nil
}

func (h *taskHandler) ListTasks(projectID int, status string) error {
	var tasks []*task.Task
	var err error

	if projectID > 0 {
		tasks, err = h.taskRepo.GetByProjectID(projectID, status)
	} else {
		tasks, err = h.taskRepo.GetAll(status)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("Nenhuma tarefa encontrada.")
		return nil
	}

	fmt.Printf("Encontradas %d tarefa(s):\n\n", len(tasks))
	for _, t := range tasks {
		statusIcon := "○"
		if t.Status == "completed" {
			statusIcon = "✓"
		} else if t.Status == "in_progress" {
			statusIcon = "◐"
		}

		fmt.Printf("%s #%d %s [%s]", statusIcon, t.ID, t.Title, t.Priority)
		if t.DueDate != nil {
			fmt.Printf(" (prazo: %s)", t.DueDate.Format("2006-01-02"))
		}
		fmt.Println()

		if t.Description != "" {
			desc := t.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			fmt.Printf("   └── %s\n", desc)
		}
		fmt.Println()
	}

	return nil
}

func (h *taskHandler) ShowTask(id int) error {
	t, err := h.taskRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch task: %w", err)
	}

	statusIcon := "○"
	if t.Status == "completed" {
		statusIcon = "✓"
	} else if t.Status == "in_progress" {
		statusIcon = "◐"
	}

	fmt.Printf("%s #%d %s [%s]\n", statusIcon, t.ID, t.Title, t.Priority)
	if t.Description != "" {
		fmt.Printf("   └── %s\n", t.Description)
	}
	if t.DueDate != nil {
		fmt.Printf("   └── Prazo: %s\n", t.DueDate.Format("2006-01-02 15:04"))
	}

	return nil
}

func (h *taskHandler) UpdateTask(id int, title, description, status, priority string, dueDate *time.Time) error {
	if err := h.taskRepo.Update(id, title, description, status, priority, dueDate); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	fmt.Printf("Tarefa atualizada com sucesso!\n")
	return nil
}

func (h *taskHandler) DeleteTask(id int) error {
	if err := h.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	fmt.Printf("Tarefa deletada com sucesso!\n")
	return nil
}

func (h *taskHandler) ToggleTaskComplete(id int) error {
	if err := h.taskRepo.ToggleComplete(id); err != nil {
		return fmt.Errorf("failed to toggle task: %w", err)
	}

	t, err := h.taskRepo.GetByID(id)
	if err != nil {
		return err
	}

	status := "pendente"
	if t.Status == "completed" {
		status = "concluída"
	}

	fmt.Printf("Tarefa marcada como %s!\n", status)
	return nil
}

