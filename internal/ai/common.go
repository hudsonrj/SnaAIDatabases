package ai

import (
	"fmt"
	"strings"
)

// Funções genéricas compartilhadas entre todos os clientes

func generateContentGeneric(client AIClient, prompt string, maxTokens int) (string, error) {
	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, maxTokens, 0.7)
}

func generateNoteContentGeneric(client AIClient, topic string, context string) (string, error) {
	prompt := fmt.Sprintf(`Você é um assistente de anotações inteligente. Crie conteúdo útil e bem estruturado sobre o tópico: "%s"

%s

Por favor, crie um conteúdo detalhado, organizado e útil sobre este tópico. Use formatação markdown quando apropriado.`, topic, context)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um assistente especializado em criar anotações bem estruturadas e úteis.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 2000, 0.7)
}

func improveSearchQueryGeneric(client AIClient, query string, notesContext []string) (string, error) {
	contextStr := ""
	if len(notesContext) > 0 {
		contextStr = fmt.Sprintf("\n\nContexto das notas existentes:\n%s",
			notesContext[0])
		if len(notesContext) > 1 {
			contextStr += fmt.Sprintf("\n... e mais %d notas relacionadas", len(notesContext)-1)
		}
	}

	prompt := fmt.Sprintf(`Melhore esta consulta de busca para encontrar notas relevantes: "%s"%s

Retorne apenas a consulta melhorada, sem explicações adicionais.`, query, contextStr)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um assistente especializado em melhorar consultas de busca para encontrar informações relevantes.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 100, 0.3)
}

func answerQuestionGeneric(client AIClient, question string, notesContext []string) (string, error) {
	contextStr := ""
	if len(notesContext) > 0 {
		contextStr = "\n\nInformações das suas notas:\n"
		for i, note := range notesContext {
			if i >= 3 {
				contextStr += fmt.Sprintf("\n... e mais %d notas", len(notesContext)-3)
				break
			}
			contextStr += note + "\n\n"
		}
	}

	prompt := fmt.Sprintf(`Responda a seguinte pergunta com base nas informações disponíveis:%s

Pergunta: %s

Se a resposta não estiver nas notas fornecidas, você pode usar seu conhecimento geral, mas mencione isso.`, contextStr, question)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um assistente inteligente que responde perguntas com base nas anotações do usuário e conhecimento geral.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 1500, 0.7)
}

func generateCodeGeneric(client AIClient, language string, description string, context string) (string, error) {
	prompt := fmt.Sprintf(`Gere código %s para: %s

%s

Por favor, forneça código completo, bem comentado e seguindo as melhores práticas.`, language, description, context)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um programador experiente que gera código limpo, bem documentado e seguindo as melhores práticas.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 2000, 0.3)
}

func generateTipsGeneric(client AIClient, topic string) (string, error) {
	prompt := fmt.Sprintf(`Forneça dicas úteis e práticas sobre: %s

Formate as dicas de forma clara e organizada, usando markdown.`, topic)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um assistente que fornece dicas práticas e úteis sobre diversos tópicos.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 1000, 0.7)
}

func generateChecklistGeneric(client AIClient, topic string, context string, numItems int) ([]string, error) {
	prompt := fmt.Sprintf(`Crie uma lista de verificação (checklist) com %d itens sobre: "%s"

%s

Retorne APENAS os itens da checklist, um por linha, sem numeração, sem marcadores, sem explicações adicionais. Cada linha deve ser um item claro e específico.`, numItems, topic, context)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um assistente especializado em criar listas de verificação práticas e bem estruturadas.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	result, err := client.Chat(messages, 500, 0.5)
	if err != nil {
		return nil, err
	}

	// Parse the result into individual items
	lines := strings.Split(result, "\n")
	var items []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Remove common prefixes like "- ", "* ", numbers, etc.
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimPrefix(line, "• ")
		// Remove numbering (e.g., "1. ", "2. ")
		if len(line) > 2 && line[1] == '.' && line[0] >= '0' && line[0] <= '9' {
			line = strings.TrimSpace(line[2:])
		}
		if line != "" && len(line) > 3 {
			items = append(items, line)
		}
	}

	// Limit to requested number
	if len(items) > numItems {
		items = items[:numItems]
	}

	return items, nil
}

func generateProjectPlanGeneric(client AIClient, projectName string, description string) (string, error) {
	prompt := fmt.Sprintf(`Crie um plano de projeto detalhado para: "%s"

Descrição: %s

O plano deve incluir:
1. Objetivos principais
2. Tarefas principais organizadas por fase
3. Prioridades sugeridas
4. Marcos importantes

Formate o resultado em markdown.`, projectName, description)

	messages := []Message{
		{
			Role:    "system",
			Content: "Você é um gerente de projetos experiente que cria planos detalhados e práticos.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return client.Chat(messages, 2000, 0.7)
}

