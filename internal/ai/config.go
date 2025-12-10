package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AIConfig representa a configuração do provedor de IA
type AIConfig struct {
	Provider Provider `json:"provider"`
	Model    string   `json:"model"`
	APIKey   string   `json:"api_key"`
}

// GetConfigPath retorna o caminho do arquivo de configuração
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".snip")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "ai_config.json"), nil
}

// LoadConfig carrega a configuração do arquivo
func LoadConfig() (*AIConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Se o arquivo não existir, retornar configuração vazia
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Tentar carregar de variável de ambiente (compatibilidade)
		apiKey := os.Getenv("GROQ_API_KEY")
		if apiKey != "" {
			return &AIConfig{
				Provider: ProviderGroq,
				Model:    "openai/gpt-oss-120b",
				APIKey:   apiKey,
			}, nil
		}
		return &AIConfig{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var config AIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao parsear configuração: %w", err)
	}

	return &config, nil
}

// SaveConfig salva a configuração no arquivo
func SaveConfig(config *AIConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("erro ao escrever arquivo de configuração: %w", err)
	}

	return nil
}

// GetAPIKey retorna a API key da configuração ou variável de ambiente (compatibilidade)
func GetAPIKey() string {
	config, err := LoadConfig()
	if err == nil && config.APIKey != "" {
		return config.APIKey
	}

	// Fallback para variável de ambiente (compatibilidade)
	return os.Getenv("GROQ_API_KEY")
}

// GetAvailableModels retorna os modelos disponíveis para cada provedor
func GetAvailableModels(provider Provider) []string {
	switch provider {
	case ProviderGroq:
		return []string{
			"openai/gpt-oss-120b",
			"llama-3.1-70b-versatile",
			"llama-3.1-8b-instant",
			"mixtral-8x7b-32768",
			"gemma-7b-it",
		}
	case ProviderOpenAI:
		return []string{
			"gpt-4o",
			"gpt-4o-mini",
			"gpt-4-turbo",
			"gpt-4",
			"gpt-3.5-turbo",
		}
	case ProviderAnthropic:
		return []string{
			"claude-3-5-sonnet-20241022",
			"claude-3-5-haiku-20241022",
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
		}
	case ProviderDeepSeek:
		return []string{
			"deepseek-chat",
			"deepseek-coder",
		}
	case ProviderGrok:
		return []string{
			"grok-beta",
			"grok-2",
		}
	case ProviderOpenRouter:
		return []string{
			"openai/gpt-4o",
			"openai/gpt-4-turbo",
			"anthropic/claude-3.5-sonnet",
			"google/gemini-pro",
			"meta-llama/llama-3.1-70b-instruct",
			"mistralai/mixtral-8x7b-instruct",
		}
	default:
		return []string{}
	}
}
