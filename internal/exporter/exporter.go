package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/snip/internal/formatter"
)

// ExportToMarkdown exporta uma análise para arquivo markdown
func ExportToMarkdown(title string, dbType string, analysisType string, result string, aiInsights string, chart string, createdAt time.Time, filename string) (string, error) {
	// Se não fornecido, gerar nome baseado no título
	if filename == "" {
		// Limpar título para nome de arquivo
		safeTitle := strings.ToLower(title)
		safeTitle = strings.ReplaceAll(safeTitle, " ", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "/", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "\\", "_")
		safeTitle = strings.ReplaceAll(safeTitle, ":", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "*", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "?", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "\"", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "<", "_")
		safeTitle = strings.ReplaceAll(safeTitle, ">", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "|", "_")
		
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s.md", safeTitle, timestamp)
	}

	// Garantir extensão .md
	if !strings.HasSuffix(strings.ToLower(filename), ".md") {
		filename = filename + ".md"
	}

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}

	// Criar diretório de exports se não existir
	exportDir := filepath.Join(homeDir, ".snip", "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de exports: %w", err)
	}

	// Caminho completo do arquivo
	filePath := filepath.Join(exportDir, filename)

	// Formatar conteúdo
	content := formatter.FormatAnalysisResult(title, dbType, analysisType, result, aiInsights, chart, createdAt)

	// Escrever arquivo
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return filePath, nil
}

// GetExportPath retorna o caminho do diretório de exports
func GetExportPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	exportDir := filepath.Join(homeDir, ".snip", "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", err
	}
	return exportDir, nil
}

