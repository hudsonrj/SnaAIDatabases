package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var taskDescription string
var taskStatus string
var taskPriority string
var taskDueDate string
var taskProjectID int

func init() {
	taskCreateCmd.Flags().StringVarP(&taskDescription, "description", "d", "", "Descrição da tarefa")
	taskCreateCmd.Flags().StringVarP(&taskPriority, "priority", "p", "medium", "Prioridade (low, medium, high)")
	taskCreateCmd.Flags().StringVarP(&taskDueDate, "due", "", "", "Data de vencimento (YYYY-MM-DD)")
	taskCreateCmd.Flags().IntVarP(&taskProjectID, "project", "", 0, "ID do projeto")
	
	taskListCmd.Flags().StringVarP(&taskStatus, "status", "s", "", "Filtrar por status (pending, in_progress, completed)")
	taskListCmd.Flags().IntVarP(&taskProjectID, "project", "", 0, "ID do projeto")
	
	taskUpdateCmd.Flags().StringVarP(&taskDescription, "description", "d", "", "Descrição da tarefa")
	taskUpdateCmd.Flags().StringVarP(&taskStatus, "status", "s", "", "Status (pending, in_progress, completed)")
	taskUpdateCmd.Flags().StringVarP(&taskPriority, "priority", "p", "", "Prioridade (low, medium, high)")
	taskUpdateCmd.Flags().StringVarP(&taskDueDate, "due", "", "", "Data de vencimento (YYYY-MM-DD)")
	
	rootCmd.AddCommand(taskCmd)
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Gerenciar tarefas",
	Long:  `Gerencie tarefas dos seus projetos.`,
}

var taskCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Criar uma nova tarefa",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			title := strings.Join(args, " ")
			var dueDate *time.Time
			if taskDueDate != "" {
				parsed, err := time.Parse("2006-01-02", taskDueDate)
				if err != nil {
					return fmt.Errorf("data inválida: %w", err)
				}
				dueDate = &parsed
			}
			return h.CreateTask(taskProjectID, title, taskDescription, taskPriority, dueDate)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar tarefas",
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			projectID := 0
			if taskProjectID > 0 {
				projectID = taskProjectID
			}
			return h.ListTasks(projectID, taskStatus)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var taskShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Mostrar detalhes de uma tarefa",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.ShowTask(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update [id] [title]",
	Short: "Atualizar uma tarefa",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			title := strings.Join(args[1:], " ")
			var dueDate *time.Time
			if taskDueDate != "" {
				parsed, err := time.Parse("2006-01-02", taskDueDate)
				if err != nil {
					return fmt.Errorf("data inválida: %w", err)
				}
				dueDate = &parsed
			}
			return h.UpdateTask(id, title, taskDescription, taskStatus, taskPriority, dueDate)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Deletar uma tarefa",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.DeleteTask(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var taskToggleCmd = &cobra.Command{
	Use:   "toggle [id]",
	Short: "Marcar/desmarcar tarefa como concluída",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithTaskHandler(func(h handler.TaskHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.ToggleTaskComplete(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

func init() {
	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskToggleCmd)
}

