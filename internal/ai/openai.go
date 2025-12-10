package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIClient implementa o cliente OpenAI
type OpenAIClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenAIClient cria um novo cliente OpenAI
func NewOpenAIClient(config *AIConfig) (*OpenAIClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "gpt-4o"
	}

	return &OpenAIClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://api.openai.com/v1/chat/completions",
		client:  createHTTPClient(),
	}, nil
}

func (o *OpenAIClient) GetProvider() Provider {
	return ProviderOpenAI
}

func (o *OpenAIClient) GetModel() string {
	return o.model
}

func (o *OpenAIClient) SetModel(model string) {
	o.model = model
}

func (o *OpenAIClient) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
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

func (o *OpenAIClient) GenerateContent(prompt string, maxTokens int) (string, error) {
	return generateContentGeneric(o, prompt, maxTokens)
}

func (o *OpenAIClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(o, topic, context)
}

func (o *OpenAIClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(o, query, notesContext)
}

func (o *OpenAIClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(o, question, notesContext)
}

func (o *OpenAIClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(o, language, description, context)
}

func (o *OpenAIClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(o, topic)
}

func (o *OpenAIClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(o, topic, context, numItems)
}

func (o *OpenAIClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(o, projectName, description)
}

