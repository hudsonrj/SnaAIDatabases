package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snip/internal/jira"
	"github.com/spf13/cobra"
)

var (
	jiraConfigURL      string
	jiraConfigEmail    string
	jiraConfigAPIToken string
	jiraConfigProject  string
	jiraConfigShow     bool
)

func init() {
	jiraConfigCmd.Flags().StringVar(&jiraConfigURL, "url", "", "URL do Jira (ex: https://empresa.atlassian.net)")
	jiraConfigCmd.Flags().StringVar(&jiraConfigEmail, "email", "", "Email do usuário")
	jiraConfigCmd.Flags().StringVar(&jiraConfigAPIToken, "api-token", "", "API Token do Jira")
	jiraConfigCmd.Flags().StringVar(&jiraConfigProject, "project", "", "Chave do projeto (ex: PROJ)")
	jiraConfigCmd.Flags().BoolVar(&jiraConfigShow, "show", false, "Mostrar configuração atual")

	jiraCreateEpicCmd.Flags().StringVarP(&jiraEpicSummary, "summary", "s", "", "Resumo do Epic")
	jiraCreateEpicCmd.Flags().StringVarP(&jiraEpicDescription, "description", "d", "", "Descrição do Epic")

	jiraCreateIssueCmd.Flags().StringVarP(&jiraIssueSummary, "summary", "s", "", "Resumo da Issue")
	jiraCreateIssueCmd.Flags().StringVarP(&jiraIssueDescription, "description", "d", "", "Descrição da Issue")
	jiraCreateIssueCmd.Flags().StringVarP(&jiraIssueType, "type", "t", "Task", "Tipo da Issue (Task, Bug, Story, etc.)")
	jiraCreateIssueCmd.Flags().StringVarP(&jiraEpicKey, "epic", "e", "", "Chave do Epic (opcional)")

	rootCmd.AddCommand(jiraCmd)
	jiraCmd.AddCommand(jiraConfigCmd)
	jiraCmd.AddCommand(jiraCreateEpicCmd)
	jiraCmd.AddCommand(jiraCreateIssueCmd)
}

var (
	jiraEpicSummary     string
	jiraEpicDescription string
	jiraIssueSummary    string
	jiraIssueDescription string
	jiraIssueType       string
	jiraEpicKey         string
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Integração com Jira",
	Long: `Comandos para integração com Jira da empresa.

Permite criar Epicos e Issues (cards) no Jira a partir de análises e problemas identificados.`,
}

var jiraConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configurar integração com Jira",
	Long: `Configura a integração com Jira.

Exemplos:
  snip jira config --url "https://empresa.atlassian.net" --email "usuario@empresa.com" --api-token "token" --project "PROJ"
  snip jira config --show  # Mostrar configuração atual
  snip jira config         # Modo interativo`,
	Run: func(cmd *cobra.Command, args []string) {
		if jiraConfigShow {
			showJiraConfig()
			return
		}

		config, err := jira.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		// Modo interativo se nenhum parâmetro foi fornecido
		if jiraConfigURL == "" && jiraConfigEmail == "" && jiraConfigAPIToken == "" && jiraConfigProject == "" {
			interactiveJiraConfig(config)
			return
		}

		// Atualizar configuração
		if jiraConfigURL != "" {
			config.URL = jiraConfigURL
		}
		if jiraConfigEmail != "" {
			config.Email = jiraConfigEmail
		}
		if jiraConfigAPIToken != "" {
			config.APIToken = jiraConfigAPIToken
		}
		if jiraConfigProject != "" {
			config.Project = jiraConfigProject
		}

		if err := jira.SaveConfig(config); err != nil {
			fmt.Printf("Erro ao salvar configuração: %v\n", err)
			return
		}

		fmt.Println("✓ Configuração do Jira salva com sucesso!")
		showJiraConfig()
	},
}

var jiraCreateEpicCmd = &cobra.Command{
	Use:   "create-epic",
	Short: "Criar um Epic no Jira",
	Long: `Cria um Epic no Jira.

Exemplos:
  snip jira create-epic --summary "Problemas de Performance PostgreSQL" --description "Epic para resolver problemas identificados"
  snip jira create-epic -s "Otimização Banco" -d "Melhorias necessárias"`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := jira.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Project == "" {
			fmt.Println("❌ Configuração do Jira incompleta. Execute: snip jira config")
			return
		}

		if jiraEpicSummary == "" {
			fmt.Println("❌ Resumo do Epic é obrigatório (use --summary)")
			return
		}

		client, err := jira.NewClient(config)
		if err != nil {
			fmt.Printf("Erro ao criar cliente Jira: %v\n", err)
			return
		}

		description := jiraEpicDescription
		if description == "" {
			description = "Epic criado automaticamente pelo SnipAI"
		}

		epic, err := client.CreateEpic(jiraEpicSummary, description)
		if err != nil {
			fmt.Printf("Erro ao criar Epic: %v\n", err)
			return
		}

		fmt.Printf("✅ Epic criado com sucesso!\n")
		fmt.Printf("  Key: %s\n", epic.Key)
		fmt.Printf("  Título: %s\n", epic.Fields.Summary)
		fmt.Printf("  URL: %s/browse/%s\n", config.URL, epic.Key)
	},
}

var jiraCreateIssueCmd = &cobra.Command{
	Use:   "create-issue",
	Short: "Criar uma Issue (card) no Jira",
	Long: `Cria uma Issue (card) no Jira, opcionalmente vinculada a um Epic.

Exemplos:
  snip jira create-issue --summary "Corrigir configuração shared_buffers" --description "Ajustar parâmetro" --epic "PROJ-123"
  snip jira create-issue -s "Tarefa" -d "Descrição" -t "Bug" -e "PROJ-100"`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := jira.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Project == "" {
			fmt.Println("❌ Configuração do Jira incompleta. Execute: snip jira config")
			return
		}

		if jiraIssueSummary == "" {
			fmt.Println("❌ Resumo da Issue é obrigatório (use --summary)")
			return
		}

		client, err := jira.NewClient(config)
		if err != nil {
			fmt.Printf("Erro ao criar cliente Jira: %v\n", err)
			return
		}

		description := jiraIssueDescription
		if description == "" {
			description = "Issue criada automaticamente pelo SnipAI"
		}

		issue, err := client.CreateIssue(jiraIssueSummary, description, jiraIssueType, jiraEpicKey)
		if err != nil {
			fmt.Printf("Erro ao criar Issue: %v\n", err)
			return
		}

		fmt.Printf("✅ Issue criada com sucesso!\n")
		fmt.Printf("  Key: %s\n", issue.Key)
		fmt.Printf("  Título: %s\n", issue.Fields.Summary)
		fmt.Printf("  Tipo: %s\n", issue.Fields.IssueType.Name)
		if jiraEpicKey != "" {
			fmt.Printf("  Epic: %s\n", jiraEpicKey)
		}
		fmt.Printf("  URL: %s/browse/%s\n", config.URL, issue.Key)
	},
}

func showJiraConfig() {
	config, err := jira.LoadConfig()
	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		return
	}

	fmt.Println("📋 Configuração Atual do Jira:")
	fmt.Printf("  URL: %s\n", config.URL)
	fmt.Printf("  Email: %s\n", config.Email)
	tokenDisplay := config.APIToken
	if len(tokenDisplay) > 8 {
		tokenDisplay = tokenDisplay[:4] + "..." + tokenDisplay[len(tokenDisplay)-4:]
	}
	fmt.Printf("  API Token: %s\n", tokenDisplay)
	fmt.Printf("  Projeto: %s\n", config.Project)
}

func interactiveJiraConfig(config *jira.JiraConfig) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Configuração Interativa do Jira ---")

	fmt.Printf("URL do Jira atual (%s): ", config.URL)
	urlInput, _ := reader.ReadString('\n')
	urlInput = strings.TrimSpace(urlInput)
	if urlInput != "" {
		config.URL = urlInput
	}

	fmt.Printf("Email atual (%s): ", config.Email)
	emailInput, _ := reader.ReadString('\n')
	emailInput = strings.TrimSpace(emailInput)
	if emailInput != "" {
		config.Email = emailInput
	}

	tokenDisplay := config.APIToken
	if len(tokenDisplay) > 8 {
		tokenDisplay = tokenDisplay[:4] + "..." + tokenDisplay[len(tokenDisplay)-4:]
	}
	fmt.Printf("API Token atual (%s): ", tokenDisplay)
	tokenInput, _ := reader.ReadString('\n')
	tokenInput = strings.TrimSpace(tokenInput)
	if tokenInput != "" {
		config.APIToken = tokenInput
	}

	fmt.Printf("Chave do Projeto atual (%s): ", config.Project)
	projectInput, _ := reader.ReadString('\n')
	projectInput = strings.TrimSpace(projectInput)
	if projectInput != "" {
		config.Project = projectInput
	}

	if err := jira.SaveConfig(config); err != nil {
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}
	fmt.Println("✓ Configuração salva com sucesso!")
	showJiraConfig()
}

