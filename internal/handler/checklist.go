package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/checklist"
	"github.com/snip/internal/repository"
)

type ChecklistHandler interface {
	CreateChecklist(title, description string, taskID, projectID *int) error
	CreateChecklistWithAI(topic, context string, numItems int, taskID, projectID *int) error
	ListChecklists(taskID, projectID *int) error
	ShowChecklist(id int) error
	AddChecklistItem(checklistID int, title, description string) error
	ToggleChecklistItem(id int) error
	DeleteChecklistItem(id int) error
	DeleteChecklist(id int) error
}

type checklistHandler struct {
	checklistRepo     repository.ChecklistRepository
	checklistItemRepo repository.ChecklistItemRepository
	groqClient        *ai.GroqClient
}

func NewChecklistHandler(checklistRepo repository.ChecklistRepository, checklistItemRepo repository.ChecklistItemRepository) ChecklistHandler {
	groqClient, _ := ai.NewGroqClient()
	return &checklistHandler{
		checklistRepo:     checklistRepo,
		checklistItemRepo:  checklistItemRepo,
		groqClient:         groqClient,
	}
}

func (h *checklistHandler) CreateChecklist(title, description string, taskID, projectID *int) error {
	c := checklist.NewChecklist(title, description)
	c.TaskID = taskID
	c.ProjectID = projectID

	if err := h.checklistRepo.Create(c); err != nil {
		return fmt.Errorf("failed to create checklist: %w", err)
	}

	fmt.Printf("Checklist criada com sucesso!\n")
	fmt.Printf("● #%d  %s\n", c.ID, c.Title)
	return nil
}

func (h *checklistHandler) CreateChecklistWithAI(topic, context string, numItems int, taskID, projectID *int) error {
	if h.groqClient == nil {
		return fmt.Errorf("AI client not available")
	}

	if numItems <= 0 {
		numItems = 5
	}

	fmt.Printf("Gerando checklist com IA (%d itens)...\n", numItems)
	items, err := h.groqClient.GenerateChecklist(topic, context, numItems)
	if err != nil {
		return fmt.Errorf("failed to generate checklist: %w", err)
	}

	c := checklist.NewChecklist(topic, context)
	c.TaskID = taskID
	c.ProjectID = projectID

	if err := h.checklistRepo.Create(c); err != nil {
		return fmt.Errorf("failed to create checklist: %w", err)
	}

	// Add items
	itemCount := 0
	for i, itemTitle := range items {
		item := checklist.NewChecklistItem(c.ID, itemTitle, "", i+1)
		if err := h.checklistItemRepo.Create(item); err != nil {
			return fmt.Errorf("failed to create checklist item: %w", err)
		}
		itemCount++
	}

	fmt.Printf("Checklist criada com sucesso!\n")
	fmt.Printf("● #%d  %s (%d itens)\n", c.ID, c.Title, itemCount)
	return nil
}

func (h *checklistHandler) ListChecklists(taskID, projectID *int) error {
	var checklists []*checklist.Checklist
	var err error

	if taskID != nil {
		checklists, err = h.checklistRepo.GetByTaskID(*taskID)
	} else if projectID != nil {
		checklists, err = h.checklistRepo.GetByProjectID(*projectID)
	} else {
		checklists, err = h.checklistRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to fetch checklists: %w", err)
	}

	if len(checklists) == 0 {
		fmt.Println("Nenhuma checklist encontrada.")
		return nil
	}

	fmt.Printf("Encontradas %d checklist(s):\n\n", len(checklists))
	for _, c := range checklists {
		fmt.Printf("● #%d %s\n", c.ID, c.Title)
		if c.Description != "" {
			desc := c.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			fmt.Printf("   └── %s\n", desc)
		}
		fmt.Println()
	}

	return nil
}

func (h *checklistHandler) ShowChecklist(id int) error {
	c, err := h.checklistRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch checklist: %w", err)
	}

	fmt.Printf("● #%d %s\n", c.ID, c.Title)
	if c.Description != "" {
		fmt.Printf("   └── %s\n\n", c.Description)
	}

	// Load items
	items, err := h.checklistItemRepo.GetByChecklistID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch checklist items: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("Nenhum item nesta checklist.")
		return nil
	}

	fmt.Printf("Itens (%d):\n", len(items))
	completedCount := 0
	for _, item := range items {
		icon := "○"
		if item.Completed {
			icon = "✓"
			completedCount++
		}
		fmt.Printf("  %s %s\n", icon, item.Title)
		if item.Description != "" {
			fmt.Printf("     └── %s\n", item.Description)
		}
	}

	fmt.Printf("\nProgresso: %d/%d concluído(s)\n", completedCount, len(items))
	return nil
}

func (h *checklistHandler) AddChecklistItem(checklistID int, title, description string) error {
	items, err := h.checklistItemRepo.GetByChecklistID(checklistID)
	if err != nil {
		return fmt.Errorf("failed to get checklist items: %w", err)
	}

	order := len(items) + 1
	item := checklist.NewChecklistItem(checklistID, title, description, order)
	if err := h.checklistItemRepo.Create(item); err != nil {
		return fmt.Errorf("failed to create checklist item: %w", err)
	}

	fmt.Printf("Item adicionado com sucesso!\n")
	fmt.Printf("  ○ %s\n", item.Title)
	return nil
}

func (h *checklistHandler) ToggleChecklistItem(id int) error {
	if err := h.checklistItemRepo.ToggleComplete(id); err != nil {
		return fmt.Errorf("failed to toggle checklist item: %w", err)
	}

	fmt.Printf("Item atualizado com sucesso!\n")
	return nil
}

func (h *checklistHandler) DeleteChecklistItem(id int) error {
	if err := h.checklistItemRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete checklist item: %w", err)
	}

	fmt.Printf("Item deletado com sucesso!\n")
	return nil
}

func (h *checklistHandler) DeleteChecklist(id int) error {
	if err := h.checklistRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete checklist: %w", err)
	}

	fmt.Printf("Checklist deletada com sucesso!\n")
	return nil
}

func parseIDs(idStr string) ([]int, error) {
	parts := strings.Split(idStr, ",")
	var ids []int
	for _, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %s", part)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

