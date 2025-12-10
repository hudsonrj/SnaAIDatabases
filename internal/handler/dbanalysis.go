package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/snip/internal/dbanalysis"
	"github.com/snip/internal/exporter"
	"github.com/snip/internal/repository"
)

type DBAnalysisHandler interface {
	CreateAnalysis(title string, dbType dbanalysis.DatabaseType, analysisType dbanalysis.AnalysisType,
		outputType dbanalysis.OutputType, config *dbanalysis.ConnectionConfig, logFilePath string) error
	ListAnalyses(limit int, dbType dbanalysis.DatabaseType, analysisType dbanalysis.AnalysisType) error
	GetAnalysis(idStr string, verbose bool) error
	GetAnalysisByID(id int) (*dbanalysis.DBAnalysis, error)
	DeleteAnalysis(idStr string) error
	RunAnalysis(idStr string, exportFilename string) error
	ExportAnalysisToMarkdown(idStr string, filename string) (string, error)
}

type dbAnalysisHandler struct {
	analysisRepo repository.DBAnalysisRepository
	analyzer     *dbanalysis.Analyzer
}

func NewDBAnalysisHandler(analysisRepo repository.DBAnalysisRepository) (DBAnalysisHandler, error) {
	analyzer, err := dbanalysis.NewAnalyzer()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar analisador: %w", err)
	}

	return &dbAnalysisHandler{
		analysisRepo: analysisRepo,
		analyzer:     analyzer,
	}, nil
}

func (h *dbAnalysisHandler) CreateAnalysis(title string, dbType dbanalysis.DatabaseType,
	analysisType dbanalysis.AnalysisType, outputType dbanalysis.OutputType,
	config *dbanalysis.ConnectionConfig, logFilePath string) error {

	analysis := dbanalysis.NewDBAnalysis(title, dbType, analysisType, outputType)
	analysis.LogFilePath = logFilePath

	// Serializar configuração de conexão
	configJSON, err := dbanalysis.SerializeConnectionConfig(config)
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}
	analysis.ConnectionConfig = configJSON

	// Criar análise no banco
	if err := h.analysisRepo.Create(analysis); err != nil {
		return fmt.Errorf("erro ao criar análise: %w", err)
	}

	fmt.Printf("✓ Análise criada com sucesso!\n")
	fmt.Printf("  ID: #%d\n", analysis.ID)
	fmt.Printf("  Título: %s\n", analysis.Title)
	fmt.Printf("  Tipo de Banco: %s\n", analysis.DatabaseType)
	fmt.Printf("  Tipo de Análise: %s\n", analysis.AnalysisType)
	fmt.Printf("  Status: %s\n", analysis.Status)

	return nil
}

func (h *dbAnalysisHandler) ListAnalyses(limit int, dbType dbanalysis.DatabaseType,
	analysisType dbanalysis.AnalysisType) error {

	analyses, err := h.analysisRepo.GetAll(limit, dbType, analysisType)
	if err != nil {
		return fmt.Errorf("erro ao buscar análises: %w", err)
	}

	if len(analyses) == 0 {
		fmt.Println("Nenhuma análise encontrada.")
		return nil
	}

	fmt.Printf("Encontradas %d análise(s):\n\n", len(analyses))

	for _, analysis := range analyses {
		statusIcon := "⏳"
		if analysis.Status == "completed" {
			statusIcon = "✅"
		} else if analysis.Status == "error" {
			statusIcon = "❌"
		}

		fmt.Printf("%s #%d %s\n", statusIcon, analysis.ID, analysis.Title)
		fmt.Printf("  └── Banco: %s | Tipo: %s | Status: %s\n",
			analysis.DatabaseType, analysis.AnalysisType, analysis.Status)
		fmt.Printf("  └── Criado: %s\n", analysis.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

func (h *dbAnalysisHandler) GetAnalysis(idStr string, verbose bool) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID de análise inválido: %s", idStr)
	}

	analysis, err := h.analysisRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar análise: %w", err)
	}

	statusIcon := "⏳"
	if analysis.Status == "completed" {
		statusIcon = "✅"
	} else if analysis.Status == "error" {
		statusIcon = "❌"
	}

	fmt.Printf("%s #%d %s\n", statusIcon, analysis.ID, analysis.Title)
	fmt.Printf("  └── Tipo de Banco: %s\n", analysis.DatabaseType)
	fmt.Printf("  └── Tipo de Análise: %s\n", analysis.AnalysisType)
	fmt.Printf("  └── Formato de Saída: %s\n", analysis.OutputType)
	fmt.Printf("  └── Status: %s\n", analysis.Status)

	if analysis.LogFilePath != "" {
		fmt.Printf("  └── Arquivo de Log: %s\n", analysis.LogFilePath)
	}

	if verbose {
		// Deserializar configuração
		config, err := dbanalysis.DeserializeConnectionConfig(analysis.ConnectionConfig)
		if err == nil {
			fmt.Printf("  └── Host: %s\n", config.Host)
			fmt.Printf("  └── Database: %s\n", config.Database)
			fmt.Printf("  └── Remoto: %v\n", config.IsRemote)
		}

		fmt.Printf("  └── Criado: %s\n", analysis.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  └── Atualizado: %s\n", analysis.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// Extrair gráfico do resultado se houver
	chart := ""
	result := analysis.Result
	if strings.Contains(result, "## Visualização") {
		parts := strings.Split(result, "## Visualização")
		if len(parts) > 1 {
			chart = strings.TrimSpace(parts[1])
			// Remover gráfico do resultado principal
			result = strings.TrimSpace(strings.Split(result, "## Visualização")[0])
		}
	}

	if result != "" {
		fmt.Println("\n" + strings.Repeat("═", 70))
		fmt.Println("📊 RESULTADO DA ANÁLISE")
		fmt.Println(strings.Repeat("═", 70) + "\n")
		fmt.Println(result)
		fmt.Println()
	}

	if chart != "" {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("📈 VISUALIZAÇÃO")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		fmt.Println(chart)
		fmt.Println()
	}

	if analysis.AIInsights != "" {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("🤖 INSIGHTS DA IA")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		fmt.Println(analysis.AIInsights)
		fmt.Println()
	}

	if analysis.ErrorMessage != "" {
		fmt.Printf("\n❌ Erro: %s\n", analysis.ErrorMessage)
	}

	return nil
}

func (h *dbAnalysisHandler) DeleteAnalysis(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID de análise inválido: %s", idStr)
	}

	if err := h.analysisRepo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar análise: %w", err)
	}

	fmt.Printf("✓ Análise #%d deletada com sucesso!\n", id)
	return nil
}

func (h *dbAnalysisHandler) RunAnalysis(idStr string, exportFilename string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID de análise inválido: %s", idStr)
	}

	analysis, err := h.analysisRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar análise: %w", err)
	}

	fmt.Printf("Executando análise #%d: %s\n", analysis.ID, analysis.Title)
	fmt.Println("Aguarde...\n")

	// Deserializar configuração
	config, err := dbanalysis.DeserializeConnectionConfig(analysis.ConnectionConfig)
	if err != nil {
		return fmt.Errorf("erro ao deserializar configuração: %w", err)
	}

	// Executar análise
	if err := h.analyzer.PerformAnalysis(analysis, config); err != nil {
		analysis.Status = "error"
		analysis.ErrorMessage = err.Error()
		h.analysisRepo.Update(analysis)
		return fmt.Errorf("erro ao executar análise: %w", err)
	}

	// Atualizar análise
	if err := h.analysisRepo.Update(analysis); err != nil {
		return fmt.Errorf("erro ao atualizar análise: %w", err)
	}

	fmt.Printf("✓ Análise concluída com sucesso!\n")
	fmt.Printf("  Status: %s\n", analysis.Status)

	// Extrair gráfico do resultado se houver
	chart := ""
	result := analysis.Result
	if strings.Contains(result, "## Visualização") {
		parts := strings.Split(result, "## Visualização")
		if len(parts) > 1 {
			chart = parts[1]
			// Remover gráfico do resultado principal
			result = strings.Split(result, "## Visualização")[0]
		}
	}

	// Formatar e exibir resultado melhorado
	if result != "" {
		fmt.Println("\n" + strings.Repeat("═", 70))
		fmt.Println("📊 RESULTADO DA ANÁLISE")
		fmt.Println(strings.Repeat("═", 70) + "\n")
		fmt.Println(result)
		fmt.Println()
	}

	if chart != "" {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("📈 VISUALIZAÇÃO")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		fmt.Println(chart)
		fmt.Println()
	}

	if analysis.AIInsights != "" {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("🤖 INSIGHTS DA IA")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		fmt.Println(analysis.AIInsights)
		fmt.Println()
	}

	// Exportar para markdown se solicitado
	if exportFilename != "" {
		filePath, err := h.ExportAnalysisToMarkdown(idStr, exportFilename)
		if err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao exportar: %v\n", err)
		} else {
			fmt.Printf("✅ Relatório exportado para: %s\n", filePath)
		}
	}

	return nil
}

// GetAnalysisByID retorna uma análise por ID
func (h *dbAnalysisHandler) GetAnalysisByID(id int) (*dbanalysis.DBAnalysis, error) {
	analysis, err := h.analysisRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar análise: %w", err)
	}
	return analysis, nil
}

// ExportAnalysisToMarkdown exporta uma análise para markdown
func (h *dbAnalysisHandler) ExportAnalysisToMarkdown(idStr string, filename string) (string, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return "", fmt.Errorf("ID de análise inválido: %s", idStr)
	}

	analysis, err := h.analysisRepo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar análise: %w", err)
	}

	// Extrair gráfico do resultado se houver
	chart := ""
	result := analysis.Result
	if strings.Contains(result, "## Visualização") {
		parts := strings.Split(result, "## Visualização")
		if len(parts) > 1 {
			chart = strings.TrimSpace(parts[1])
			// Remover gráfico do resultado principal
			result = strings.TrimSpace(strings.Split(result, "## Visualização")[0])
		}
	}

	filePath, err := exporter.ExportToMarkdown(
		analysis.Title,
		string(analysis.DatabaseType),
		string(analysis.AnalysisType),
		result,
		analysis.AIInsights,
		chart,
		analysis.CreatedAt,
		filename,
	)

	return filePath, err
}

// FormatOutput formata a saída conforme o tipo especificado
func FormatOutput(content string, outputType dbanalysis.OutputType) string {
	switch outputType {
	case dbanalysis.OutputTypeJSON:
		// Converter para JSON se necessário
		data := map[string]string{"content": content}
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			return string(jsonBytes)
		}
		return content
	case dbanalysis.OutputTypeMarkdown:
		return content
	case dbanalysis.OutputTypeHTML:
		// Converter markdown para HTML (requer biblioteca adicional)
		return fmt.Sprintf("<html><body><pre>%s</pre></body></html>", content)
	case dbanalysis.OutputTypeText:
		// Remover formatação markdown básica
		text := strings.ReplaceAll(content, "#", "")
		text = strings.ReplaceAll(text, "**", "")
		text = strings.ReplaceAll(text, "*", "")
		return text
	default:
		return content
	}
}

