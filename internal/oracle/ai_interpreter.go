package oracle

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/snip/internal/ai"
)

// AINaturalLanguageInterpreter interpreta requisições em linguagem natural
type AINaturalLanguageInterpreter struct {
	aiClient *ai.GroqClient
}

// NewAINaturalLanguageInterpreter cria um novo interpretador
func NewAINaturalLanguageInterpreter() (*AINaturalLanguageInterpreter, error) {
	aiClient, err := ai.NewGroqClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &AINaturalLanguageInterpreter{aiClient: aiClient}, nil
}

// ParseAWRRequest interpreta uma requisição em linguagem natural para AWR
func (i *AINaturalLanguageInterpreter) ParseAWRRequest(naturalLanguage string, availableSnapshots []Snapshot) (*AWRReportRequest, error) {
	prompt := fmt.Sprintf(`Você é um especialista em Oracle Database. Interprete a seguinte requisição do usuário e extraia os parâmetros para gerar um relatório AWR.

Requisição do usuário: "%s"

Snapshots disponíveis:
%s

Retorne APENAS um JSON com os seguintes campos (use null para valores não especificados):
{
  "begin_time": "YYYY-MM-DD HH:MM:SS" ou null,
  "end_time": "YYYY-MM-DD HH:MM:SS" ou null,
  "begin_snapshot": número ou null,
  "end_snapshot": número ou null,
  "duration_hours": número de horas ou null,
  "report_type": "html" ou "text"
}

Exemplos de interpretação:
- "relatório das últimas 2 horas" -> duration_hours: 2, end_time: agora
- "relatório de ontem" -> begin_time: início de ontem, end_time: fim de ontem
- "relatório entre os snaps 100 e 200" -> begin_snapshot: 100, end_snapshot: 200
- "relatório de hoje de manhã" -> begin_time: início de hoje, end_time: meio-dia

Retorne APENAS o JSON, sem explicações.`, naturalLanguage, formatSnapshots(availableSnapshots))

	response, err := i.aiClient.GenerateContent(prompt, 500)
	if err != nil {
		return nil, fmt.Errorf("erro ao interpretar requisição: %w", err)
	}

	// Parsear JSON da resposta
	req, err := parseAWRRequestFromJSON(response)
	if err != nil {
		// Tentar parsing manual se JSON falhar
		req = parseAWRRequestManual(naturalLanguage, availableSnapshots)
	}

	return req, nil
}

// ParseASHRequest interpreta uma requisição em linguagem natural para ASH
func (i *AINaturalLanguageInterpreter) ParseASHRequest(naturalLanguage string) (*ASHRequest, error) {
	prompt := fmt.Sprintf(`Você é um especialista em Oracle Database. Interprete a seguinte requisição do usuário e extraia os parâmetros para análise ASH.

Requisição do usuário: "%s"

Retorne APENAS um JSON com os seguintes campos (use null para valores não especificados):
{
  "sql_id": "string" ou null,
  "sid": número ou null,
  "serial": número ou null,
  "begin_time": "YYYY-MM-DD HH:MM:SS" ou null,
  "end_time": "YYYY-MM-DD HH:MM:SS" ou null,
  "duration_minutes": número de minutos ou null
}

Exemplos:
- "sessão 123" -> sid: 123
- "SQL ID abc123def" -> sql_id: "abc123def"
- "últimos 30 minutos" -> duration_minutes: 30
- "sessão 456 serial 789" -> sid: 456, serial: 789

Retorne APENAS o JSON, sem explicações.`, naturalLanguage)

	response, err := i.aiClient.GenerateContent(prompt, 500)
	if err != nil {
		return nil, fmt.Errorf("erro ao interpretar requisição: %w", err)
	}

	req, err := parseASHRequestFromJSON(response)
	if err != nil {
		req = parseASHRequestManual(naturalLanguage)
	}

	return req, nil
}

// InterpretAnalysisResults interpreta resultados de análise e fornece recomendações
func (i *AINaturalLanguageInterpreter) InterpretAnalysisResults(analysisType string, rawResults string) (string, error) {
	prompt := fmt.Sprintf(`Você é um DBA experiente em Oracle Database. Analise os seguintes resultados e forneça:

1. Resumo executivo em português
2. Principais problemas identificados
3. Recomendações de ação prioritárias
4. Explicação técnica simplificada

Tipo de Análise: %s

Resultados:
%s

Formate a resposta em markdown, sendo claro e objetivo. Use linguagem técnica mas acessível.`, analysisType, rawResults)

	return i.aiClient.GenerateContent(prompt, 2000)
}

// formatSnapshots formata snapshots para exibição
func formatSnapshots(snapshots []Snapshot) string {
	if len(snapshots) == 0 {
		return "Nenhum snapshot disponível"
	}

	var sb strings.Builder
	sb.WriteString("ID | Instância | Início | Fim\n")
	sb.WriteString("---|-----------|--------|----\n")
	for _, snap := range snapshots {
		sb.WriteString(fmt.Sprintf("%d | %d | %s | %s\n",
			snap.ID, snap.InstanceNumber,
			snap.BeginTime.Format("2006-01-02 15:04:05"),
			snap.EndTime.Format("2006-01-02 15:04:05")))
	}
	return sb.String()
}

// parseAWRRequestFromJSON parseia JSON da resposta da IA
func parseAWRRequestFromJSON(jsonStr string) (*AWRReportRequest, error) {
	// Limpar resposta (pode conter markdown code blocks)
	jsonStr = strings.TrimSpace(jsonStr)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	req := &AWRReportRequest{
		ReportType: "html",
		InstanceID: 1,
	}

	// Parsing básico (implementação simplificada)
	// Em produção, usar encoding/json

	// Extrair begin_time
	if match := regexp.MustCompile(`"begin_time"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		if t, err := time.Parse("2006-01-02 15:04:05", match[1]); err == nil {
			req.BeginTime = t
		}
	}

	// Extrair end_time
	if match := regexp.MustCompile(`"end_time"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		if t, err := time.Parse("2006-01-02 15:04:05", match[1]); err == nil {
			req.EndTime = t
		}
	}

	// Extrair begin_snapshot
	if match := regexp.MustCompile(`"begin_snapshot"\s*:\s*(\d+)`).FindStringSubmatch(jsonStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.BeginSnapshot)
	}

	// Extrair end_snapshot
	if match := regexp.MustCompile(`"end_snapshot"\s*:\s*(\d+)`).FindStringSubmatch(jsonStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.EndSnapshot)
	}

	// Extrair report_type
	if match := regexp.MustCompile(`"report_type"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		req.ReportType = match[1]
	}

	return req, nil
}

// parseAWRRequestManual parsing manual como fallback
func parseAWRRequestManual(naturalLanguage string, snapshots []Snapshot) *AWRReportRequest {
	req := &AWRReportRequest{
		ReportType: "html",
		InstanceID: 1,
	}

	now := time.Now()
	lower := strings.ToLower(naturalLanguage)

	// Detectar "últimas X horas"
	if match := regexp.MustCompile(`últim[ao]s?\s+(\d+)\s+horas?`).FindStringSubmatch(lower); len(match) > 1 {
		var hours int
		fmt.Sscanf(match[1], "%d", &hours)
		req.EndTime = now
		req.BeginTime = now.Add(-time.Duration(hours) * time.Hour)
		return req
	}

	// Detectar "ontem"
	if strings.Contains(lower, "ontem") {
		req.EndTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		req.BeginTime = req.EndTime.Add(-24 * time.Hour)
		return req
	}

	// Detectar "hoje"
	if strings.Contains(lower, "hoje") {
		req.BeginTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		req.EndTime = now
		return req
	}

	// Detectar snapshots específicos
	if match := regexp.MustCompile(`snap\s+(\d+)\s+e\s+(\d+)`).FindStringSubmatch(lower); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.BeginSnapshot)
		fmt.Sscanf(match[2], "%d", &req.EndSnapshot)
		return req
	}

	// Default: últimas 2 horas
	req.EndTime = now
	req.BeginTime = now.Add(-2 * time.Hour)
	return req
}

// parseASHRequestFromJSON parseia JSON da resposta da IA para ASH
func parseASHRequestFromJSON(jsonStr string) (*ASHRequest, error) {
	jsonStr = strings.TrimSpace(jsonStr)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	req := &ASHRequest{}

	// Extrair sql_id
	if match := regexp.MustCompile(`"sql_id"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		req.SQLID = match[1]
	}

	// Extrair sid
	if match := regexp.MustCompile(`"sid"\s*:\s*(\d+)`).FindStringSubmatch(jsonStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.SID)
	}

	// Extrair serial
	if match := regexp.MustCompile(`"serial"\s*:\s*(\d+)`).FindStringSubmatch(jsonStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.Serial)
	}

	// Extrair begin_time
	if match := regexp.MustCompile(`"begin_time"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		if t, err := time.Parse("2006-01-02 15:04:05", match[1]); err == nil {
			req.BeginTime = t
		}
	}

	// Extrair end_time
	if match := regexp.MustCompile(`"end_time"\s*:\s*"([^"]+)"`).FindStringSubmatch(jsonStr); len(match) > 1 {
		if t, err := time.Parse("2006-01-02 15:04:05", match[1]); err == nil {
			req.EndTime = t
		}
	}

	// Extrair duration_minutes
	if match := regexp.MustCompile(`"duration_minutes"\s*:\s*(\d+)`).FindStringSubmatch(jsonStr); len(match) > 1 {
		var minutes int
		fmt.Sscanf(match[1], "%d", &minutes)
		now := time.Now()
		req.EndTime = now
		req.BeginTime = now.Add(-time.Duration(minutes) * time.Minute)
	}

	return req, nil
}

// parseASHRequestManual parsing manual como fallback
func parseASHRequestManual(naturalLanguage string) *ASHRequest {
	req := &ASHRequest{}
	lower := strings.ToLower(naturalLanguage)

	// Detectar SQL ID
	if match := regexp.MustCompile(`sql\s+id\s+([a-z0-9]{13})`).FindStringSubmatch(lower); len(match) > 1 {
		req.SQLID = match[1]
	}

	// Detectar SID
	if match := regexp.MustCompile(`sid\s+(\d+)`).FindStringSubmatch(lower); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.SID)
	}

	// Detectar serial
	if match := regexp.MustCompile(`serial\s+(\d+)`).FindStringSubmatch(lower); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &req.Serial)
	}

	// Detectar duração
	if match := regexp.MustCompile(`últim[ao]s?\s+(\d+)\s+minutos?`).FindStringSubmatch(lower); len(match) > 1 {
		var minutes int
		fmt.Sscanf(match[1], "%d", &minutes)
		now := time.Now()
		req.EndTime = now
		req.BeginTime = now.Add(-time.Duration(minutes) * time.Minute)
	}

	return req
}

