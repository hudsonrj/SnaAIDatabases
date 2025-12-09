package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	GroqAPIURL   = "https://api.groq.com/openai/v1/chat/completions"
	DefaultModel = "openai/gpt-oss-120b" // Modelo mais recente e poderoso disponível na Groq
)

type GroqClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type GroqError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func NewGroqClient() (*GroqClient, error) {
	apiKey := GetAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}

	return &GroqClient{
		apiKey:  apiKey,
		model:   DefaultModel,
		baseURL: GroqAPIURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (g *GroqClient) SetModel(model string) {
	g.model = model
}

func (g *GroqClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
	reqBody := ChatRequest{
		Model:    g.model,
		Messages: messages,
	}

	if maxTokens > 0 {
		reqBody.MaxTokens = maxTokens
	}

	if temperature > 0 {
		reqBody.Temperature = temperature
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", g.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var groqErr GroqError
		if err := json.Unmarshal(body, &groqErr); err == nil {
			return "", fmt.Errorf("groq API error: %s (type: %s, code: %s)",
				groqErr.Error.Message, groqErr.Error.Type, groqErr.Error.Code)
		}
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (g *GroqClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return g.Chat(messages, maxTokens, 0.7)
}

func (g *GroqClient) GenerateNoteContent(topic string, context string) (string, error) {
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

	return g.Chat(messages, 2000, 0.7)
}

func (g *GroqClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
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

	return g.Chat(messages, 100, 0.3)
}

func (g *GroqClient) AnswerQuestion(question string, notesContext []string) (string, error) {
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

	return g.Chat(messages, 1500, 0.7)
}

func (g *GroqClient) GenerateCode(language string, description string, context string) (string, error) {
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

	return g.Chat(messages, 2000, 0.3)
}

func (g *GroqClient) GenerateTips(topic string) (string, error) {
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

	return g.Chat(messages, 1000, 0.7)
}

func (g *GroqClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
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

	result, err := g.Chat(messages, 500, 0.5)
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

func (g *GroqClient) GenerateProjectPlan(projectName string, description string) (string, error) {
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

	return g.Chat(messages, 2000, 0.7)
}
