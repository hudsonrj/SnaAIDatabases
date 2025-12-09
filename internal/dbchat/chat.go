package dbchat

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/dbtypes"
)

// ChatSession representa uma sessão de chat com o banco de dados
type ChatSession struct {
	ID           int
	DatabaseType dbtypes.DatabaseType
	Config       *dbtypes.ConnectionConfig
	Messages     []ChatMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ChatMessage representa uma mensagem no chat
type ChatMessage struct {
	Role      string // "user" ou "assistant"
	Content   string
	Query     string // Query SQL gerada (se aplicável)
	Result    string // Resultado da query (se aplicável)
	Timestamp time.Time
}

// DBChat é o gerenciador de chat com banco de dados
type DBChat struct {
	aiClient *ai.GroqClient
	db       *sql.DB
	dbType   dbtypes.DatabaseType
	config   *dbtypes.ConnectionConfig
	session  *ChatSession
}

// NewDBChat cria uma nova sessão de chat
func NewDBChat(dbType dbtypes.DatabaseType, config *dbtypes.ConnectionConfig, db *sql.DB) (*DBChat, error) {
	aiClient, err := ai.NewGroqClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}

	session := &ChatSession{
		DatabaseType: dbType,
		Config:       config,
		Messages:     []ChatMessage{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return &DBChat{
		aiClient: aiClient,
		db:       db,
		dbType:   dbType,
		config:   config,
		session:  session,
	}, nil
}

// SendMessage envia uma mensagem e recebe resposta
func (c *DBChat) SendMessage(userMessage string) (string, error) {
	// Adicionar mensagem do usuário
	c.session.Messages = append(c.session.Messages, ChatMessage{
		Role:      "user",
		Content:   userMessage,
		Timestamp: time.Now(),
	})

	// Construir contexto do chat
	context := c.buildContext()

	// Gerar resposta com IA
	response, query, result, err := c.generateResponse(userMessage, context)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar resposta: %w", err)
	}

	// Adicionar resposta do assistente
	c.session.Messages = append(c.session.Messages, ChatMessage{
		Role:      "assistant",
		Content:   response,
		Query:     query,
		Result:    result,
		Timestamp: time.Now(),
	})

	c.session.UpdatedAt = time.Now()

	return response, nil
}

// buildContext constrói o contexto para a IA
func (c *DBChat) buildContext() string {
	var context strings.Builder

	context.WriteString(fmt.Sprintf("Você é um assistente especializado em %s.\n\n", c.dbType))
	context.WriteString("Contexto do banco de dados:\n")
	context.WriteString(fmt.Sprintf("- Tipo: %s\n", c.dbType))
	context.WriteString(fmt.Sprintf("- Host: %s\n", c.config.Host))
	if c.config.Database != "" {
		context.WriteString(fmt.Sprintf("- Database: %s\n", c.config.Database))
	}

	// Adicionar histórico de mensagens recentes
	if len(c.session.Messages) > 0 {
		context.WriteString("\nHistórico recente da conversa:\n")
		start := len(c.session.Messages) - 4
		if start < 0 {
			start = 0
		}
		for i := start; i < len(c.session.Messages); i++ {
			msg := c.session.Messages[i]
			context.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
			if msg.Query != "" {
				context.WriteString(fmt.Sprintf("  Query executada: %s\n", msg.Query))
			}
		}
	}

	return context.String()
}

// generateResponse gera resposta da IA e executa queries se necessário
func (c *DBChat) generateResponse(userMessage, context string) (response, query, result string, err error) {
	// Determinar se a mensagem requer execução de query
	needsQuery := c.needsQueryExecution(userMessage)

	if needsQuery {
		// Gerar query com IA
		query, err = c.generateQuery(userMessage, context)
		if err != nil {
			return "", "", "", fmt.Errorf("erro ao gerar query: %w", err)
		}

		// Executar query
		result, err = c.executeQuery(query)
		if err != nil {
			// Se a query falhar, pedir à IA para explicar o erro
			errorContext := fmt.Sprintf("%s\n\nErro ao executar query: %v\nQuery: %s", context, err, query)
			response, err = c.generateErrorExplanation(userMessage, errorContext, query, err.Error())
			return response, query, "", err
		}

		// Gerar resposta interpretando o resultado
		response, err = c.interpretQueryResult(userMessage, context, query, result)
		if err != nil {
			return "", query, result, fmt.Errorf("erro ao interpretar resultado: %w", err)
		}
	} else {
		// Resposta conversacional sem query
		response, err = c.generateConversationalResponse(userMessage, context)
		if err != nil {
			return "", "", "", fmt.Errorf("erro ao gerar resposta: %w", err)
		}
	}

	return response, query, result, nil
}

// needsQueryExecution determina se a mensagem requer execução de query
func (c *DBChat) needsQueryExecution(message string) bool {
	message = strings.ToLower(message)
	
	// Palavras-chave que indicam necessidade de query
	queryKeywords := []string{
		"mostre", "mostrar", "liste", "listar", "consulte", "consultar",
		"quantos", "quais", "quem", "onde", "quando", "como",
		"select", "query", "tabela", "tabelas", "dados", "informações",
		"estatísticas", "status", "verificar", "analisar",
	}

	for _, keyword := range queryKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}

	return false
}

// generateQuery gera uma query SQL baseada na mensagem do usuário
func (c *DBChat) generateQuery(userMessage, context string) (string, error) {
	prompt := fmt.Sprintf(`%s

O usuário solicitou: "%s"

Gere uma query SQL apropriada para %s que atenda à solicitação do usuário.

IMPORTANTE:
- Retorne APENAS a query SQL, sem explicações
- Use sintaxe correta para %s
- Seja específico e preciso
- Inclua apenas colunas necessárias
- Use LIMIT/TOP quando apropriado para evitar resultados muito grandes

Query SQL:`, context, userMessage, c.dbType, c.dbType)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf("Você é um especialista em SQL para %s. Gere queries SQL precisas e otimizadas.", c.dbType),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	query, err := c.aiClient.Chat(messages, 500, 0.3)
	if err != nil {
		return "", err
	}

	// Limpar a query (remover markdown code blocks se houver)
	query = strings.TrimSpace(query)
	query = strings.TrimPrefix(query, "```sql")
	query = strings.TrimPrefix(query, "```")
	query = strings.TrimSuffix(query, "```")
	query = strings.TrimSpace(query)

	return query, nil
}

// executeQuery executa a query no banco de dados
func (c *DBChat) executeQuery(query string) (string, error) {
	if c.db == nil {
		return "", fmt.Errorf("conexão com banco de dados não disponível")
	}

	rows, err := c.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("erro ao obter colunas: %w", err)
	}

	var result strings.Builder

	// Escrever cabeçalho
	result.WriteString(strings.Join(columns, " | "))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("-", len(strings.Join(columns, " | "))))
	result.WriteString("\n")

	// Ler resultados
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			continue
		}

		row := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				row[i] = fmt.Sprintf("%v", val)
			} else {
				row[i] = "NULL"
			}
		}

		result.WriteString(strings.Join(row, " | "))
		result.WriteString("\n")

		count++
		if count >= 100 { // Limitar resultados
			result.WriteString(fmt.Sprintf("\n... e mais resultados (limitado a 100)\n"))
			break
		}
	}

	if count == 0 {
		return "Nenhum resultado encontrado.", nil
	}

	return result.String(), nil
}

// interpretQueryResult interpreta o resultado da query usando IA
func (c *DBChat) interpretQueryResult(userMessage, context, query, result string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nO usuário perguntou: \"%s\"\n\nA seguinte query foi executada:\n```sql\n%s\n```\n\nResultado:\n```\n%s\n```\n\nInterprete os resultados de forma clara e útil para o usuário. Explique o que os dados significam e forneça insights relevantes.", context, userMessage, query, result)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um analista de dados que interpreta resultados de queries SQL de forma clara e útil.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 1500, 0.7)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateErrorExplanation gera explicação de erro usando IA
func (c *DBChat) generateErrorExplanation(userMessage, context, query, errorMsg string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nO usuário tentou executar a seguinte query:\n```sql\n%s\n```\n\nMas ocorreu um erro: %s\n\nExplique o erro de forma clara e sugira como corrigir.", context, query, errorMsg)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em SQL que ajuda a entender e corrigir erros em queries.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 1000, 0.7)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateConversationalResponse gera resposta conversacional sem query
func (c *DBChat) generateConversationalResponse(userMessage, context string) (string, error) {
	prompt := fmt.Sprintf(`%s

O usuário disse: "%s"

Responda de forma útil e conversacional sobre o banco de dados %s.`, context, userMessage, c.dbType)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf("Você é um assistente especializado em %s. Responda perguntas de forma clara e útil.", c.dbType),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 1500, 0.7)
	if err != nil {
		return "", err
	}

	return response, nil
}

// GetSession retorna a sessão atual
func (c *DBChat) GetSession() *ChatSession {
	return c.session
}

// Close fecha a conexão
func (c *DBChat) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

