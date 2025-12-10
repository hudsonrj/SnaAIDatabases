package jira

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JiraConfig armazena a configuração do Jira
type JiraConfig struct {
	URL      string `json:"url"`      // URL base do Jira (ex: https://empresa.atlassian.net)
	Email    string `json:"email"`    // Email do usuário
	APIToken string `json:"api_token"` // API Token do Jira
	Project  string `json:"project"`   // Chave do projeto (ex: PROJ)
}

// GetConfigPath retorna o caminho do arquivo de configuração do Jira
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}
	configDir := filepath.Join(homeDir, ".snip")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de configuração: %w", err)
	}
	return filepath.Join(configDir, "jira_config.json"), nil
}

// LoadConfig carrega a configuração do Jira
func LoadConfig() (*JiraConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &JiraConfig{}, nil
		}
		return nil, fmt.Errorf("erro ao ler configuração: %w", err)
	}

	var config JiraConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao parsear configuração: %w", err)
	}

	return &config, nil
}

// SaveConfig salva a configuração do Jira
func SaveConfig(config *JiraConfig) error {
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

