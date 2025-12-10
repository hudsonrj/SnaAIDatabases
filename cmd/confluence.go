package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snip/internal/confluence"
	"github.com/spf13/cobra"
)

var (
	confluenceConfigURL      string
	confluenceConfigEmail     string
	confluenceConfigAPIToken  string
	confluenceConfigSpace     string
	confluenceConfigShow      bool
	confluencePageTitle       string
	confluencePageContent     string
	confluencePageParentID     string
)

func init() {
	confluenceConfigCmd.Flags().StringVar(&confluenceConfigURL, "url", "", "URL do Confluence (ex: https://empresa.atlassian.net)")
	confluenceConfigCmd.Flags().StringVar(&confluenceConfigEmail, "email", "", "Email do usuário")
	confluenceConfigCmd.Flags().StringVar(&confluenceConfigAPIToken, "api-token", "", "API Token do Confluence")
	confluenceConfigCmd.Flags().StringVar(&confluenceConfigSpace, "space", "", "Chave do espaço (ex: DB)")
	confluenceConfigCmd.Flags().BoolVar(&confluenceConfigShow, "show", false, "Mostrar configuração atual")

	confluenceCreatePageCmd.Flags().StringVarP(&confluencePageTitle, "title", "t", "", "Título da página")
	confluenceCreatePageCmd.Flags().StringVarP(&confluencePageContent, "content", "c", "", "Conteúdo da página (markdown ou HTML)")
	confluenceCreatePageCmd.Flags().StringVarP(&confluencePageParentID, "parent", "p", "", "ID da página pai (opcional)")

	rootCmd.AddCommand(confluenceCmd)
	confluenceCmd.AddCommand(confluenceConfigCmd)
	confluenceCmd.AddCommand(confluenceCreatePageCmd)
}

var confluenceCmd = &cobra.Command{
	Use:   "confluence",
	Short: "Integração com Confluence",
	Long: `Comandos para integração com Confluence da empresa.

Permite exportar documentações, configurações de banco de dados e análises para o Confluence.`,
}

var confluenceConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configurar integração com Confluence",
	Long: `Configura a integração com Confluence.

Exemplos:
  snip confluence config --url "https://empresa.atlassian.net" --email "usuario@empresa.com" --api-token "token" --space "DB"
  snip confluence config --show  # Mostrar configuração atual
  snip confluence config         # Modo interativo`,
	Run: func(cmd *cobra.Command, args []string) {
		if confluenceConfigShow {
			showConfluenceConfig()
			return
		}

		config, err := confluence.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		// Modo interativo se nenhum parâmetro foi fornecido
		if confluenceConfigURL == "" && confluenceConfigEmail == "" && confluenceConfigAPIToken == "" && confluenceConfigSpace == "" {
			interactiveConfluenceConfig(config)
			return
		}

		// Atualizar configuração
		if confluenceConfigURL != "" {
			config.URL = confluenceConfigURL
		}
		if confluenceConfigEmail != "" {
			config.Email = confluenceConfigEmail
		}
		if confluenceConfigAPIToken != "" {
			config.APIToken = confluenceConfigAPIToken
		}
		if confluenceConfigSpace != "" {
			config.Space = confluenceConfigSpace
		}

		if err := confluence.SaveConfig(config); err != nil {
			fmt.Printf("Erro ao salvar configuração: %v\n", err)
			return
		}

		fmt.Println("✓ Configuração do Confluence salva com sucesso!")
		showConfluenceConfig()
	},
}

var confluenceCreatePageCmd = &cobra.Command{
	Use:   "create-page",
	Short: "Criar uma página no Confluence",
	Long: `Cria uma página no Confluence.

Exemplos:
  snip confluence create-page --title "Documentação PostgreSQL" --content "# Título\n\nConteúdo..."
  snip confluence create-page -t "Configurações" -c "Conteúdo" -p "123456"  # Com página pai`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := confluence.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Space == "" {
			fmt.Println("❌ Configuração do Confluence incompleta. Execute: snip confluence config")
			return
		}

		if confluencePageTitle == "" {
			fmt.Println("❌ Título da página é obrigatório (use --title)")
			return
		}

		client, err := confluence.NewClient(config)
		if err != nil {
			fmt.Printf("Erro ao criar cliente Confluence: %v\n", err)
			return
		}

		content := confluencePageContent
		if content == "" {
			content = "Página criada automaticamente pelo SnipAI"
		}

		// Converter markdown para HTML básico do Confluence
		content = convertMarkdownToConfluence(content)

		page, err := client.CreatePage(confluencePageTitle, content, confluencePageParentID)
		if err != nil {
			fmt.Printf("Erro ao criar página: %v\n", err)
			return
		}

		fmt.Printf("✅ Página criada com sucesso!\n")
		fmt.Printf("  ID: %s\n", page.ID)
		fmt.Printf("  Título: %s\n", page.Title)
		fmt.Printf("  URL: %s/wiki%s\n", config.URL, page.ID)
	},
}

func showConfluenceConfig() {
	config, err := confluence.LoadConfig()
	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		return
	}

	fmt.Println("📋 Configuração Atual do Confluence:")
	fmt.Printf("  URL: %s\n", config.URL)
	fmt.Printf("  Email: %s\n", config.Email)
	tokenDisplay := config.APIToken
	if len(tokenDisplay) > 8 {
		tokenDisplay = tokenDisplay[:4] + "..." + tokenDisplay[len(tokenDisplay)-4:]
	}
	fmt.Printf("  API Token: %s\n", tokenDisplay)
	fmt.Printf("  Espaço: %s\n", config.Space)
}

func interactiveConfluenceConfig(config *confluence.ConfluenceConfig) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Configuração Interativa do Confluence ---")

	fmt.Printf("URL do Confluence atual (%s): ", config.URL)
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

	fmt.Printf("Chave do Espaço atual (%s): ", config.Space)
	spaceInput, _ := reader.ReadString('\n')
	spaceInput = strings.TrimSpace(spaceInput)
	if spaceInput != "" {
		config.Space = spaceInput
	}

	if err := confluence.SaveConfig(config); err != nil {
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}
	fmt.Println("✓ Configuração salva com sucesso!")
	showConfluenceConfig()
}

// convertMarkdownToConfluence converte markdown básico para formato do Confluence
func convertMarkdownToConfluence(markdown string) string {
	// Conversão básica - pode ser expandida
	html := markdown
	html = strings.ReplaceAll(html, "# ", "<h1>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = strings.ReplaceAll(html, "**", "<strong>")
	html = strings.ReplaceAll(html, "*", "<em>")
	html = strings.ReplaceAll(html, "`", "<code>")
	
	// Formato do Confluence Storage Format
	return fmt.Sprintf(`<p>%s</p>`, html)
}

