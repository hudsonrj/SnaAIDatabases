package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GetProvider retorna o provedor
func (g *GroqClient) GetProvider() Provider {
	return ProviderGroq
}

// GetModel retorna o modelo atual
func (g *GroqClient) GetModel() string {
	return g.model
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
	return generateContentGeneric(g, prompt, maxTokens)
}

func (g *GroqClient) GenerateNoteContent(topic string, context string) (string, error) {
	return generateNoteContentGeneric(g, topic, context)
}

func (g *GroqClient) ImproveSearchQuery(query string, notesContext []string) (string, error) {
	return improveSearchQueryGeneric(g, query, notesContext)
}

func (g *GroqClient) AnswerQuestion(question string, notesContext []string) (string, error) {
	return answerQuestionGeneric(g, question, notesContext)
}

func (g *GroqClient) GenerateCode(language string, description string, context string) (string, error) {
	return generateCodeGeneric(g, language, description, context)
}

func (g *GroqClient) GenerateTips(topic string) (string, error) {
	return generateTipsGeneric(g, topic)
}

func (g *GroqClient) GenerateChecklist(topic string, context string, numItems int) ([]string, error) {
	return generateChecklistGeneric(g, topic, context, numItems)
}

func (g *GroqClient) GenerateProjectPlan(projectName string, description string) (string, error) {
	return generateProjectPlanGeneric(g, projectName, description)
}
