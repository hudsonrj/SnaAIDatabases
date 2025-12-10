package ai

import (
	"fmt"
	"net/http"
	"time"
)

// createHTTPClient cria um cliente HTTP padrão
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// Provider representa um provedor de IA
type Provider string

const (
	ProviderGroq       Provider = "groq"
	ProviderOpenAI     Provider = "openai"
	ProviderAnthropic   Provider = "anthropic"
	ProviderDeepSeek   Provider = "deepseek"
	ProviderGrok       Provider = "grok"
	ProviderOpenRouter Provider = "openrouter"
)

// AIClient é a interface genérica para clientes de IA
type AIClient interface {
	Chat(messages []Message, maxTokens int, temperature float64) (string, error)
	GenerateContent(prompt string, maxTokens int) (string, error)
	GenerateNoteContent(topic string, context string) (string, error)
	ImproveSearchQuery(query string, notesContext []string) (string, error)
	AnswerQuestion(question string, notesContext []string) (string, error)
	GenerateCode(language string, description string, context string) (string, error)
	GenerateTips(topic string) (string, error)
	GenerateChecklist(topic string, context string, numItems int) ([]string, error)
	GenerateProjectPlan(projectName string, description string) (string, error)
	SetModel(model string)
	GetProvider() Provider
	GetModel() string
}

// Message representa uma mensagem no chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewAIClient cria um novo cliente de IA baseado na configuração
func NewAIClient() (AIClient, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	if config.Provider == "" {
		// Fallback para Groq se não configurado
		config.Provider = ProviderGroq
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada. Execute: snip ai config")
	}

	switch config.Provider {
	case ProviderGroq:
		return NewGroqClientWithConfig(config)
	case ProviderOpenAI:
		return NewOpenAIClient(config)
	case ProviderAnthropic:
		return NewAnthropicClient(config)
	case ProviderDeepSeek:
		return NewDeepSeekClient(config)
	case ProviderGrok:
		return NewGrokClient(config)
	case ProviderOpenRouter:
		return NewOpenRouterClient(config)
	default:
		return nil, fmt.Errorf("provedor não suportado: %s", config.Provider)
	}
}

// NewGroqClientWithConfig cria um cliente Groq com configuração
func NewGroqClientWithConfig(config *AIConfig) (*GroqClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key não configurada")
	}

	model := config.Model
	if model == "" {
		model = "openai/gpt-oss-120b" // Default Groq
	}

	return &GroqClient{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: "https://api.groq.com/openai/v1/chat/completions",
		client:  createHTTPClient(),
	}, nil
}

