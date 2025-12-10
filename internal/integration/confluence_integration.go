package integration

import (
	"fmt"
	"strings"

	"github.com/snip/internal/confluence"
	"github.com/snip/internal/dbanalysis"
)

// ExportAnalysisToConfluence exporta uma análise para o Confluence
func ExportAnalysisToConfluence(analysis *dbanalysis.DBAnalysis, pageTitle string, parentPageID string) (string, error) {
	config, err := confluence.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar configuração Confluence: %w", err)
	}

	if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Space == "" {
		return "", fmt.Errorf("configuração do Confluence incompleta. Execute: snip confluence config")
	}

	client, err := confluence.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("erro ao criar cliente Confluence: %w", err)
	}

	// Se não fornecido, usar título da análise
	if pageTitle == "" {
		pageTitle = analysis.Title
	}

	// Criar conteúdo formatado
	content := formatAnalysisForConfluence(analysis)

	page, err := client.CreatePage(pageTitle, content, parentPageID)
	if err != nil {
		return "", fmt.Errorf("erro ao criar página: %w", err)
	}

	return page.ID, nil
}

// ExportDatabaseConfigToConfluence exporta configurações de banco de dados para Confluence
func ExportDatabaseConfigToConfluence(dbType, configData, pageTitle string, parentPageID string) (string, error) {
	config, err := confluence.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar configuração Confluence: %w", err)
	}

	if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Space == "" {
		return "", fmt.Errorf("configuração do Confluence incompleta. Execute: snip confluence config")
	}

	client, err := confluence.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("erro ao criar cliente Confluence: %w", err)
	}

	// Formatar conteúdo
	var content strings.Builder
	content.WriteString(fmt.Sprintf("<h1>Configurações do Banco de Dados: %s</h1>\n", dbType))
	content.WriteString(fmt.Sprintf("<p><strong>Data:</strong> %s</p>\n", "2024-01-01"))
	content.WriteString("<h2>Configurações</h2>\n")
	content.WriteString("<pre>")
	content.WriteString(configData)
	content.WriteString("</pre>")

	page, err := client.CreatePage(pageTitle, content.String(), parentPageID)
	if err != nil {
		return "", fmt.Errorf("erro ao criar página: %w", err)
	}

	return page.ID, nil
}

// formatAnalysisForConfluence formata uma análise para o formato do Confluence
func formatAnalysisForConfluence(analysis *dbanalysis.DBAnalysis) string {
	var content strings.Builder

	// Título
	content.WriteString(fmt.Sprintf("<h1>%s</h1>\n", analysis.Title))

	// Metadados
	content.WriteString("<table><tbody>\n")
	content.WriteString(fmt.Sprintf("<tr><td><strong>Tipo de Banco:</strong></td><td>%s</td></tr>\n", analysis.DatabaseType))
	content.WriteString(fmt.Sprintf("<tr><td><strong>Tipo de Análise:</strong></td><td>%s</td></tr>\n", analysis.AnalysisType))
	content.WriteString(fmt.Sprintf("<tr><td><strong>Status:</strong></td><td>%s</td></tr>\n", analysis.Status))
	content.WriteString(fmt.Sprintf("<tr><td><strong>Data:</strong></td><td>%s</td></tr>\n", analysis.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("</tbody></table>\n\n")

	// Resultado
	if analysis.Result != "" {
		content.WriteString("<h2>Resultado da Análise</h2>\n")
		// Converter markdown básico para HTML do Confluence
		resultHTML := convertMarkdownToConfluenceHTML(analysis.Result)
		content.WriteString(resultHTML)
		content.WriteString("\n\n")
	}

	// Insights da IA
	if analysis.AIInsights != "" {
		content.WriteString("<h2>Insights da IA</h2>\n")
		insightsHTML := convertMarkdownToConfluenceHTML(analysis.AIInsights)
		content.WriteString(insightsHTML)
		content.WriteString("\n\n")
	}

	return content.String()
}

// convertMarkdownToConfluenceHTML converte markdown básico para HTML do Confluence
func convertMarkdownToConfluenceHTML(markdown string) string {
	html := markdown
	
	// Títulos
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "# ", "<h1>")
	
	// Quebras de linha
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = "<p>" + html + "</p>"
	
	// Negrito
	html = strings.ReplaceAll(html, "**", "<strong>")
	html = strings.ReplaceAll(html, "*", "<em>")
	
	// Código
	html = strings.ReplaceAll(html, "`", "<code>")
	
	// Listas
	lines := strings.Split(html, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			result = append(result, "<li>"+strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")+"</li>")
		} else {
			result = append(result, line)
		}
	}
	html = strings.Join(result, "\n")
	
	return html
}

