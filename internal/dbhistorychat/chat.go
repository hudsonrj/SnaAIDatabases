package dbhistorychat

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
)

// ChatSession representa uma sessão de chat com o histórico de análises
type ChatSession struct {
	ID        int
	Messages  []ChatMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChatMessage representa uma mensagem no chat
type ChatMessage struct {
	Role      string // "user" ou "assistant"
	Content   string
	Query     string // Query SQL gerada (se aplicável)
	Result    string // Resultado da query (se aplicável)
	Timestamp time.Time
}

// DBHistoryChat é o gerenciador de chat com o histórico de análises
type DBHistoryChat struct {
	aiClient ai.AIClient
	db       *sql.DB
	session  *ChatSession
}

// NewDBHistoryChat cria uma nova sessão de chat com o histórico
func NewDBHistoryChat(db *sql.DB) (*DBHistoryChat, error) {
	aiClient, err := ai.NewAIClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}

	session := &ChatSession{
		Messages:  []ChatMessage{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &DBHistoryChat{
		aiClient: aiClient,
		db:       db,
		session:  session,
	}, nil
}

// SendMessage envia uma mensagem e recebe resposta
func (c *DBHistoryChat) SendMessage(userMessage string) (string, error) {
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
func (c *DBHistoryChat) buildContext() string {
	var context strings.Builder

	context.WriteString("Você é um assistente especializado em análise de bancos de dados.\n\n")
	context.WriteString("Você está conversando com um banco de dados SQLite que armazena todas as análises de bancos de dados realizadas.\n\n")
	
	context.WriteString("Schema da tabela db_analyses:\n")
	context.WriteString("- id: INTEGER PRIMARY KEY\n")
	context.WriteString("- title: TEXT (título da análise)\n")
	context.WriteString("- database_type: TEXT (oracle, sqlserver, mysql, postgresql, mongodb)\n")
	context.WriteString("- analysis_type: TEXT (diagnostic, tuning, query, awr, ash, locks, etc.)\n")
	context.WriteString("- connection_config: TEXT (JSON com configuração de conexão)\n")
	context.WriteString("- log_file_path: TEXT (caminho do arquivo de log, se aplicável)\n")
	context.WriteString("- output_type: TEXT (text, json, markdown)\n")
	context.WriteString("- result: TEXT (resultado completo da análise)\n")
	context.WriteString("- ai_insights: TEXT (insights gerados pela IA)\n")
	context.WriteString("- status: TEXT (pending, completed, error)\n")
	context.WriteString("- error_message: TEXT (mensagem de erro, se houver)\n")
	context.WriteString("- created_at: DATETIME (data de criação)\n")
	context.WriteString("- updated_at: DATETIME (data de atualização)\n\n")

	context.WriteString("Você pode ajudar o usuário a:\n")
	context.WriteString("- Listar análises por tipo de banco, tipo de análise, data, etc.\n")
	context.WriteString("- Comparar análises de diferentes datas para ver evolução\n")
	context.WriteString("- Identificar problemas e insights das análises\n")
	context.WriteString("- Rastrear a evolução ou degradação dos bancos ao longo do tempo\n")
	context.WriteString("- Analisar tendências e padrões nas análises\n")
	context.WriteString("- Responder perguntas sobre resultados específicos\n\n")

	// Adicionar TODO o histórico da conversa (mantém contexto completo)
	if len(c.session.Messages) > 0 {
		context.WriteString("Histórico completo da conversa:\n")
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
func (c *DBHistoryChat) generateResponse(userMessage, context string) (response, query, result string, err error) {
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
func (c *DBHistoryChat) needsQueryExecution(message string) bool {
	message = strings.ToLower(strings.TrimSpace(message))
	
	// Se já é uma query SQL, sempre executar
	if strings.HasPrefix(message, "select") || 
	   strings.HasPrefix(message, "show") ||
	   strings.HasPrefix(message, "describe") ||
	   strings.HasPrefix(message, "desc") ||
	   strings.HasPrefix(message, "explain") {
		return true
	}
	
	// Palavras-chave que indicam necessidade de query
	queryKeywords := []string{
		"mostre", "mostrar", "liste", "listar", "consulte", "consultar",
		"quantos", "quais", "quem", "onde", "quando", "como",
		"select", "query", "tabela", "tabelas", "dados", "informações",
		"estatísticas", "status", "verificar", "analisar", "buscar",
		"encontrar", "exibir", "veja", "veja-me", "me mostre", "me liste",
		"conte", "conte-me", "diga", "diga-me", "informe", "informe-me",
		"existe", "existem", "há", "tem", "tenho", "temos",
		"análises", "analises", "histórico", "historico", "resultados",
		"comparar", "comparação", "evolução", "evolucao", "tendências", "tendencias",
		"problemas", "insights", "gráficos", "graficos", "datas", "coletas",
		"bancos", "databases", "tipos", "datas de coleta",
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
	// (contém palavras relacionadas a análises)
	analysisWords := []string{"análise", "analise", "resultado", "insight", "problema", "erro", "status"}
	for _, word := range analysisWords {
		if strings.Contains(message, word) {
			return true
		}
	}

	return false
}

// generateQuery gera uma query SQL baseada na mensagem do usuário
func (c *DBHistoryChat) generateQuery(userMessage, context string) (string, error) {
	prompt := fmt.Sprintf(`%s

O usuário solicitou: "%s"

Gere uma query SQL apropriada para SQLite que atenda EXATAMENTE à solicitação do usuário.

IMPORTANTE:
- Retorne APENAS a query SQL, sem explicações, sem markdown, sem código de bloco
- Use sintaxe SQLite correta
- A tabela se chama 'db_analyses'
- Seja específico e preciso
- Inclua apenas colunas necessárias
- Use LIMIT quando apropriado para evitar resultados muito grandes (máximo 100 linhas)
- Se o usuário não especificar limite, use LIMIT 100
- Para datas, use funções SQLite como date(), datetime(), strftime()
- Para comparações de texto, use LIKE ou = conforme apropriado
- Para análises de evolução, compare created_at entre diferentes períodos
- Para insights e problemas, busque em ai_insights e error_message

Query SQL:`, context, userMessage)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em SQL para SQLite. Gere queries SQL precisas e otimizadas para consultar análises de bancos de dados armazenadas na tabela db_analyses.",
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

// executeQuery executa a query no banco de dados SQLite
func (c *DBHistoryChat) executeQuery(query string) (string, error) {
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
				// Truncar valores muito longos (especialmente result e ai_insights)
				if len(valStr) > 200 {
					valStr = valStr[:200] + "..."
				}
				// Escapar pipes em valores
				valStr = strings.ReplaceAll(valStr, "|", "\\|")
				valStr = strings.ReplaceAll(valStr, "\n", " ")
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
func (c *DBHistoryChat) interpretQueryResult(userMessage, context, query, result string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nO usuário perguntou: \"%s\"\n\nA seguinte query foi executada automaticamente:\n```sql\n%s\n```\n\nResultado obtido:\n```\n%s\n```\n\nIMPORTANTE:\n- Você DEVE responder baseado nos resultados REAIS obtidos da query\n- Formate a resposta de forma clara, natural e bem estruturada\n- Use os dados reais para responder a pergunta do usuário\n- Se houver tabelas ou listas, formate-as de forma legível\n- Forneça insights relevantes baseados nos dados obtidos\n- Seja direto e objetivo, mas completo\n- Use formatação markdown para melhorar a legibilidade (tabelas, listas, etc.)\n- Para comparações de datas, destaque tendências e evoluções\n- Para problemas e erros, destaque-os claramente\n- Para insights, apresente-os de forma organizada\n\nResponda de forma natural e bem formatada:", context, userMessage, query, result)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um analista especializado em interpretar resultados de análises de bancos de dados. Você interpreta resultados de queries SQL de forma clara, natural e bem formatada. Você SEMPRE responde baseado nos dados reais obtidos, não em sugestões ou exemplos. Você ajuda a identificar tendências, problemas, evoluções e insights importantes.",
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
func (c *DBHistoryChat) generateErrorExplanation(userMessage, context, query, errorMsg string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nO usuário tentou executar a seguinte query:\n```sql\n%s\n```\n\nMas ocorreu um erro: %s\n\nExplique o erro de forma clara e sugira como corrigir.", context, query, errorMsg)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em SQL que ajuda a entender e corrigir erros em queries SQLite.",
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
func (c *DBHistoryChat) generateConversationalResponse(userMessage, context string) (string, error) {
	prompt := fmt.Sprintf(`%s

O usuário disse: "%s"

Responda de forma útil e conversacional sobre o histórico de análises de bancos de dados.`, context, userMessage)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um assistente especializado em análise de bancos de dados. Responda perguntas sobre o histórico de análises de forma clara e útil.",
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
func (c *DBHistoryChat) GetSession() *ChatSession {
	return c.session
}

// Close fecha a conexão (não fecha o DB pois é compartilhado)
func (c *DBHistoryChat) Close() error {
	// Não fechar o DB aqui pois pode ser compartilhado
	return nil
}

