package handler

import (
	"fmt"
	"strings"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/project"
	"github.com/snip/internal/repository"
)

type ProjectHandler interface {
	CreateProject(name, description string) error
	ListProjects(status string) error
	ShowProject(id int) error
	UpdateProject(id int, name, description, status string) error
	DeleteProject(id int) error
	CreateProjectWithAI(name, description string) error
}

type projectHandler struct {
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
	groqClient  *ai.GroqClient
}

func NewProjectHandler(projectRepo repository.ProjectRepository, taskRepo repository.TaskRepository) ProjectHandler {
	groqClient, _ := ai.NewGroqClient()
	return &projectHandler{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		groqClient:  groqClient,
	}
}

func (h *projectHandler) CreateProject(name, description string) error {
	p := project.NewProject(name, description)
	if err := h.projectRepo.Create(p); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Projeto criado com sucesso!\n")
	fmt.Printf("● #%d  %s\n", p.ID, p.Name)
	return nil
}

func (h *projectHandler) ListProjects(status string) error {
	projects, err := h.projectRepo.GetAll(status)
	if err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("Nenhum projeto encontrado.")
		return nil
	}

	fmt.Printf("Encontrados %d projeto(s):\n\n", len(projects))
	for _, p := range projects {
		fmt.Printf("● #%d %s [%s]\n", p.ID, p.Name, p.Status)
		if p.Description != "" {
			desc := p.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			fmt.Printf("   └── %s\n", desc)
		}
		fmt.Println()
	}

	return nil
}

func (h *projectHandler) ShowProject(id int) error {
	p, err := h.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	fmt.Printf("● #%d %s [%s]\n", p.ID, p.Name, p.Status)
	if p.Description != "" {
		fmt.Printf("   └── %s\n\n", p.Description)
	}

	// Show tasks
	tasks, err := h.taskRepo.GetByProjectID(id, "")
	if err == nil && len(tasks) > 0 {
		fmt.Printf("Tarefas (%d):\n", len(tasks))
		for _, t := range tasks {
			statusIcon := "○"
			if t.Status == "completed" {
				statusIcon = "✓"
			} else if t.Status == "in_progress" {
				statusIcon = "◐"
			}
			fmt.Printf("  %s #%d %s [%s]\n", statusIcon, t.ID, t.Title, t.Priority)
		}
	}

	return nil
}

func (h *projectHandler) UpdateProject(id int, name, description, status string) error {
	if err := h.projectRepo.Update(id, name, description, status); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	fmt.Printf("Projeto atualizado com sucesso!\n")
	return nil
}

func (h *projectHandler) DeleteProject(id int) error {
	if err := h.projectRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	fmt.Printf("Projeto deletado com sucesso!\n")
	return nil
}

func (h *projectHandler) CreateProjectWithAI(name, description string) error {
	if h.groqClient == nil {
		return fmt.Errorf("AI client not available")
	}

	fmt.Println("Gerando plano de projeto com IA...")
	plan, err := h.groqClient.GenerateProjectPlan(name, description)
	if err != nil {
		return fmt.Errorf("failed to generate project plan: %w", err)
	}

	p := project.NewProject(name, description)
	if err := h.projectRepo.Create(p); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Projeto criado com sucesso!\n")
	fmt.Printf("● #%d  %s\n\n", p.ID, p.Name)
	fmt.Println("Plano gerado pela IA:")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(plan)

	return nil
}

