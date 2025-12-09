package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "snip",
	Short: "A fast and lightweight note-taking CLI application with AI and project management",
	Long: `Snip é uma aplicação completa de gerenciamento de notas, projetos, tarefas e checklists.

Funcionalidades principais:
  📝 Notas: Crie, edite, busque e organize suas notas
  🤖 IA: Gere conteúdo, código, checklists e planejamentos com IA
  📁 Projetos: Organize seus projetos e tarefas
  ✅ Checklists: Crie listas de verificação e acompanhe o progresso
  🏷️ Tags: Organize tudo com tags personalizadas

Exemplos rápidos:
  snip create "Minha Nota"
  snip ai-create "Python Básico" --tag "programming"
  snip project create "Meu Projeto"
  snip task create "Nova Tarefa" --project 1
  snip checklist ai-create "Preparação" --items 5`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(patchCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editorCmd)
	rootCmd.AddCommand(recentCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}
