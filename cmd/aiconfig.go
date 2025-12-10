package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snip/internal/ai"
	"github.com/spf13/cobra"
)

var (
	aiConfigProvider string
	aiConfigModel    string
	aiConfigAPIKey   string
	aiConfigShow     bool
)

func init() {
	aiConfigCmd.Flags().StringVarP(&aiConfigProvider, "provider", "p", "", "Provedor de IA (groq, openai, anthropic, deepseek, grok, openrouter)")
	aiConfigCmd.Flags().StringVarP(&aiConfigModel, "model", "m", "", "Modelo a ser usado")
	aiConfigCmd.Flags().StringVarP(&aiConfigAPIKey, "api-key", "k", "", "API Key")
	aiConfigCmd.Flags().BoolVar(&aiConfigShow, "show", false, "Mostrar configuração atual")

	rootCmd.AddCommand(aiConfigCmd)
}

var aiConfigCmd = &cobra.Command{
	Use:   "ai config",
	Short: "Configurar provedor de IA e API key",
	Long: `Configura o provedor de IA, modelo e API key.

Provedores suportados:
  - groq: Groq API
  - openai: OpenAI
  - anthropic: Anthropic (Claude)
  - deepseek: DeepSeek
  - grok: Grok (xAI)
  - openrouter: OpenRouter

Exemplos:
  snip ai config --provider groq --model "openai/gpt-oss-120b" --api-key "sua-chave"
  snip ai config --provider openai --model "gpt-4o" --api-key "sua-chave"
  snip ai config --show  # Mostrar configuração atual
  snip ai config         # Modo interativo`,
	Run: func(cmd *cobra.Command, args []string) {
		if aiConfigShow {
			showConfig()
			return
		}

		config, err := ai.LoadConfig()
		if err != nil {
			fmt.Printf("Erro ao carregar configuração: %v\n", err)
			return
		}

		// Modo interativo se nenhum parâmetro foi fornecido
		if aiConfigProvider == "" && aiConfigModel == "" && aiConfigAPIKey == "" {
			interactiveConfig(config)
			return
		}

		// Atualizar configuração
		if aiConfigProvider != "" {
			config.Provider = ai.Provider(aiConfigProvider)
		}

		if aiConfigModel != "" {
			config.Model = aiConfigModel
		}

		if aiConfigAPIKey != "" {
			config.APIKey = aiConfigAPIKey
		}

		// Se o modelo não foi especificado, usar o padrão do provedor
		if config.Model == "" {
			models := ai.GetAvailableModels(config.Provider)
			if len(models) > 0 {
				config.Model = models[0]
			}
		}

		if err := ai.SaveConfig(config); err != nil {
			fmt.Printf("Erro ao salvar configuração: %v\n", err)
			return
		}

		fmt.Println("✓ Configuração salva com sucesso!")
		fmt.Printf("  Provedor: %s\n", config.Provider)
		fmt.Printf("  Modelo: %s\n", config.Model)
		fmt.Printf("  API Key: %s\n", maskAPIKey(config.APIKey))
	},
}

func showConfig() {
	config, err := ai.LoadConfig()
	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		return
	}

	fmt.Println("📋 Configuração Atual de IA:")
	fmt.Printf("  Provedor: %s\n", config.Provider)
	fmt.Printf("  Modelo: %s\n", config.Model)
	fmt.Printf("  API Key: %s\n", maskAPIKey(config.APIKey))

	if config.Provider != "" {
		models := ai.GetAvailableModels(config.Provider)
		if len(models) > 0 {
			fmt.Println("\n📦 Modelos disponíveis para este provedor:")
			for i, model := range models {
				marker := " "
				if model == config.Model {
					marker = "✓"
				}
				fmt.Printf("  %s %d. %s\n", marker, i+1, model)
			}
		}
	}
}

func interactiveConfig(currentConfig *ai.AIConfig) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔧 Configuração Interativa de IA")
	fmt.Println("")

	// Provedor
	fmt.Println("Provedores disponíveis:")
	providers := []ai.Provider{
		ai.ProviderGroq,
		ai.ProviderOpenAI,
		ai.ProviderAnthropic,
		ai.ProviderDeepSeek,
		ai.ProviderGrok,
		ai.ProviderOpenRouter,
	}

	for i, p := range providers {
		marker := " "
		if string(p) == string(currentConfig.Provider) {
			marker = "✓"
		}
		fmt.Printf("  %s %d. %s\n", marker, i+1, p)
	}

	fmt.Print("\nEscolha o provedor (1-6) [padrão: groq]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selectedProvider ai.Provider
	if input == "" {
		selectedProvider = ai.ProviderGroq
	} else {
		var idx int
		fmt.Sscanf(input, "%d", &idx)
		if idx >= 1 && idx <= len(providers) {
			selectedProvider = providers[idx-1]
		} else {
			selectedProvider = ai.ProviderGroq
		}
	}

	// Modelo
	models := ai.GetAvailableModels(selectedProvider)
	if len(models) > 0 {
		fmt.Println("\nModelos disponíveis:")
		for i, model := range models {
			marker := " "
			if model == currentConfig.Model {
				marker = "✓"
			}
			fmt.Printf("  %s %d. %s\n", marker, i+1, model)
		}

		fmt.Printf("\nEscolha o modelo (1-%d) [padrão: %s]: ", len(models), models[0])
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selectedModel string
		if input == "" {
			selectedModel = models[0]
		} else {
			var idx int
			fmt.Sscanf(input, "%d", &idx)
			if idx >= 1 && idx <= len(models) {
				selectedModel = models[idx-1]
			} else {
				selectedModel = models[0]
			}
		}
		currentConfig.Model = selectedModel
	}

	// API Key
	fmt.Print("\nDigite a API Key (ou Enter para manter a atual): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		currentConfig.APIKey = input
	}

	currentConfig.Provider = selectedProvider

	if err := ai.SaveConfig(currentConfig); err != nil {
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}

	fmt.Println("\n✓ Configuração salva com sucesso!")
	fmt.Printf("  Provedor: %s\n", currentConfig.Provider)
	fmt.Printf("  Modelo: %s\n", currentConfig.Model)
	fmt.Printf("  API Key: %s\n", maskAPIKey(currentConfig.APIKey))
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(não configurada)"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

