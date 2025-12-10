package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GrokClient implementa o cliente Grok (xAI)
type GrokClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewGrokClient cria um novo cliente Grok
func NewGrokClient(config *AIConfig) (*GrokClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "grok-beta"
	}

	return &GrokClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://api.x.ai/v1/chat/completions",
		client:  createHTTPClient(),
	}, nil
}

func (g *GrokClient) GetProvider() Provider {
	return ProviderGrok
}

func (g *GrokClient) GetModel() string {
	return g.model
}

func (g *GrokClient) SetModel(model string) {
	g.model = model
}

func (g *GrokClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
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

func (g *GrokClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	return generateContentGeneric(g, prompt, maxTokens)
}

func (g *GrokClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(g, topic, context)
}

func (g *GrokClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(g, query, notesContext)
}

func (g *GrokClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(g, question, notesContext)
}

func (g *GrokClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(g, language, description, context)
}

func (g *GrokClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(g, topic)
}

func (g *GrokClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(g, topic, context, numItems)
}

func (g *GrokClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(g, projectName, description)
}

