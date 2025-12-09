package dbdynamic

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/dbtypes"
)

// DynamicAnalyzer realiza análises dinâmicas gerando queries com IA
type DynamicAnalyzer struct {
	aiClient *ai.GroqClient
	db       *sql.DB
	dbType   dbtypes.DatabaseType
	config   *dbtypes.ConnectionConfig
}

// NewDynamicAnalyzer cria um novo analisador dinâmico
func NewDynamicAnalyzer(dbType dbtypes.DatabaseType, config *dbtypes.ConnectionConfig, db *sql.DB) (*DynamicAnalyzer, error) {
	aiClient, err := ai.NewGroqClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}

	return &DynamicAnalyzer{
		aiClient: aiClient,
		db:       db,
		dbType:   dbType,
		config:   config,
	}, nil
}

// AnalyzeWithNaturalLanguage realiza análise baseada em linguagem natural
func (d *DynamicAnalyzer) AnalyzeWithNaturalLanguage(request string) (string, error) {
	// Obter informações do schema para contexto
	schemaInfo, err := d.getSchemaInfo()
	if err != nil {
		// Continuar mesmo sem schema info
		schemaInfo = "Informações de schema não disponíveis"
	}

	// Gerar query com IA
	query, err := d.generateAnalysisQuery(request, schemaInfo)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar query: %w", err)
	}

	// Executar query
	result, err := d.executeQuery(query)
	if err != nil {
		return "", fmt.Errorf("erro ao executar query: %w", err)
	}

	// Interpretar resultados com IA
	interpretation, err := d.interpretResults(request, query, result)
	if err != nil {
		// Retornar resultado mesmo sem interpretação
		return fmt.Sprintf("Query executada:\n```sql\n%s\n```\n\nResultados:\n%s", query, result), nil
	}

	// Combinar tudo
	var output strings.Builder
	output.WriteString("# Análise Dinâmica\n\n")
	output.WriteString(fmt.Sprintf("**Solicitação:** %s\n\n", request))
	output.WriteString(fmt.Sprintf("**Query Gerada:**\n```sql\n%s\n```\n\n", query))
	output.WriteString(fmt.Sprintf("**Resultados:**\n```\n%s\n```\n\n", result))
	output.WriteString(fmt.Sprintf("**Interpretação da IA:**\n%s\n", interpretation))

	return output.String(), nil
}

// getSchemaInfo obtém informações básicas do schema
func (d *DynamicAnalyzer) getSchemaInfo() (string, error) {
	if d.db == nil {
		return "", fmt.Errorf("conexão não disponível")
	}

	var schemaInfo strings.Builder

	switch d.dbType {
	case dbtypes.DatabaseTypePostgreSQL:
		query := `
			SELECT table_schema, table_name, column_name, data_type
			FROM information_schema.columns
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name, ordinal_position
			LIMIT 100
		`
		rows, err := d.db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		schemaInfo.WriteString("Schema PostgreSQL:\n")
		for rows.Next() {
			var schema, table, column, dataType string
			if err := rows.Scan(&schema, &table, &column, &dataType); err == nil {
				schemaInfo.WriteString(fmt.Sprintf("- %s.%s.%s (%s)\n", schema, table, column, dataType))
			}
		}

	case dbtypes.DatabaseTypeMySQL:
		query := `
			SELECT table_schema, table_name, column_name, data_type
			FROM information_schema.columns
			WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
			ORDER BY table_schema, table_name, ordinal_position
			LIMIT 100
		`
		rows, err := d.db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		schemaInfo.WriteString("Schema MySQL:\n")
		for rows.Next() {
			var schema, table, column, dataType string
			if err := rows.Scan(&schema, &table, &column, &dataType); err == nil {
				schemaInfo.WriteString(fmt.Sprintf("- %s.%s.%s (%s)\n", schema, table, column, dataType))
			}
		}

	case dbtypes.DatabaseTypeSQLServer:
		query := `
			SELECT TOP 100
				TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, DATA_TYPE
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA NOT IN ('sys', 'information_schema')
			ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION
		`
		rows, err := d.db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		schemaInfo.WriteString("Schema SQL Server:\n")
		for rows.Next() {
			var schema, table, column, dataType string
			if err := rows.Scan(&schema, &table, &column, &dataType); err == nil {
				schemaInfo.WriteString(fmt.Sprintf("- %s.%s.%s (%s)\n", schema, table, column, dataType))
			}
		}

	case dbtypes.DatabaseTypeOracle:
		query := `
			SELECT owner, table_name, column_name, data_type
			FROM all_tab_columns
			WHERE owner NOT IN ('SYS', 'SYSTEM', 'SYSAUX')
			AND ROWNUM <= 100
			ORDER BY owner, table_name, column_id
		`
		rows, err := d.db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		schemaInfo.WriteString("Schema Oracle:\n")
		for rows.Next() {
			var owner, table, column, dataType string
			if err := rows.Scan(&owner, &table, &column, &dataType); err == nil {
				schemaInfo.WriteString(fmt.Sprintf("- %s.%s.%s (%s)\n", owner, table, column, dataType))
			}
		}
	}

	return schemaInfo.String(), nil
}

// generateAnalysisQuery gera query SQL baseada na solicitação
func (d *DynamicAnalyzer) generateAnalysisQuery(request, schemaInfo string) (string, error) {
	prompt := fmt.Sprintf(`Você é um especialista em SQL para %s.

Schema disponível:
%s

O usuário solicitou: "%s"

Gere uma query SQL apropriada que atenda à solicitação. 

IMPORTANTE:
- Retorne APENAS a query SQL, sem explicações
- Use sintaxe correta para %s
- Seja específico e preciso
- Use LIMIT/TOP/ROWNUM quando apropriado
- Se não souber a estrutura exata, use nomes genéricos comuns

Query SQL:`, d.dbType, schemaInfo, request, d.dbType)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf("Você é um especialista em SQL para %s. Gere queries SQL precisas e otimizadas.", d.dbType),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	query, err := d.aiClient.Chat(messages, 1000, 0.3)
	if err != nil {
		return "", err
	}

	// Limpar a query
	query = strings.TrimSpace(query)
	query = strings.TrimPrefix(query, "```sql")
	query = strings.TrimPrefix(query, "```")
	query = strings.TrimSuffix(query, "```")
	query = strings.TrimSpace(query)

	return query, nil
}

// executeQuery executa a query
func (d *DynamicAnalyzer) executeQuery(query string) (string, error) {
	if d.db == nil {
		return "", fmt.Errorf("conexão não disponível")
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("erro ao obter colunas: %w", err)
	}

	var result strings.Builder

	// Cabeçalho
	result.WriteString(strings.Join(columns, " | "))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("-", len(strings.Join(columns, " | "))))
	result.WriteString("\n")

	// Valores
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
		if count >= 1000 {
			result.WriteString("\n... e mais resultados (limitado a 1000)\n")
			break
		}
	}

	if count == 0 {
		return "Nenhum resultado encontrado.", nil
	}

	return result.String(), nil
}

// interpretResults interpreta os resultados com IA
func (d *DynamicAnalyzer) interpretResults(request, query, result string) (string, error) {
	prompt := fmt.Sprintf("O usuário solicitou: \"%s\"\n\nA seguinte query foi executada:\n```sql\n%s\n```\n\nResultados:\n```\n%s\n```\n\nInterprete os resultados de forma clara e útil. Explique o que os dados significam, identifique padrões, anomalias ou insights relevantes.", request, query, result)

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

	interpretation, err := d.aiClient.Chat(messages, 2000, 0.7)
	if err != nil {
		return "", err
	}

	return interpretation, nil
}

// Close fecha a conexão
func (d *DynamicAnalyzer) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

