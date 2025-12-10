package checklist

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ChecklistType representa o tipo de checklist
type ChecklistType string

const (
	ChecklistTypeGeneric     ChecklistType = "generic"
	ChecklistTypeDaily       ChecklistType = "daily"
	ChecklistTypeWeekly      ChecklistType = "weekly"
	ChecklistTypeDeep        ChecklistType = "deep"
	ChecklistTypeBackup      ChecklistType = "backup"
	ChecklistTypeSecurity    ChecklistType = "security"
	ChecklistTypePerformance ChecklistType = "performance"
	ChecklistTypeMaintenance ChecklistType = "maintenance"
)

// BulkChecklistItem representa um item de checklist para processamento em massa
type BulkChecklistItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

// BulkChecklistResult representa o resultado de um checklist em massa
type BulkChecklistResult struct {
	TotalItems    int                `json:"total_items"`
	Completed     int                `json:"completed"`
	Pending       int                `json:"pending"`
	Failed        int                `json:"failed"`
	Items         []BulkChecklistItem `json:"items"`
	GeneratedAt   time.Time          `json:"generated_at"`
	ExecutionTime time.Duration       `json:"execution_time"`
}

// ProcessBulkChecklistFromCSV processa um checklist em massa a partir de um arquivo CSV
func ProcessBulkChecklistFromCSV(csvPath string) (*BulkChecklistResult, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV deve ter pelo menos um cabeçalho e uma linha de dados")
	}

	// Parsear cabeçalho
	header := records[0]
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Verificar campos obrigatórios
	requiredFields := []string{"title", "description", "category"}
	for _, field := range requiredFields {
		if _, ok := headerMap[field]; !ok {
			return nil, fmt.Errorf("campo obrigatório '%s' não encontrado no CSV", field)
		}
	}

	startTime := time.Now()
	result := &BulkChecklistResult{
		Items:       []BulkChecklistItem{},
		GeneratedAt: time.Now(),
	}

	// Processar linhas
	for i, record := range records[1:] {
		if len(record) < len(header) {
			continue // Pular linhas incompletas
		}

		item := BulkChecklistItem{
			Title:       getField(record, headerMap, "title"),
			Description: getField(record, headerMap, "description"),
			Category:    getField(record, headerMap, "category"),
			Priority:    getField(record, headerMap, "priority"),
			Status:      getField(record, headerMap, "status"),
			Notes:       getField(record, headerMap, "notes"),
		}

		// Valores padrão
		if item.Priority == "" {
			item.Priority = "medium"
		}
		if item.Status == "" {
			item.Status = "pending"
		}

		// Processar item (aqui você pode adicionar lógica de validação/execução)
		result.Items = append(result.Items, item)
		result.TotalItems++

		switch strings.ToLower(item.Status) {
		case "completed", "done", "ok":
			result.Completed++
		case "failed", "error":
			result.Failed++
		default:
			result.Pending++
		}

		_ = i // Evitar erro de variável não usada
	}

	result.ExecutionTime = time.Since(startTime)

	return result, nil
}

// getField obtém um campo do record baseado no headerMap
func getField(record []string, headerMap map[string]int, fieldName string) string {
	if idx, ok := headerMap[strings.ToLower(fieldName)]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx])
	}
	return ""
}

// ExportBulkChecklistToMarkdown exporta o resultado do checklist em massa para markdown
func ExportBulkChecklistToMarkdown(result *BulkChecklistResult, filename string) (string, error) {
	if filename == "" {
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("bulk_checklist_%s.md", timestamp)
	}

	if !strings.HasSuffix(strings.ToLower(filename), ".md") {
		filename = filename + ".md"
	}

	// Criar conteúdo markdown
	var content strings.Builder

	content.WriteString(fmt.Sprintf("# Checklist em Massa - Resultados\n\n"))
	content.WriteString(fmt.Sprintf("**Data de Geração:** %s\n", result.GeneratedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("**Tempo de Execução:** %s\n\n", result.ExecutionTime.Round(time.Second)))
	content.WriteString("---\n\n")

	// Estatísticas
	content.WriteString("## 📊 Estatísticas\n\n")
	content.WriteString(fmt.Sprintf("| Métrica | Valor |\n"))
	content.WriteString(fmt.Sprintf("|---------|-------|\n"))
	content.WriteString(fmt.Sprintf("| Total de Itens | %d |\n", result.TotalItems))
	content.WriteString(fmt.Sprintf("| ✅ Concluídos | %d |\n", result.Completed))
	content.WriteString(fmt.Sprintf("| ⏳ Pendentes | %d |\n", result.Pending))
	content.WriteString(fmt.Sprintf("| ❌ Falhas | %d |\n", result.Failed))
	content.WriteString("\n")

	// Agrupar por categoria
	categories := make(map[string][]BulkChecklistItem)
	for _, item := range result.Items {
		category := item.Category
		if category == "" {
			category = "Sem Categoria"
		}
		categories[category] = append(categories[category], item)
	}

	// Itens por categoria
	content.WriteString("## 📋 Itens por Categoria\n\n")
	for category, items := range categories {
		content.WriteString(fmt.Sprintf("### %s\n\n", category))
		content.WriteString("| Título | Prioridade | Status | Notas |\n")
		content.WriteString("|--------|------------|--------|-------|\n")

		for _, item := range items {
			priorityIcon := getPriorityIcon(item.Priority)
			statusIcon := getStatusIcon(item.Status)
			notes := item.Notes
			if len(notes) > 50 {
				notes = notes[:50] + "..."
			}
			content.WriteString(fmt.Sprintf("| %s | %s %s | %s %s | %s |\n",
				item.Title, priorityIcon, item.Priority, statusIcon, item.Status, notes))
		}
		content.WriteString("\n")
	}

	// Detalhes completos
	content.WriteString("## 📝 Detalhes Completos\n\n")
	for i, item := range result.Items {
		content.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, item.Title))
		content.WriteString(fmt.Sprintf("**Categoria:** %s  \n", item.Category))
		content.WriteString(fmt.Sprintf("**Prioridade:** %s  \n", item.Priority))
		content.WriteString(fmt.Sprintf("**Status:** %s  \n", item.Status))
		content.WriteString(fmt.Sprintf("**Descrição:** %s  \n\n", item.Description))
		if item.Notes != "" {
			content.WriteString(fmt.Sprintf("**Notas:** %s  \n\n", item.Notes))
		}
		content.WriteString("---\n\n")
	}

	// Exportar arquivo
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}

	exportDir := filepath.Join(homeDir, ".snip", "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de exports: %w", err)
	}

	filePath := filepath.Join(exportDir, filename)
	if err := os.WriteFile(filePath, []byte(content.String()), 0644); err != nil {
		return "", fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return filePath, nil
}

func getPriorityIcon(priority string) string {
	switch strings.ToLower(priority) {
	case "high", "alta":
		return "🔴"
	case "medium", "media":
		return "🟡"
	case "low", "baixa":
		return "🟢"
	default:
		return "⚪"
	}
}

func getStatusIcon(status string) string {
	switch strings.ToLower(status) {
	case "completed", "done", "ok", "concluido":
		return "✅"
	case "pending", "pendente":
		return "⏳"
	case "failed", "error", "falha":
		return "❌"
	case "in_progress", "em_andamento":
		return "🔄"
	default:
		return "⚪"
	}
}

// GenerateCSVTemplate gera um template CSV para um tipo específico de checklist
func GenerateCSVTemplate(checklistType ChecklistType, outputPath string) error {
	var headers []string
	var sampleRows [][]string

	switch checklistType {
	case ChecklistTypeDaily:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "check_time", "result"}
		sampleRows = [][]string{
			{"Verificar conexões ativas", "Contar número de conexões ativas no banco", "Monitoramento", "high", "pending", "", "09:00", ""},
			{"Verificar espaço em disco", "Verificar espaço disponível em tablespaces", "Storage", "high", "pending", "", "09:00", ""},
			{"Verificar logs de erro", "Revisar logs de erro das últimas 24h", "Logs", "medium", "pending", "", "09:00", ""},
			{"Verificar backups", "Confirmar execução dos backups agendados", "Backup", "high", "pending", "", "10:00", ""},
		}

	case ChecklistTypeWeekly:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "week", "assigned_to"}
		sampleRows = [][]string{
			{"Análise de performance", "Revisar queries lentas e índices", "Performance", "high", "pending", "", "Semana 1", ""},
			{"Revisão de segurança", "Auditar usuários e permissões", "Security", "high", "pending", "", "Semana 1", ""},
			{"Atualização de estatísticas", "Executar ANALYZE/UPDATE STATISTICS", "Manutenção", "medium", "pending", "", "Semana 1", ""},
			{"Revisão de fragmentação", "Verificar fragmentação de tabelas", "Manutenção", "medium", "pending", "", "Semana 1", ""},
		}

	case ChecklistTypeDeep:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "impact", "effort", "risk_level"}
		sampleRows = [][]string{
			{"Auditoria completa de segurança", "Revisão completa de segurança do banco", "Security", "high", "pending", "", "Alto", "Alto", "Médio"},
			{"Análise de capacidade", "Projeção de crescimento e capacidade", "Capacidade", "high", "pending", "", "Alto", "Médio", "Baixo"},
			{"Otimização de índices", "Análise e otimização de todos os índices", "Performance", "medium", "pending", "", "Médio", "Alto", "Baixo"},
			{"Revisão de arquitetura", "Avaliar arquitetura e sugerir melhorias", "Arquitetura", "high", "pending", "", "Alto", "Alto", "Médio"},
		}

	case ChecklistTypeBackup:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "backup_type", "retention_days", "last_backup"}
		sampleRows = [][]string{
			{"Backup completo", "Verificar execução de backup completo", "Backup", "high", "pending", "", "Full", "30", ""},
			{"Backup incremental", "Verificar execução de backup incremental", "Backup", "high", "pending", "", "Incremental", "7", ""},
			{"Teste de restore", "Testar procedimento de restore", "Backup", "high", "pending", "", "Test", "", ""},
			{"Verificar retenção", "Confirmar políticas de retenção", "Backup", "medium", "pending", "", "Policy", "", ""},
		}

	case ChecklistTypeSecurity:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "severity", "compliance", "remediation"}
		sampleRows = [][]string{
			{"Auditar usuários", "Revisar usuários e remover inativos", "Security", "high", "pending", "", "Alta", "SOX", "Remover usuários inativos"},
			{"Revisar permissões", "Auditar permissões e privilégios", "Security", "high", "pending", "", "Alta", "PCI-DSS", "Aplicar least privilege"},
			{"Verificar criptografia", "Confirmar criptografia de dados sensíveis", "Security", "high", "pending", "", "Alta", "GDPR", "Habilitar TDE"},
			{"Auditar logs de acesso", "Revisar logs de acesso e autenticação", "Security", "medium", "pending", "", "Média", "SOX", "Configurar alertas"},
		}

	case ChecklistTypePerformance:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "metric", "threshold", "current_value"}
		sampleRows = [][]string{
			{"CPU Utilization", "Monitorar utilização de CPU", "Performance", "high", "pending", "", "CPU %", "80%", ""},
			{"Memory Usage", "Monitorar uso de memória", "Performance", "high", "pending", "", "Memory %", "85%", ""},
			{"Disk I/O", "Monitorar I/O de disco", "Performance", "medium", "pending", "", "IOPS", "1000", ""},
			{"Query Performance", "Identificar queries lentas", "Performance", "high", "pending", "", "Query Time", "5s", ""},
		}

	case ChecklistTypeMaintenance:
		headers = []string{"title", "description", "category", "priority", "status", "notes", "frequency", "last_execution", "next_execution"}
		sampleRows = [][]string{
			{"Vacuum/Analyze", "Executar VACUUM e ANALYZE", "Manutenção", "medium", "pending", "", "Semanal", "", ""},
			{"Reindex", "Reindexar tabelas fragmentadas", "Manutenção", "low", "pending", "", "Mensal", "", ""},
			{"Update Statistics", "Atualizar estatísticas do otimizador", "Manutenção", "medium", "pending", "", "Semanal", "", ""},
			{"Cleanup Logs", "Limpar logs antigos", "Manutenção", "low", "pending", "", "Mensal", "", ""},
		}

	default: // ChecklistTypeGeneric
		headers = []string{"title", "description", "category", "priority", "status", "notes"}
		sampleRows = [][]string{
			{"Item 1", "Descrição do item 1", "Categoria 1", "high", "pending", "Notas adicionais"},
			{"Item 2", "Descrição do item 2", "Categoria 2", "medium", "pending", ""},
			{"Item 3", "Descrição do item 3", "Categoria 1", "low", "completed", "Item concluído"},
		}
	}

	// Criar arquivo CSV
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escrever cabeçalho
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho: %w", err)
	}

	// Escrever linhas de exemplo
	for _, row := range sampleRows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("erro ao escrever linha: %w", err)
		}
	}

	return nil
}

// GetAvailableChecklistTypes retorna os tipos de checklist disponíveis
func GetAvailableChecklistTypes() []ChecklistType {
	return []ChecklistType{
		ChecklistTypeGeneric,
		ChecklistTypeDaily,
		ChecklistTypeWeekly,
		ChecklistTypeDeep,
		ChecklistTypeBackup,
		ChecklistTypeSecurity,
		ChecklistTypePerformance,
		ChecklistTypeMaintenance,
	}
}

// GetChecklistTypeDescription retorna a descrição de um tipo de checklist
func GetChecklistTypeDescription(checklistType ChecklistType) string {
	descriptions := map[ChecklistType]string{
		ChecklistTypeGeneric:     "Checklist genérico com campos básicos",
		ChecklistTypeDaily:       "Checklist diário para tarefas rotineiras",
		ChecklistTypeWeekly:      "Checklist semanal para revisões periódicas",
		ChecklistTypeDeep:        "Checklist profundo para análises detalhadas",
		ChecklistTypeBackup:      "Checklist específico para backups e restores",
		ChecklistTypeSecurity:    "Checklist de segurança e compliance",
		ChecklistTypePerformance: "Checklist de performance e monitoramento",
		ChecklistTypeMaintenance: "Checklist de manutenção preventiva",
	}
	return descriptions[checklistType]
}

