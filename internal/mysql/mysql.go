package mysql

import (
	"database/sql"
	"fmt"
	"strings"
)

// MySQLAnalyzer realiza análises específicas do MySQL
type MySQLAnalyzer struct {
	db *sql.DB
}

// NewMySQLAnalyzer cria um novo analisador MySQL
func NewMySQLAnalyzer(db *sql.DB) *MySQLAnalyzer {
	return &MySQLAnalyzer{db: db}
}

// AnalyzeReplication analisa replicação MySQL
func (m *MySQLAnalyzer) AnalyzeReplication() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Replicação MySQL\n\n")

	// Status do slave
	query := `SHOW SLAVE STATUS`

	rows, err := m.db.Query(query)
	if err != nil {
		// Pode ser master ou não ter replicação configurada
		return m.analyzeMasterStatus(result)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	result.WriteString("## Status de Replicação (Slave)\n\n")

	if rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return "", err
		}

		result.WriteString("| Parâmetro | Valor |\n")
		result.WriteString("|-----------|-------|\n")

		for i, col := range cols {
			val := "NULL"
			if values[i] != nil {
				val = fmt.Sprintf("%v", values[i])
			}
			result.WriteString(fmt.Sprintf("| %s | %s |\n", col, val))
		}

		// Verificar lag
		if len(cols) > 0 {
			secondsBehindMaster := getValueByName(cols, values, "Seconds_Behind_Master")
			if secondsBehindMaster != "" && secondsBehindMaster != "NULL" {
				result.WriteString(fmt.Sprintf("\n**Lag de Replicação:** %s segundos\n", secondsBehindMaster))
			}
		}
	}

	return result.String(), nil
}

// analyzeMasterStatus analisa status do master
func (m *MySQLAnalyzer) analyzeMasterStatus(result strings.Builder) (string, error) {
	query := `SHOW MASTER STATUS`

	rows, err := m.db.Query(query)
	if err != nil {
		result.WriteString("⚠️ Replicação não configurada ou sem permissões\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("## Status de Replicação (Master)\n\n")

	if rows.Next() {
		var file, position sql.NullString
		var binlogDoDB, binlogIgnoreDB sql.NullString
		var executedGtidSet sql.NullString

		err := rows.Scan(&file, &position, &binlogDoDB, &binlogIgnoreDB, &executedGtidSet)
		if err == nil {
			result.WriteString("| Parâmetro | Valor |\n")
			result.WriteString("|-----------|-------|\n")
			result.WriteString(fmt.Sprintf("| File | %s |\n", getString(file)))
			result.WriteString(fmt.Sprintf("| Position | %s |\n", getString(position)))
			result.WriteString(fmt.Sprintf("| Binlog Do DB | %s |\n", getString(binlogDoDB)))
			result.WriteString(fmt.Sprintf("| Binlog Ignore DB | %s |\n", getString(binlogIgnoreDB)))
		}
	}

	return result.String(), nil
}

// AnalyzeLocks analisa locks MySQL
func (m *MySQLAnalyzer) AnalyzeLocks() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Locks MySQL\n\n")

	// Processos e locks
	query := `
		SELECT 
			id,
			user,
			host,
			db,
			command,
			time,
			state,
			info
		FROM information_schema.PROCESSLIST
		WHERE command != 'Sleep' OR time > 0
		ORDER BY time DESC
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar processos: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Processos Ativos\n\n")
	result.WriteString("| ID | Usuário | Host | DB | Command | Tempo (s) | Estado | Query |\n")
	result.WriteString("|----|---------|------|----|---------|----------|--------|-------|\n")

	for rows.Next() {
		var id, time int
		var user, host, db, command, state sql.NullString
		var info sql.NullString

		err := rows.Scan(&id, &user, &host, &db, &command, &time, &state, &info)
		if err != nil {
			continue
		}

		queryText := "N/A"
		if info.Valid && len(info.String) > 50 {
			queryText = info.String[:50] + "..."
		} else if info.Valid {
			queryText = info.String
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %d | %s | %s |\n",
			id, getString(user), getString(host), getString(db),
			getString(command), time, getString(state), queryText))
	}

	// InnoDB Locks
	result.WriteString("\n## Locks InnoDB\n\n")
	innodbLocks, err := m.analyzeInnoDBLocks()
	if err == nil {
		result.WriteString(innodbLocks)
	}

	return result.String(), nil
}

// analyzeInnoDBLocks analisa locks InnoDB
func (m *MySQLAnalyzer) analyzeInnoDBLocks() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			r.trx_id waiting_trx_id,
			r.trx_mysql_thread_id waiting_thread,
			r.trx_query waiting_query,
			b.trx_id blocking_trx_id,
			b.trx_mysql_thread_id blocking_thread,
			b.trx_query blocking_query
		FROM information_schema.INNODB_LOCK_WAITS w
		INNER JOIN information_schema.INNODB_TRX b ON b.trx_id = w.blocking_trx_id
		INNER JOIN information_schema.INNODB_TRX r ON r.trx_id = w.requesting_trx_id
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Bloqueado (Thread) | Bloqueador (Thread) | Query Bloqueada | Query Bloqueadora |\n")
	result.WriteString("|--------------------|---------------------|-----------------|-------------------|\n")

	for rows.Next() {
		var waitingTrxID, waitingThread, blockingTrxID, blockingThread int64
		var waitingQuery, blockingQuery sql.NullString

		err := rows.Scan(&waitingTrxID, &waitingThread, &waitingQuery,
			&blockingTrxID, &blockingThread, &blockingQuery)
		if err != nil {
			continue
		}

		waitingQ := "N/A"
		if waitingQuery.Valid && len(waitingQuery.String) > 50 {
			waitingQ = waitingQuery.String[:50] + "..."
		} else if waitingQuery.Valid {
			waitingQ = waitingQuery.String
		}

		blockingQ := "N/A"
		if blockingQuery.Valid && len(blockingQuery.String) > 50 {
			blockingQ = blockingQuery.String[:50] + "..."
		} else if blockingQuery.Valid {
			blockingQ = blockingQuery.String
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %s | %s |\n",
			waitingThread, blockingThread, waitingQ, blockingQ))
	}

	return result.String(), nil
}

// AnalyzeFragmentation analisa fragmentação MySQL
func (m *MySQLAnalyzer) AnalyzeFragmentation() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Fragmentação MySQL\n\n")

	query := `
		SELECT 
			table_schema,
			table_name,
			ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'size_mb',
			ROUND((data_free / 1024 / 1024), 2) AS 'free_mb',
			ROUND((data_free / (data_length + index_length + data_free)) * 100, 2) AS 'frag_percent',
			table_rows,
			avg_row_length
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
		AND data_free > 0
		ORDER BY data_free DESC
		LIMIT 20
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar fragmentação: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Tabelas com Fragmentação\n\n")
	result.WriteString("| Schema | Tabela | Tamanho (MB) | Espaço Livre (MB) | % Fragmentação | Linhas |\n")
	result.WriteString("|--------|--------|---------------|-------------------|----------------|--------|\n")

	for rows.Next() {
		var tableSchema, tableName string
		var sizeMB, freeMB, fragPercent sql.NullFloat64
		var tableRows, avgRowLength sql.NullInt64

		err := rows.Scan(&tableSchema, &tableName, &sizeMB, &freeMB, &fragPercent, &tableRows, &avgRowLength)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %.2f | %.2f | %.2f%% | %d |\n",
			tableSchema, tableName, getFloat64(sizeMB), getFloat64(freeMB),
			getFloat64(fragPercent), getInt64(tableRows)))
	}

	return result.String(), nil
}

// GetExecutionPlan obtém plano de execução
func (m *MySQLAnalyzer) GetExecutionPlan(query string) (string, error) {
	var result strings.Builder

	result.WriteString("# Plano de Execução MySQL\n\n")

	explainQuery := "EXPLAIN " + query

	rows, err := m.db.Query(explainQuery)
	if err != nil {
		return "", fmt.Errorf("erro ao obter plano de execução: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	result.WriteString("| " + strings.Join(cols, " | ") + " |\n")
	result.WriteString("|" + strings.Repeat("---|", len(cols)) + "\n")

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			continue
		}

		var row []string
		for _, val := range values {
			if val == nil {
				row = append(row, "NULL")
			} else {
				row = append(row, fmt.Sprintf("%v", val))
			}
		}
		result.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}

	return result.String(), nil
}

func getString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "N/A"
}

func getInt64(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}
	return 0
}

func getFloat64(f sql.NullFloat64) float64 {
	if f.Valid {
		return f.Float64
	}
	return 0
}

func getValueByName(cols []string, values []interface{}, name string) string {
	for i, col := range cols {
		if col == name {
			if values[i] != nil {
				return fmt.Sprintf("%v", values[i])
			}
			return "NULL"
		}
	}
	return ""
}

