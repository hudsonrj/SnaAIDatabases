package formatter

import (
	"fmt"
	"strings"
	"time"
)

// FormatAnalysisResult formata o resultado de uma análise de forma bonita
func FormatAnalysisResult(title string, dbType string, analysisType string, result string, aiInsights string, chart string, createdAt time.Time) string {
	var output strings.Builder

	// Cabeçalho
	output.WriteString(fmt.Sprintf("# %s\n\n", title))
	output.WriteString(fmt.Sprintf("**Tipo de Banco:** %s  \n", dbType))
	output.WriteString(fmt.Sprintf("**Tipo de Análise:** %s  \n", analysisType))
	output.WriteString(fmt.Sprintf("**Data:** %s  \n\n", createdAt.Format("2006-01-02 15:04:05")))
	output.WriteString("---\n\n")

	// Resultado principal
	if result != "" {
		output.WriteString("## 📊 Resultado da Análise\n\n")
		// Melhorar formatação do resultado
		formattedResult := improveMarkdownFormatting(result)
		output.WriteString(formattedResult)
		output.WriteString("\n\n")
	}

	// Gráfico
	if chart != "" {
		output.WriteString("## 📈 Visualização\n\n")
		output.WriteString(chart)
		output.WriteString("\n\n")
	}

	// Insights da IA
	if aiInsights != "" {
		output.WriteString("## 🤖 Insights da IA\n\n")
		formattedInsights := improveMarkdownFormatting(aiInsights)
		output.WriteString(formattedInsights)
		output.WriteString("\n\n")
	}

	// Rodapé
	output.WriteString("---\n\n")
	output.WriteString(fmt.Sprintf("*Relatório gerado em %s*\n", time.Now().Format("2006-01-02 15:04:05")))

	return output.String()
}

// improveMarkdownFormatting melhora a formatação markdown
func improveMarkdownFormatting(text string) string {
	// Garantir que títulos tenham espaço após #
	text = strings.ReplaceAll(text, "\n#", "\n\n#")
	text = strings.ReplaceAll(text, "\n##", "\n\n##")
	text = strings.ReplaceAll(text, "\n###", "\n\n###")

	// Melhorar formatação de listas
	lines := strings.Split(text, "\n")
	var formattedLines []string
	inList := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Detectar início de lista
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || 
		   strings.HasPrefix(trimmed, "• ") || (len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' && strings.Contains(trimmed, ".")) {
			if !inList && i > 0 && strings.TrimSpace(lines[i-1]) != "" {
				formattedLines = append(formattedLines, "")
			}
			inList = true
		} else if trimmed == "" {
			inList = false
		} else if inList {
			// Fim da lista
			inList = false
			if i > 0 && strings.TrimSpace(lines[i-1]) != "" {
				formattedLines = append(formattedLines, "")
			}
		}

		formattedLines = append(formattedLines, line)
	}

	// Remover linhas vazias duplicadas
	var result []string
	prevEmpty := false
	for _, line := range formattedLines {
		isEmpty := strings.TrimSpace(line) == ""
		if !(isEmpty && prevEmpty) {
			result = append(result, line)
		}
		prevEmpty = isEmpty
	}

	return strings.Join(result, "\n")
}

// FormatTable cria uma tabela markdown formatada
func FormatTable(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	var output strings.Builder

	// Cabeçalho
	output.WriteString("|")
	for _, header := range headers {
		output.WriteString(fmt.Sprintf(" %s |", header))
	}
	output.WriteString("\n")

	// Separador
	output.WriteString("|")
	for range headers {
		output.WriteString(" --- |")
	}
	output.WriteString("\n")

	// Linhas
	for _, row := range rows {
		output.WriteString("|")
		for i, cell := range row {
			if i < len(headers) {
				output.WriteString(fmt.Sprintf(" %s |", cell))
			}
		}
		// Preencher células faltantes
		for i := len(row); i < len(headers); i++ {
			output.WriteString(" |")
		}
		output.WriteString("\n")
	}

	return output.String()
}

// FormatCodeBlock formata um bloco de código
func FormatCodeBlock(code string, language string) string {
	return fmt.Sprintf("```%s\n%s\n```", language, code)
}

// FormatAlert cria um alerta formatado
func FormatAlert(message string, alertType string) string {
	icons := map[string]string{
		"info":    "ℹ️",
		"success": "✅",
		"warning": "⚠️",
		"error":   "❌",
	}

	icon := icons[alertType]
	if icon == "" {
		icon = "ℹ️"
	}

	return fmt.Sprintf("> %s **%s:** %s", icon, strings.ToUpper(alertType), message)
}

