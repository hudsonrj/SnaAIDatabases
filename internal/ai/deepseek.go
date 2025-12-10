package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeepSeekClient implementa o cliente DeepSeek
type DeepSeekClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewDeepSeekClient cria um novo cliente DeepSeek
func NewDeepSeekClient(config *AIConfig) (*DeepSeekClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "deepseek-chat"
	}

	return &DeepSeekClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://api.deepseek.com/v1/chat/completions",
		client:  createHTTPClient(),
	}, nil
}

func (d *DeepSeekClient) GetProvider() Provider {
	return ProviderDeepSeek
}

func (d *DeepSeekClient) GetModel() string {
	return d.model
}

func (d *DeepSeekClient) SetModel(model string) {
	d.model = model
}

func (d *DeepSeekClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
	reqBody := ChatRequest{
		Model:    d.model,
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

	req, err := http.NewRequest("POST", d.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
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

func (d *DeepSeekClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	return generateContentGeneric(d, prompt, maxTokens)
}

func (d *DeepSeekClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(d, topic, context)
}

func (d *DeepSeekClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(d, query, notesContext)
}

func (d *DeepSeekClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(d, question, notesContext)
}

func (d *DeepSeekClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(d, language, description, context)
}

func (d *DeepSeekClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(d, topic)
}

func (d *DeepSeekClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(d, topic, context, numItems)
}

func (d *DeepSeekClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(d, projectName, description)
}

