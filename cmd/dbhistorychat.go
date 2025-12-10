package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snip/internal/database"
	"github.com/snip/internal/dbhistorychat"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbHistoryChatCmd)
}

var dbHistoryChatCmd = &cobra.Command{
	Use:   "db-history chat",
	Short: "Chat interativo com o histórico de análises usando IA",
	Long: `Inicia uma sessão de chat interativa com o banco SQLite que armazena todas as análises de bancos de dados.

A IA pode gerar queries SQL baseadas em suas perguntas e interpretar os resultados para:
- Listar análises por tipo de banco, tipo de análise, data, etc.
- Comparar análises de diferentes datas para ver evolução
- Identificar problemas e insights das análises
- Rastrear a evolução ou degradação dos bancos ao longo do tempo
- Analisar tendências e padrões nas análises
- Responder perguntas sobre resultados específicos

Exemplos de perguntas:
  - "Liste todas as análises do Oracle"
  - "Quantas análises foram feitas este mês?"
  - "Compare as análises de diagnóstico do PostgreSQL entre janeiro e fevereiro"
  - "Quais problemas foram identificados nas análises do SQL Server?"
  - "Mostre a evolução das análises de tuning do MySQL"
  - "Quais insights a IA gerou sobre o MongoDB?"

Para sair do chat, digite 'exit', 'quit' ou 'sair'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Conectar ao banco SQLite interno
		db, err := database.Connect()
		if err != nil {
			fmt.Printf("Erro ao conectar ao banco de dados: %v\n", err)
			return
		}
		defer db.Close()

		// Criar sessão de chat
		chat, err := dbhistorychat.NewDBHistoryChat(db)
		if err != nil {
			fmt.Printf("Erro ao criar sessão de chat: %v\n", err)
			return
		}
		defer chat.Close()

		fmt.Println("🤖 Chat com Histórico de Análises iniciado!")
		fmt.Println("Digite suas perguntas sobre as análises armazenadas.")
		fmt.Println("A IA executará queries automaticamente e responderá com os resultados.")
		fmt.Println("Digite 'exit', 'quit' ou 'sair' para sair.\n")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Você: ")
			if !scanner.Scan() {
				break
			}

			userInput := strings.TrimSpace(scanner.Text())
			if userInput == "" {
				continue
			}

			// Verificar comandos de saída
			userInputLower := strings.ToLower(userInput)
			if userInputLower == "exit" || userInputLower == "quit" || userInputLower == "sair" {
				fmt.Println("\nAté logo! 👋")
				break
			}

			// Enviar mensagem e receber resposta
			fmt.Print("\n🤖 Assistente: ")
			response, err := chat.SendMessage(userInput)
			if err != nil {
				fmt.Printf("❌ Erro: %v\n\n", err)
				continue
			}

			fmt.Println(response)
			fmt.Println()
		}
	},
}

