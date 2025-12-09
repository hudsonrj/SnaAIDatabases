package loganalyzer

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/snip/internal/dbtypes"
)

// LogAnalyzer analisa logs de banco de dados
type LogAnalyzer struct {
	DatabaseType dbtypes.DatabaseType
}

// NewLogAnalyzer cria um novo analisador de logs
func NewLogAnalyzer(dbType dbtypes.DatabaseType) *LogAnalyzer {
	return &LogAnalyzer{DatabaseType: dbType}
}

// AnalyzeLog analisa um arquivo de log
func (la *LogAnalyzer) AnalyzeLog(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	// Detecta o tipo de arquivo pela extensão
	if strings.HasSuffix(strings.ToLower(filePath), ".xml") {
		return la.analyzeXMLLog(file)
	}
	return la.analyzeTextLog(file)
}

// analyzeTextLog analisa logs em formato texto
func (la *LogAnalyzer) analyzeTextLog(file io.Reader) (string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	text := string(content)
	var analysis strings.Builder

	analysis.WriteString("# Análise de Log de Texto\n\n")
	analysis.WriteString(fmt.Sprintf("**Tipo de Banco:** %s\n\n", la.DatabaseType))

	// Análise específica por tipo de banco
	switch la.DatabaseType {
	case dbtypes.DatabaseTypeOracle:
		analysis.WriteString(la.analyzeOracleLog(text))
	case dbtypes.DatabaseTypeSQLServer:
		analysis.WriteString(la.analyzeSQLServerLog(text))
	case dbtypes.DatabaseTypeMySQL:
		analysis.WriteString(la.analyzeMySQLLog(text))
	case dbtypes.DatabaseTypePostgreSQL:
		analysis.WriteString(la.analyzePostgreSQLLog(text))
	case dbtypes.DatabaseTypeMongoDB:
		analysis.WriteString(la.analyzeMongoDBLog(text))
	default:
		analysis.WriteString(la.analyzeGenericLog(text))
	}

	return analysis.String(), nil
}

// analyzeXMLLog analisa logs em formato XML
func (la *LogAnalyzer) analyzeXMLLog(file io.Reader) (string, error) {
	decoder := xml.NewDecoder(file)
	var analysis strings.Builder

	analysis.WriteString("# Análise de Log XML\n\n")
	analysis.WriteString(fmt.Sprintf("**Tipo de Banco:** %s\n\n", la.DatabaseType))

	var inElement string
	var errorCount, warningCount, infoCount int
	var errors []string
	var warnings []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("erro ao decodificar XML: %w", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "error" || inElement == "Error" || inElement == "ERROR" {
				errorCount++
			} else if inElement == "warning" || inElement == "Warning" || inElement == "WARNING" {
				warningCount++
			} else if inElement == "info" || inElement == "Info" || inElement == "INFO" {
				infoCount++
			}
		case xml.CharData:
			if inElement == "error" || inElement == "Error" || inElement == "ERROR" {
				errors = append(errors, string(se))
			} else if inElement == "warning" || inElement == "Warning" || inElement == "WARNING" {
				warnings = append(warnings, string(se))
			}
		}
	}

	analysis.WriteString(fmt.Sprintf("## Estatísticas\n\n"))
	analysis.WriteString(fmt.Sprintf("- **Erros:** %d\n", errorCount))
	analysis.WriteString(fmt.Sprintf("- **Avisos:** %d\n", warningCount))
	analysis.WriteString(fmt.Sprintf("- **Informações:** %d\n\n", infoCount))

	if len(errors) > 0 {
		analysis.WriteString("## Erros Encontrados\n\n")
		for i, err := range errors {
			if i < 10 { // Limita a 10 erros
				analysis.WriteString(fmt.Sprintf("%d. %s\n", i+1, err))
			}
		}
		if len(errors) > 10 {
			analysis.WriteString(fmt.Sprintf("\n... e mais %d erros\n", len(errors)-10))
		}
	}

	if len(warnings) > 0 {
		analysis.WriteString("\n## Avisos Encontrados\n\n")
		for i, warn := range warnings {
			if i < 10 {
				analysis.WriteString(fmt.Sprintf("%d. %s\n", i+1, warn))
			}
		}
		if len(warnings) > 10 {
			analysis.WriteString(fmt.Sprintf("\n... e mais %d avisos\n", len(warnings)-10))
		}
	}

	return analysis.String(), nil
}

// analyzeOracleLog analisa logs específicos do Oracle
func (la *LogAnalyzer) analyzeOracleLog(text string) string {
	var analysis strings.Builder

	// Padrões comuns do Oracle
	errorPattern := regexp.MustCompile(`(?i)ORA-\d+`)
	errors := errorPattern.FindAllString(text, -1)

	analysis.WriteString("## Análise Oracle\n\n")
	analysis.WriteString(fmt.Sprintf("**Erros ORA encontrados:** %d\n\n", len(errors)))

	if len(errors) > 0 {
		errorMap := make(map[string]int)
		for _, err := range errors {
			errorMap[err]++
		}

		analysis.WriteString("### Erros mais frequentes:\n\n")
		for err, count := range errorMap {
			analysis.WriteString(fmt.Sprintf("- %s: %d ocorrências\n", err, count))
		}
	}

	// Detectar problemas de tablespace
	if strings.Contains(strings.ToLower(text), "tablespace") {
		analysis.WriteString("\n### Problemas de Tablespace detectados\n")
	}

	// Detectar problemas de conexão
	if strings.Contains(strings.ToLower(text), "connection") || strings.Contains(strings.ToLower(text), "connect") {
		analysis.WriteString("\n### Problemas de conexão detectados\n")
	}

	return analysis.String()
}

// analyzeSQLServerLog analisa logs específicos do SQL Server
func (la *LogAnalyzer) analyzeSQLServerLog(text string) string {
	var analysis strings.Builder

	analysis.WriteString("## Análise SQL Server\n\n")

	// Detectar erros comuns
	errorPattern := regexp.MustCompile(`(?i)error\s+\d+`)
	errors := errorPattern.FindAllString(text, -1)
	analysis.WriteString(fmt.Sprintf("**Erros encontrados:** %d\n\n", len(errors)))

	// Detectar deadlocks
	if strings.Contains(strings.ToLower(text), "deadlock") {
		analysis.WriteString("### Deadlocks detectados\n")
	}

	// Detectar problemas de bloqueio
	if strings.Contains(strings.ToLower(text), "blocked") {
		analysis.WriteString("### Bloqueios detectados\n")
	}

	return analysis.String()
}

// analyzeMySQLLog analisa logs específicos do MySQL
func (la *LogAnalyzer) analyzeMySQLLog(text string) string {
	var analysis strings.Builder

	analysis.WriteString("## Análise MySQL\n\n")

	// Detectar erros
	errorPattern := regexp.MustCompile(`(?i)\[ERROR\]`)
	errors := errorPattern.FindAllString(text, -1)
	analysis.WriteString(fmt.Sprintf("**Erros encontrados:** %d\n\n", len(errors)))

	// Detectar slow queries
	if strings.Contains(strings.ToLower(text), "slow query") {
		analysis.WriteString("### Slow queries detectadas\n")
	}

	return analysis.String()
}

// analyzePostgreSQLLog analisa logs específicos do PostgreSQL
func (la *LogAnalyzer) analyzePostgreSQLLog(text string) string {
	var analysis strings.Builder

	analysis.WriteString("## Análise PostgreSQL\n\n")

	// Detectar erros
	errorPattern := regexp.MustCompile(`(?i)ERROR:\s+`)
	errors := errorPattern.FindAllString(text, -1)
	analysis.WriteString(fmt.Sprintf("**Erros encontrados:** %d\n\n", len(errors)))

	// Detectar problemas de conexão
	if strings.Contains(strings.ToLower(text), "connection") {
		analysis.WriteString("### Problemas de conexão detectados\n")
	}

	return analysis.String()
}

// analyzeMongoDBLog analisa logs específicos do MongoDB
func (la *LogAnalyzer) analyzeMongoDBLog(text string) string {
	var analysis strings.Builder

	analysis.WriteString("## Análise MongoDB\n\n")

	// Detectar erros
	errorPattern := regexp.MustCompile(`(?i)\"E\"`)
	errors := errorPattern.FindAllString(text, -1)
	analysis.WriteString(fmt.Sprintf("**Erros encontrados:** %d\n\n", len(errors)))

	return analysis.String()
}

// analyzeGenericLog analisa logs genéricos
func (la *LogAnalyzer) analyzeGenericLog(text string) string {
	var analysis strings.Builder

	analysis.WriteString("## Análise Genérica\n\n")

	lines := strings.Split(text, "\n")
	analysis.WriteString(fmt.Sprintf("**Total de linhas:** %d\n\n", len(lines)))

	// Contar erros, avisos, etc.
	errorCount := strings.Count(strings.ToLower(text), "error")
	warningCount := strings.Count(strings.ToLower(text), "warning")
	infoCount := strings.Count(strings.ToLower(text), "info")

	analysis.WriteString(fmt.Sprintf("- **Erros:** %d\n", errorCount))
	analysis.WriteString(fmt.Sprintf("- **Avisos:** %d\n", warningCount))
	analysis.WriteString(fmt.Sprintf("- **Informações:** %d\n", infoCount))

	return analysis.String()
}

// GetLogSummary retorna um resumo do log
func (la *LogAnalyzer) GetLogSummary(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	summary := make(map[string]interface{})
	summary["file_path"] = filePath
	summary["file_size"] = info.Size()
	summary["modified_time"] = info.ModTime().Format(time.RFC3339)
	summary["database_type"] = string(la.DatabaseType)

	return summary, nil
}
