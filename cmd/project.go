package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var projectDescription string
var projectStatus string

func init() {
	projectCreateCmd.Flags().StringVarP(&projectDescription, "description", "d", "", "Descrição do projeto")
	projectUpdateCmd.Flags().StringVarP(&projectDescription, "description", "d", "", "Descrição do projeto")
	projectUpdateCmd.Flags().StringVarP(&projectStatus, "status", "s", "", "Status do projeto (active, completed, archived)")
	rootCmd.AddCommand(projectCmd)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Gerenciar projetos",
	Long:  `Gerencie seus projetos, tarefas e checklists.`,
}

var projectCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Criar um novo projeto",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			name := strings.Join(args, " ")
			return h.CreateProject(name, projectDescription)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar projetos",
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			return h.ListProjects(projectStatus)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var projectShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Mostrar detalhes de um projeto",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.ShowProject(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [id] [name]",
	Short: "Atualizar um projeto",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			name := strings.Join(args[1:], " ")
			return h.UpdateProject(id, name, projectDescription, projectStatus)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Deletar um projeto",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
			return h.DeleteProject(id)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var projectAICreateCmd = &cobra.Command{
	Use:   "ai-create [name]",
	Short: "Criar projeto com plano gerado por IA",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithProjectHandler(func(h handler.ProjectHandler) error {
			name := strings.Join(args, " ")
			return h.CreateProjectWithAI(name, projectDescription)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectShowCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectAICreateCmd)
}

