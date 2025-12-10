package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snip/internal/dbchat"
	"github.com/snip/internal/dbconnection"
	"github.com/snip/internal/dbanalysis"
	"github.com/spf13/cobra"
)

var (
	dbChatDBType      string
	dbChatHost        string
	dbChatPort        int
	dbChatDatabase    string
	dbChatUsername    string
	dbChatPassword    string
	dbChatIsRemote    bool
	dbChatJDBCURL     string
	dbChatConnString  string
)

func init() {
	dbChatCmd.Flags().StringVarP(&dbChatDBType, "db-type", "d", "", "Tipo de banco (oracle, sqlserver, mysql, postgresql, mongodb)")
	dbChatCmd.Flags().StringVar(&dbChatHost, "host", "localhost", "Host do banco de dados")
	dbChatCmd.Flags().IntVar(&dbChatPort, "port", 0, "Porta do banco de dados")
	dbChatCmd.Flags().StringVar(&dbChatDatabase, "database", "", "Nome do banco de dados")
	dbChatCmd.Flags().StringVarP(&dbChatUsername, "username", "u", "", "Usuário do banco de dados")
	dbChatCmd.Flags().StringVarP(&dbChatPassword, "password", "p", "", "Senha do banco de dados")
	dbChatCmd.Flags().BoolVar(&dbChatIsRemote, "remote", false, "Conexão remota")
	dbChatCmd.Flags().StringVar(&dbChatJDBCURL, "jdbc-url", "", "URL JDBC completa")
	dbChatCmd.Flags().StringVar(&dbChatConnString, "conn-string", "", "String de conexão completa")

	rootCmd.AddCommand(dbChatCmd)
}

var dbChatCmd = &cobra.Command{
	Use:   "db-chat",
	Short: "Chat interativo com banco de dados usando IA",
	Long: `Inicia uma sessão de chat interativa com um banco de dados usando IA.

A IA pode gerar queries SQL baseadas em suas perguntas e interpretar os resultados.

Exemplos:
  snip db-chat --db-type postgresql --host localhost --port 5432 --database mydb --username user --password pass
  snip db-chat --db-type mysql --jdbc-url "jdbc:mysql://localhost:3306/db"
  snip db-chat --db-type sqlserver --conn-string "Server=localhost;Database=AdventureWorks;User Id=sa;Password=senha;"

Para sair do chat, digite 'exit', 'quit' ou 'sair'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbChatDBType == "" {
			fmt.Println("Erro: tipo de banco é obrigatório (use --db-type)")
			return
		}

		dbType := parseDatabaseType(dbChatDBType)
		config := &dbanalysis.ConnectionConfig{
			Type:            dbType,
			Host:            dbChatHost,
			Port:            dbChatPort,
			Database:        dbChatDatabase,
			Username:        dbChatUsername,
			Password:        dbChatPassword,
			IsRemote:        dbChatIsRemote,
			JDBCURL:         dbChatJDBCURL,
			ConnectionString: dbChatConnString,
		}

		// Conectar ao banco
		connector, err := dbconnection.GetConnector(dbType)
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			return
		}

		db, err := connector.Connect(config)
		if err != nil {
			fmt.Printf("Erro ao conectar: %v\n", err)
			return
		}
		defer db.Close()

		// Criar sessão de chat
		chat, err := dbchat.NewDBChat(dbType, config, db)
		if err != nil {
			fmt.Printf("Erro ao criar sessão de chat: %v\n", err)
			return
		}
		defer chat.Close()

		fmt.Println("🤖 Chat com Banco de Dados iniciado!")
		fmt.Println("Digite suas perguntas ou solicitações. A IA executará queries automaticamente e responderá com os resultados.")
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

