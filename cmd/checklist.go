package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var checklistDescription string
var checklistTaskID int
var checklistProjectID int
var checklistNumItems int
var checklistItemDescription string

func init() {
	checklistCreateCmd.Flags().StringVarP(&checklistDescription, "description", "d", "", "Descrição da checklist")
	checklistCreateCmd.Flags().IntVarP(&checklistTaskID, "task", "", 0, "ID da tarefa")
	checklistCreateCmd.Flags().IntVarP(&checklistProjectID, "project", "", 0, "ID do projeto")
	
	checklistAICreateCmd.Flags().StringVarP(&checklistDescription, "description", "d", "", "Contexto para geração")
	checklistAICreateCmd.Flags().IntVarP(&checklistNumItems, "items", "n", 5, "Número de itens")
	checklistAICreateCmd.Flags().IntVarP(&checklistTaskID, "task", "", 0, "ID da tarefa")
	checklistAICreateCmd.Flags().IntVarP(&checklistProjectID, "project", "", 0, "ID do projeto")
	
	checklistListCmd.Flags().IntVarP(&checklistTaskID, "task", "", 0, "ID da tarefa")
	checklistListCmd.Flags().IntVarP(&checklistProjectID, "project", "", 0, "ID do projeto")
	
	checklistItemAddCmd.Flags().StringVarP(&checklistItemDescription, "description", "d", "", "Descrição do item")
	
	rootCmd.AddCommand(checklistCmd)
}

var checklistCmd = &cobra.Command{
	Use:   "checklist",
	Short: "Gerenciar checklists",
	Long:  `Gerencie checklists e seus itens.`,
}

var checklistCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Criar uma nova checklist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			title := strings.Join(args, " ")
			var taskID, projectID *int
			if checklistTaskID > 0 {
				taskID = &checklistTaskID
			}
			if checklistProjectID > 0 {
				projectID = &checklistProjectID
			}
			return h.CreateChecklist(title, checklistDescription, taskID, projectID)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistAICreateCmd = &cobra.Command{
	Use:   "ai-create [topic]",
	Short: "Criar checklist com itens gerados por IA",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			topic := strings.Join(args, " ")
			var taskID, projectID *int
			if checklistTaskID > 0 {
				taskID = &checklistTaskID
			}
			if checklistProjectID > 0 {
				projectID = &checklistProjectID
			}
			return h.CreateChecklistWithAI(topic, checklistDescription, checklistNumItems, taskID, projectID)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar checklists",
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			var taskID, projectID *int
			if checklistTaskID > 0 {
				taskID = &checklistTaskID
			}
			if checklistProjectID > 0 {
				projectID = &checklistProjectID
			}
			return h.ListChecklists(taskID, projectID)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Mostrar detalhes de uma checklist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.ShowChecklist(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Deletar uma checklist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.DeleteChecklist(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistItemAddCmd = &cobra.Command{
	Use:   "item-add [checklist_id] [title]",
	Short: "Adicionar item a uma checklist",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			checklistID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			title := strings.Join(args[1:], " ")
			return h.AddChecklistItem(checklistID, title, checklistItemDescription)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistItemToggleCmd = &cobra.Command{
	Use:   "item-toggle [item_id]",
	Short: "Marcar/desmarcar item como concluído",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.ToggleChecklistItem(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var checklistItemDeleteCmd = &cobra.Command{
	Use:   "item-delete [item_id]",
	Short: "Deletar um item de checklist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithChecklistHandler(func(h handler.ChecklistHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.DeleteChecklistItem(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

func init() {
	checklistCmd.AddCommand(checklistCreateCmd)
	checklistCmd.AddCommand(checklistAICreateCmd)
	checklistCmd.AddCommand(checklistListCmd)
	checklistCmd.AddCommand(checklistShowCmd)
	checklistCmd.AddCommand(checklistDeleteCmd)
	checklistCmd.AddCommand(checklistItemAddCmd)
	checklistCmd.AddCommand(checklistItemToggleCmd)
	checklistCmd.AddCommand(checklistItemDeleteCmd)
}

