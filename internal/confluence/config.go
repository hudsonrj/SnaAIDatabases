package confluence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfluenceConfig armazena a configuração do Confluence
type ConfluenceConfig struct {
	URL      string `json:"url"`       // URL base do Confluence (ex: https://empresa.atlassian.net)
	Email    string `json:"email"`     // Email do usuário
	APIToken string `json:"api_token"` // API Token
	Space    string `json:"space"`     // Chave do espaço (ex: DB)
}

// GetConfigPath retorna o caminho do arquivo de configuração do Confluence
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}
	configDir := filepath.Join(homeDir, ".snip")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de configuração: %w", err)
	}
	return filepath.Join(configDir, "confluence_config.json"), nil
}

// LoadConfig carrega a configuração do Confluence
func LoadConfig() (*ConfluenceConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfluenceConfig{}, nil
		}
		return nil, fmt.Errorf("erro ao ler configuração: %w", err)
	}

	var config ConfluenceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao parsear configuração: %w", err)
	}

	return &config, nil
}

// SaveConfig salva a configuração do Confluence
func SaveConfig(config *ConfluenceConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("erro ao salvar configuração: %w", err)
	}

	return nil
}

