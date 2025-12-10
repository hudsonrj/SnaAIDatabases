package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AnthropicClient implementa o cliente Anthropic (Claude)
type AnthropicClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewAnthropicClient cria um novo cliente Anthropic
func NewAnthropicClient(config *AIConfig) (*AnthropicClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	return &AnthropicClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://api.anthropic.com/v1/messages",
		client:  createHTTPClient(),
	}, nil
}

func (a *AnthropicClient) GetProvider() Provider {
	return ProviderAnthropic
}

func (a *AnthropicClient) GetModel() string {
	return a.model
}

func (a *AnthropicClient) SetModel(model string) {
	a.model = model
}

func (a *AnthropicClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
	// Anthropic usa formato diferente
	anthropicMessages := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		anthropicMessages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	reqBody := map[string]interface{}{
		"model":     a.model,
		"messages":  anthropicMessages,
		"max_tokens": maxTokens,
	}

	if temperature > 0 {
		reqBody["temperature"] = temperature
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", a.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.client.Do(req)
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

	var response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
}

func (a *AnthropicClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	return generateContentGeneric(a, prompt, maxTokens)
}

func (a *AnthropicClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(a, topic, context)
}

func (a *AnthropicClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(a, query, notesContext)
}

func (a *AnthropicClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(a, question, notesContext)
}

func (a *AnthropicClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(a, language, description, context)
}

func (a *AnthropicClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(a, topic)
}

func (a *AnthropicClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(a, topic, context, numItems)
}

func (a *AnthropicClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(a, projectName, description)
}

