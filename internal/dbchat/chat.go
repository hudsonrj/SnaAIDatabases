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
	aiClient ai.AIClient
	db       *sql.DB
	dbType   dbtypes.DatabaseType
	config   *dbtypes.ConnectionConfig
	session  *ChatSession
}

// NewDBChat cria uma nova sessão de chat
func NewDBChat(dbType dbtypes.DatabaseType, config *dbtypes.ConnectionConfig, db *sql.DB) (*DBChat, error) {
	aiClient, err := ai.NewAIClient()
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

	// Adicionar TODO o histórico da conversa (mantém contexto completo)
	if len(c.session.Messages) > 0 {
		context.WriteString("\nHistórico completo da conversa:\n")
		for _, msg := range c.session.Messages {
			if msg.Role == "user" {
				context.WriteString(fmt.Sprintf("\n[Usuário]: %s\n", msg.Content))
			} else {
				context.WriteString(fmt.Sprintf("[Assistente]: %s\n", msg.Content))
				if msg.Query != "" {
					context.WriteString(fmt.Sprintf("  [Query executada]: %s\n", msg.Query))
				}
				if msg.Result != "" && msg.Result != "Nenhum resultado encontrado." {
					// Incluir resultado resumido para contexto
					resultLines := strings.Split(msg.Result, "\n")
					if len(resultLines) > 0 {
						context.WriteString(fmt.Sprintf("  [Resultado]: %s...\n", resultLines[0]))
					}
				}
			}
		}
		context.WriteString("\n")
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
// Agora é mais agressivo - sempre tenta executar quando faz sentido
func (c *DBChat) needsQueryExecution(message string) bool {
	message = strings.ToLower(strings.TrimSpace(message))
	
	// Se já é uma query SQL, sempre executar
	if strings.HasPrefix(message, "select") || 
	   strings.HasPrefix(message, "show") ||
	   strings.HasPrefix(message, "describe") ||
	   strings.HasPrefix(message, "desc") ||
	   strings.HasPrefix(message, "explain") {
		return true
	}
	
	// Palavras-chave que indicam necessidade de query (lista expandida)
	queryKeywords := []string{
		"mostre", "mostrar", "liste", "listar", "consulte", "consultar",
		"quantos", "quais", "quem", "onde", "quando", "como",
		"select", "query", "tabela", "tabelas", "dados", "informações",
		"estatísticas", "status", "verificar", "analisar", "buscar",
		"encontrar", "exibir", "veja", "veja-me", "me mostre", "me liste",
		"conte", "conte-me", "diga", "diga-me", "informe", "informe-me",
		"existe", "existem", "há", "tem", "tenho", "temos",
		"versão", "versao", "tamanho", "tamanhos", "conexões", "conexoes",
		"usuários", "usuarios", "schemas", "índices", "indices",
	}

	// Verificar palavras-chave
	for _, keyword := range queryKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}

	// Se a mensagem contém perguntas diretas (começa com palavras interrogativas)
	interrogatives := []string{"quantos", "quais", "quem", "onde", "quando", "como", "qual", "que"}
	for _, inter := range interrogatives {
		if strings.HasPrefix(message, inter) {
			return true
		}
	}

	// Se não encontrou nada específico, mas a mensagem parece ser uma solicitação de dados
	// (contém palavras relacionadas a banco de dados)
	dbWords := []string{"tabela", "coluna", "registro", "linha", "banco", "database", "schema"}
	for _, word := range dbWords {
		if strings.Contains(message, word) {
			return true
		}
	}

	return false
}

// generateQuery gera uma query SQL baseada na mensagem do usuário
func (c *DBChat) generateQuery(userMessage, context string) (string, error) {
	// Determinar sintaxe SQL baseada no tipo de banco
	limitClause := "LIMIT"
	if c.dbType == "sqlserver" {
		limitClause = "TOP"
	} else if c.dbType == "oracle" {
		limitClause = "ROWNUM"
	}

	prompt := fmt.Sprintf(`%s

O usuário solicitou: "%s"

Gere uma query SQL apropriada para %s que atenda EXATAMENTE à solicitação do usuário.

IMPORTANTE:
- Retorne APENAS a query SQL, sem explicações, sem markdown, sem código de bloco
- Use sintaxe correta para %s
- Seja específico e preciso
- Inclua apenas colunas necessárias
- Use %s quando apropriado para evitar resultados muito grandes (máximo 100 linhas)
- Se o usuário não especificar limite, use %s 100
- Use nomes de tabelas e colunas corretos baseados no contexto da conversa anterior

Query SQL:`, context, userMessage, c.dbType, c.dbType, limitClause, limitClause)

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

	// Limpar a query (remover markdown code blocks, explicações, etc.)
	query = strings.TrimSpace(query)
	
	// Remover blocos de código markdown
	query = strings.TrimPrefix(query, "```sql")
	query = strings.TrimPrefix(query, "```")
	query = strings.TrimSuffix(query, "```")
	query = strings.TrimSpace(query)
	
	// Remover explicações que possam vir antes ou depois da query
	lines := strings.Split(query, "\n")
	var queryLines []string
	inQuery := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Detectar início da query SQL
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "select") || 
		   strings.HasPrefix(lowerLine, "show") ||
		   strings.HasPrefix(lowerLine, "describe") ||
		   strings.HasPrefix(lowerLine, "desc") ||
		   strings.HasPrefix(lowerLine, "explain") ||
		   strings.HasPrefix(lowerLine, "with") {
			inQuery = true
		}
		if inQuery {
			queryLines = append(queryLines, line)
			// Parar em ponto e vírgula ou se encontrar explicação após
			if strings.HasSuffix(line, ";") {
				break
			}
		}
	}
	
	if len(queryLines) > 0 {
		query = strings.Join(queryLines, " ")
		query = strings.TrimSuffix(query, ";")
		query = strings.TrimSpace(query)
	}

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

	// Formato melhorado: tabela markdown
	result.WriteString("| ")
	result.WriteString(strings.Join(columns, " | "))
	result.WriteString(" |\n")
	
	// Separador
	result.WriteString("|")
	for range columns {
		result.WriteString(" --- |")
	}
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

		result.WriteString("| ")
		for i, val := range values {
			if val != nil {
				valStr := fmt.Sprintf("%v", val)
				// Escapar pipes em valores
				valStr = strings.ReplaceAll(valStr, "|", "\\|")
				result.WriteString(valStr)
			} else {
				result.WriteString("NULL")
			}
			if i < len(values)-1 {
				result.WriteString(" | ")
			}
		}
		result.WriteString(" |\n")

		count++
		if count >= 100 { // Limitar resultados
			result.WriteString(fmt.Sprintf("\n*... e mais resultados (limitado a 100 linhas)*\n"))
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
	prompt := fmt.Sprintf("%s\n\nO usuário perguntou: \"%s\"\n\nA seguinte query foi executada automaticamente:\n```sql\n%s\n```\n\nResultado obtido:\n```\n%s\n```\n\nIMPORTANTE:\n- Você DEVE responder baseado nos resultados REAIS obtidos da query\n- Formate a resposta de forma clara, natural e bem estruturada\n- Use os dados reais para responder a pergunta do usuário\n- Se houver tabelas ou listas, formate-as de forma legível\n- Forneça insights relevantes baseados nos dados obtidos\n- Seja direto e objetivo, mas completo\n- Use formatação markdown para melhorar a legibilidade (tabelas, listas, etc.)\n\nResponda de forma natural e bem formatada:", context, userMessage, query, result)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf("Você é um analista de dados especializado em %s que interpreta resultados de queries SQL de forma clara, natural e bem formatada. Você SEMPRE responde baseado nos dados reais obtidos, não em sugestões ou exemplos.", c.dbType),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 2000, 0.7)
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

