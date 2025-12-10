package dbmaintenance

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
)

// MaintenancePlan representa um plano de manutenção
type MaintenancePlan struct {
	ID          int
	Title       string
	Description string
	Priority    string // high, medium, low
	Status      string // pending, in_progress, completed
	Tasks       []MaintenanceTask
	EstimatedDuration time.Duration
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MaintenanceTask representa uma tarefa de manutenção
type MaintenanceTask struct {
	ID          int
	Title       string
	Description string
	Steps       []string
	Priority    string
	Status      string
	EstimatedTime time.Duration
	Dependencies []int // IDs de tarefas que devem ser executadas antes
}

// MaintenancePlanner gera planos de manutenção usando IA
type MaintenancePlanner struct {
	aiClient ai.AIClient
}

// NewMaintenancePlanner cria um novo planejador de manutenção
func NewMaintenancePlanner() (*MaintenancePlanner, error) {
	aiClient, err := ai.NewAIClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &MaintenancePlanner{aiClient: aiClient}, nil
}

// GenerateMaintenancePlan gera um plano de manutenção baseado em análise
func (m *MaintenancePlanner) GenerateMaintenancePlan(analysisResult string, analysisType string, dbType string) (*MaintenancePlan, error) {
	// Usar IA para gerar plano
	planJSON, err := m.generatePlanWithAI(analysisResult, analysisType, dbType)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar plano: %w", err)
	}

	// Parsear JSON
	var plan MaintenancePlan
	err = json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear plano: %w", err)
	}

	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()
	plan.Status = "pending"

	return &plan, nil
}

// generatePlanWithAI usa IA para gerar plano de manutenção
func (m *MaintenancePlanner) generatePlanWithAI(analysisResult, analysisType, dbType string) (string, error) {
	prompt := fmt.Sprintf(`Você é um DBA experiente especializado em %s.

Analise o seguinte resultado de análise (%s) e crie um plano de manutenção detalhado:

%s

Crie um plano de manutenção completo com:
1. Título e descrição do plano
2. Prioridade (high, medium, low)
3. Lista de tarefas com:
   - Título da tarefa
   - Descrição detalhada
   - Passo a passo para execução
   - Prioridade individual
   - Tempo estimado (em minutos)
   - Dependências (IDs de outras tarefas)

Retorne APENAS um JSON válido com a seguinte estrutura:
{
  "title": "Título do Plano",
  "description": "Descrição detalhada",
  "priority": "high|medium|low",
  "tasks": [
    {
      "title": "Título da Tarefa",
      "description": "Descrição",
      "steps": ["passo 1", "passo 2", ...],
      "priority": "high|medium|low",
      "estimated_time_minutes": 30,
      "dependencies": []
    }
  ]
}`, dbType, analysisType, analysisResult)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um DBA experiente que cria planos de manutenção detalhados e práticos para bancos de dados.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := m.aiClient.Chat(messages, 3000, 0.5)
	if err != nil {
		return "", err
	}

	// Limpar resposta
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	return response, nil
}

// SuggestMaintenanceActions sugere ações de manutenção baseadas em análise
func (m *MaintenancePlanner) SuggestMaintenanceActions(analysisResult string, dbType string) (string, error) {
	prompt := fmt.Sprintf(`Analise o seguinte resultado de análise de banco de dados %s e sugira ações de manutenção prioritárias:

%s

Forneça:
1. Problemas identificados
2. Ações recomendadas (priorizadas)
3. Impacto esperado de cada ação
4. Tempo estimado para implementação

Formate a resposta em markdown de forma clara e organizada.`, dbType, analysisResult)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um DBA experiente que identifica problemas e sugere ações de manutenção prioritárias.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := m.aiClient.Chat(messages, 2000, 0.7)
	if err != nil {
		return "", err
	}

	return response, nil
}

// FormatPlan formata o plano para exibição
func (m *MaintenancePlan) FormatPlan() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# %s\n\n", m.Title))
	result.WriteString(fmt.Sprintf("**Prioridade:** %s\n", m.Priority))
	result.WriteString(fmt.Sprintf("**Status:** %s\n", m.Status))
	result.WriteString(fmt.Sprintf("**Duração Estimada:** %s\n\n", m.EstimatedDuration.String()))
	result.WriteString(fmt.Sprintf("%s\n\n", m.Description))

	result.WriteString("## Tarefas\n\n")

	for i, task := range m.Tasks {
		result.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, task.Title))
		result.WriteString(fmt.Sprintf("**Prioridade:** %s\n", task.Priority))
		result.WriteString(fmt.Sprintf("**Tempo Estimado:** %s\n", task.EstimatedTime.String()))
		result.WriteString(fmt.Sprintf("**Status:** %s\n\n", task.Status))
		result.WriteString(fmt.Sprintf("%s\n\n", task.Description))

		if len(task.Steps) > 0 {
			result.WriteString("**Passo a passo:**\n\n")
			for j, step := range task.Steps {
				result.WriteString(fmt.Sprintf("%d. %s\n", j+1, step))
			}
			result.WriteString("\n")
		}

		if len(task.Dependencies) > 0 {
			result.WriteString(fmt.Sprintf("**Dependências:** Tarefas %v\n\n", task.Dependencies))
		}
	}

	return result.String()
}
