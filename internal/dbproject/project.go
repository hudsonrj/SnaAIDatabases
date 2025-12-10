package dbproject

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
)

// ProjectFromAnalysis representa um projeto criado a partir de uma análise
type ProjectFromAnalysis struct {
	ProjectName        string
	ProjectDescription string
	Tasks              []ProjectTask
	Priority           string
	CreatedFrom        string // ID ou título da análise
}

// ProjectTask representa uma tarefa do projeto
type ProjectTask struct {
	Title       string
	Description string
	Steps       []string
	Priority    string
	DueDate     *time.Time
	EstimatedTime time.Duration
}

// ProjectGenerator transforma análises em projetos
type ProjectGenerator struct {
	aiClient ai.AIClient
}

// NewProjectGenerator cria um novo gerador de projetos
func NewProjectGenerator() (*ProjectGenerator, error) {
	aiClient, err := ai.NewAIClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &ProjectGenerator{aiClient: aiClient}, nil
}

// GenerateProjectFromAnalysis transforma uma análise em projeto
func (p *ProjectGenerator) GenerateProjectFromAnalysis(analysisTitle string, analysisResult string, analysisType string, dbType string) (*ProjectFromAnalysis, error) {
	// Usar IA para transformar análise em projeto
	projectJSON, err := p.generateProjectWithAI(analysisTitle, analysisResult, analysisType, dbType)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar projeto: %w", err)
	}

	// Parsear JSON
	var project ProjectFromAnalysis
	err = json.Unmarshal([]byte(projectJSON), &project)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear projeto: %w", err)
	}

	project.CreatedFrom = analysisTitle

	return &project, nil
}

// generateProjectWithAI usa IA para transformar análise em projeto
func (p *ProjectGenerator) generateProjectWithAI(analysisTitle, analysisResult, analysisType, dbType string) (string, error) {
	prompt := fmt.Sprintf(`Você é um gerente de projetos especializado em bancos de dados.

Com base na seguinte análise de banco de dados %s (%s), transforme os problemas e recomendações identificados em um projeto estruturado com tarefas e passo a passo.

Análise:
Título: %s
Tipo: %s
Resultado:
%s

Crie um projeto completo com:
1. Nome do projeto (descritivo e claro)
2. Descrição do projeto
3. Prioridade geral (high, medium, low)
4. Lista de tarefas com:
   - Título da tarefa
   - Descrição detalhada
   - Passo a passo para execução
   - Prioridade individual
   - Tempo estimado (em minutos)
   - Data de vencimento sugerida (se aplicável)

Retorne APENAS um JSON válido com a seguinte estrutura:
{
  "project_name": "Nome do Projeto",
  "project_description": "Descrição detalhada",
  "priority": "high|medium|low",
  "tasks": [
    {
      "title": "Título da Tarefa",
      "description": "Descrição",
      "steps": ["passo 1", "passo 2", ...],
      "priority": "high|medium|low",
      "estimated_time_minutes": 60,
      "due_date": "2024-12-31" ou null
    }
  ]
}`, dbType, analysisType, analysisTitle, analysisType, analysisResult)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um gerente de projetos experiente que transforma análises técnicas em projetos estruturados com tarefas e passo a passo.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := p.aiClient.Chat(messages, 3000, 0.5)
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

// GenerateIncidentProject transforma um incidente em projeto
func (p *ProjectGenerator) GenerateIncidentProject(incidentDescription string, analysisResult string, dbType string) (*ProjectFromAnalysis, error) {
	prompt := fmt.Sprintf(`Você é um gerente de projetos especializado em resolução de incidentes de banco de dados.

Com base no seguinte incidente e análise de banco de dados %s, crie um projeto de resolução:

Incidente:
%s

Análise relacionada:
%s

Crie um projeto de resolução com:
1. Nome do projeto (focado na resolução)
2. Descrição do problema e objetivo
3. Prioridade (geralmente high para incidentes)
4. Tarefas de resolução com passo a passo detalhado
5. Tarefas de prevenção para evitar recorrência

Retorne APENAS um JSON válido com a mesma estrutura do GenerateProjectFromAnalysis.`, dbType, incidentDescription, analysisResult)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em resolução de incidentes que cria projetos estruturados para resolver problemas de banco de dados.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := p.aiClient.Chat(messages, 3000, 0.5)
	if err != nil {
		return nil, err
	}

	// Limpar resposta
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var project ProjectFromAnalysis
	err = json.Unmarshal([]byte(response), &project)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear projeto: %w", err)
	}

	project.CreatedFrom = "incident: " + incidentDescription
	return &project, nil
}

// FormatProject formata o projeto para exibição
func (p *ProjectFromAnalysis) FormatProject() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# %s\n\n", p.ProjectName))
	result.WriteString(fmt.Sprintf("**Prioridade:** %s\n", p.Priority))
	if p.CreatedFrom != "" {
		result.WriteString(fmt.Sprintf("**Criado a partir de:** %s\n", p.CreatedFrom))
	}
	result.WriteString("\n")
	result.WriteString(fmt.Sprintf("%s\n\n", p.ProjectDescription))

	result.WriteString("## Tarefas\n\n")

	for i, task := range p.Tasks {
		result.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, task.Title))
		result.WriteString(fmt.Sprintf("**Prioridade:** %s\n", task.Priority))
		result.WriteString(fmt.Sprintf("**Tempo Estimado:** %s\n", task.EstimatedTime.String()))
		if task.DueDate != nil {
			result.WriteString(fmt.Sprintf("**Data de Vencimento:** %s\n", task.DueDate.Format("2006-01-02")))
		}
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("%s\n\n", task.Description))

		if len(task.Steps) > 0 {
			result.WriteString("**Passo a passo:**\n\n")
			for j, step := range task.Steps {
				result.WriteString(fmt.Sprintf("%d. %s\n", j+1, step))
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}

