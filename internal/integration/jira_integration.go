package integration

import (
	"fmt"
	"strings"

	"github.com/snip/internal/dbanalysis"
	"github.com/snip/internal/jira"
)

// CreateJiraEpicFromAnalysis cria um Epic no Jira a partir de uma análise
func CreateJiraEpicFromAnalysis(analysis *dbanalysis.DBAnalysis) (*jira.Issue, error) {
	config, err := jira.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração Jira: %w", err)
	}

	if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Project == "" {
		return nil, fmt.Errorf("configuração do Jira incompleta. Execute: snip jira config")
	}

	client, err := jira.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Jira: %w", err)
	}

	// Criar descrição do Epic
	var description strings.Builder
	description.WriteString(fmt.Sprintf("Epic criado automaticamente a partir da análise: %s\n\n", analysis.Title))
	description.WriteString(fmt.Sprintf("**Tipo de Banco:** %s\n", analysis.DatabaseType))
	description.WriteString(fmt.Sprintf("**Tipo de Análise:** %s\n\n", analysis.AnalysisType))
	
	if analysis.Result != "" {
		description.WriteString("**Resultado da Análise:**\n")
		// Limitar tamanho da descrição
		result := analysis.Result
		if len(result) > 2000 {
			result = result[:2000] + "..."
		}
		description.WriteString(result)
		description.WriteString("\n\n")
	}

	if analysis.AIInsights != "" {
		description.WriteString("**Insights da IA:**\n")
		insights := analysis.AIInsights
		if len(insights) > 1000 {
			insights = insights[:1000] + "..."
		}
		description.WriteString(insights)
	}

	epic, err := client.CreateEpic(analysis.Title, description.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao criar Epic: %w", err)
	}

	return epic, nil
}

// CreateJiraIssuesFromAnalysis cria Issues (cards) no Jira a partir de problemas identificados na análise
func CreateJiraIssuesFromAnalysis(analysis *dbanalysis.DBAnalysis, epicKey string) ([]string, error) {
	config, err := jira.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração Jira: %w", err)
	}

	if config.URL == "" || config.Email == "" || config.APIToken == "" || config.Project == "" {
		return nil, fmt.Errorf("configuração do Jira incompleta. Execute: snip jira config")
	}

	client, err := jira.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Jira: %w", err)
	}

	var issueKeys []string

	// Extrair problemas da análise usando IA ou parsing
	problems := extractProblemsFromAnalysis(analysis)

	for _, problem := range problems {
		issue, err := client.CreateIssue(
			problem.Title,
			problem.Description,
			problem.IssueType,
			epicKey,
		)
		if err != nil {
			// Continuar mesmo se uma issue falhar
			continue
		}
		issueKeys = append(issueKeys, issue.Key)
	}

	return issueKeys, nil
}

// Problem representa um problema identificado na análise
type Problem struct {
	Title       string
	Description string
	IssueType   string
	Priority    string
}

// extractProblemsFromAnalysis extrai problemas da análise
func extractProblemsFromAnalysis(analysis *dbanalysis.DBAnalysis) []Problem {
	var problems []Problem

	// Analisar insights da IA para extrair problemas
	if analysis.AIInsights != "" {
		// Procurar por padrões de problemas
		insights := strings.ToLower(analysis.AIInsights)
		
		// Problemas comuns
		if strings.Contains(insights, "backup") || strings.Contains(insights, "backup") {
			problems = append(problems, Problem{
				Title:       "Revisar estratégia de backup",
				Description: "Problema identificado relacionado a backups na análise",
				IssueType:   "Task",
				Priority:    "High",
			})
		}

		if strings.Contains(insights, "performance") || strings.Contains(insights, "lento") {
			problems = append(problems, Problem{
				Title:       "Otimizar performance",
				Description: "Problemas de performance identificados na análise",
				IssueType:   "Task",
				Priority:    "High",
			})
		}

		if strings.Contains(insights, "segurança") || strings.Contains(insights, "security") {
			problems = append(problems, Problem{
				Title:       "Revisar segurança",
				Description: "Problemas de segurança identificados na análise",
				IssueType:   "Bug",
				Priority:    "High",
			})
		}

		if strings.Contains(insights, "espaço") || strings.Contains(insights, "storage") || strings.Contains(insights, "disco") {
			problems = append(problems, Problem{
				Title:       "Gerenciar espaço em disco",
				Description: "Problemas de espaço em disco identificados na análise",
				IssueType:   "Task",
				Priority:    "Medium",
			})
		}
	}

	// Se não encontrou problemas específicos, criar um genérico
	if len(problems) == 0 {
		problems = append(problems, Problem{
			Title:       fmt.Sprintf("Revisar análise: %s", analysis.Title),
			Description: fmt.Sprintf("Revisar e implementar melhorias identificadas na análise %s", analysis.Title),
			IssueType:   "Task",
			Priority:    "Medium",
		})
	}

	return problems
}

