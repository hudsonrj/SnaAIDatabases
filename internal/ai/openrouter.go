package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenRouterClient implementa o cliente OpenRouter
type OpenRouterClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenRouterClient cria um novo cliente OpenRouter
func NewOpenRouterClient(config *AIConfig) (*OpenRouterClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "openai/gpt-4o"
	}

	return &OpenRouterClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://openrouter.ai/api/v1/chat/completions",
		client:  createHTTPClient(),
	}, nil
}

func (o *OpenRouterClient) GetProvider() Provider {
	return ProviderOpenRouter
}

func (o *OpenRouterClient) GetModel() string {
	return o.model
}

func (o *OpenRouterClient) SetModel(model string) {
	o.model = model
}

func (o *OpenRouterClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
	reqBody := ChatRequest{
		Model:    o.model,
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

	req, err := http.NewRequest("POST", o.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/snip") // Opcional mas recomendado
	req.Header.Set("X-Title", "SnipAI Databases")             // Opcional mas recomendado

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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

func (o *OpenRouterClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	return generateContentGeneric(o, prompt, maxTokens)
}

func (o *OpenRouterClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(o, topic, context)
}

func (o *OpenRouterClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(o, query, notesContext)
}

func (o *OpenRouterClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(o, question, notesContext)
}

func (o *OpenRouterClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(o, language, description, context)
}

func (o *OpenRouterClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(o, topic)
}

func (o *OpenRouterClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(o, topic, context, numItems)
}

func (o *OpenRouterClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(o, projectName, description)
}

