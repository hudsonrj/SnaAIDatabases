package sqlserver

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// SQLServerAnalyzer realiza análises específicas do SQL Server
type SQLServerAnalyzer struct {
	db *sql.DB
}

// NewSQLServerAnalyzer cria um novo analisador SQL Server
func NewSQLServerAnalyzer(db *sql.DB) *SQLServerAnalyzer {
	return &SQLServerAnalyzer{db: db}
}

// AnalyzeDMV analisa Dynamic Management Views (equivalente ao AWR do Oracle)
func (s *SQLServerAnalyzer) AnalyzeDMV(dmvType string, beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Análise DMV - %s\n\n", dmvType))

	switch dmvType {
	case "wait_stats":
		return s.analyzeWaitStats(beginTime, endTime)
	case "query_stats":
		return s.analyzeQueryStats(beginTime, endTime)
	case "session_stats":
		return s.analyzeSessionStats(beginTime, endTime)
	case "io_stats":
		return s.analyzeIOStats(beginTime, endTime)
	case "locks":
		return s.AnalyzeLocks()
	case "active_sessions":
		return s.AnalyzeActiveSessions()
	case "running_queries":
		return s.AnalyzeRunningQueries()
	default:
		return s.analyzeGeneralDMV(dmvType, beginTime, endTime)
	}
}

// analyzeWaitStats analisa estatísticas de wait
func (s *SQLServerAnalyzer) analyzeWaitStats(beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Wait Statistics\n\n")

	query := `
		SELECT 
			wait_type,
			waiting_tasks_count,
			wait_time_ms,
			max_wait_time_ms,
			signal_wait_time_ms
		FROM sys.dm_os_wait_stats
		WHERE wait_time_ms > 0
		ORDER BY wait_time_ms DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar wait stats: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Top Waits\n\n")
	result.WriteString("| Tipo de Wait | Tarefas | Tempo Total (ms) | Tempo Máximo (ms) |\n")
	result.WriteString("|--------------|---------|------------------|-------------------|\n")

	count := 0
	for rows.Next() {
		var waitType string
		var waitingTasks, waitTime, maxWaitTime, signalWaitTime int64

		err := rows.Scan(&waitType, &waitingTasks, &waitTime, &maxWaitTime, &signalWaitTime)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n",
			waitType, waitingTasks, waitTime, maxWaitTime))

		count++
		if count >= 20 {
			break
		}
	}

	return result.String(), nil
}

// analyzeQueryStats analisa estatísticas de queries
func (s *SQLServerAnalyzer) analyzeQueryStats(beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Query Statistics\n\n")

	query := `
		SELECT TOP 20
			qs.sql_handle,
			qs.plan_handle,
			qs.execution_count,
			qs.total_worker_time / 1000 as total_cpu_ms,
			qs.total_elapsed_time / 1000 as total_elapsed_ms,
			qs.total_logical_reads,
			qs.total_physical_reads,
			qs.creation_time,
			qs.last_execution_time
		FROM sys.dm_exec_query_stats qs
		ORDER BY qs.total_worker_time DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar query stats: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Top Queries por CPU\n\n")
	result.WriteString("| Execuções | CPU (ms) | Tempo Total (ms) | Logical Reads | Physical Reads |\n")
	result.WriteString("|-----------|----------|------------------|---------------|----------------|\n")

	for rows.Next() {
		var sqlHandle, planHandle string
		var execCount, totalCPU, totalElapsed, logicalReads, physicalReads int64
		var creationTime, lastExecTime time.Time

		err := rows.Scan(&sqlHandle, &planHandle, &execCount, &totalCPU, &totalElapsed,
			&logicalReads, &physicalReads, &creationTime, &lastExecTime)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %d | %d | %d |\n",
			execCount, totalCPU, totalElapsed, logicalReads, physicalReads))
	}

	return result.String(), nil
}

// analyzeSessionStats analisa estatísticas de sessões usando procedures nativas
func (s *SQLServerAnalyzer) analyzeSessionStats(beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Sessões (sp_who2)\n\n")

	// Usar sp_who2 para obter informações detalhadas de sessões
	query := `
		EXEC sp_who2
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao executar sp_who2: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Todas as Sessões\n\n")
	result.WriteString("| SPID | Status | Login | Host | BlkBy | DB | Command | CPU | DiskIO | LastBatch | Program |\n")
	result.WriteString("|------|--------|-------|------|-------|----|---------|-----|--------|-----------|---------|\n")

	// sp_who2 retorna colunas dinâmicas, precisamos ler como strings
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

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

		// Converter valores para strings
		rowData := make([]string, len(cols))
		for i, val := range values {
			if val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = "NULL"
			}
		}

		// Formatar linha (sp_who2 tem colunas específicas)
		if len(rowData) >= 11 {
			result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
				rowData[0], rowData[1], rowData[2], rowData[3], rowData[4],
				rowData[5], rowData[6], rowData[7], rowData[8], rowData[9], rowData[10]))
		}
	}

	// Análise adicional usando DMVs
	result.WriteString("\n## Sessões Detalhadas (DMV)\n\n")
	return s.analyzeSessionsDMV(result)
}

// analyzeSessionsDMV analisa sessões usando DMVs detalhadas
func (s *SQLServerAnalyzer) analyzeSessionsDMV(result strings.Builder) (string, error) {
	query := `
		SELECT 
			s.session_id,
			s.login_name,
			s.host_name,
			s.program_name,
			s.status,
			s.cpu_time,
			s.memory_usage,
			s.total_scheduled_time,
			s.total_elapsed_time,
			s.last_request_start_time,
			s.last_request_end_time,
			s.reads,
			s.writes,
			s.logical_reads
		FROM sys.dm_exec_sessions s
		WHERE s.is_user_process = 1
		ORDER BY s.cpu_time DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return result.String(), fmt.Errorf("erro ao consultar sessões: %w", err)
	}
	defer rows.Close()

	result.WriteString("| Session ID | Login | Host | Programa | Status | CPU (ms) | Reads | Writes | Logical Reads |\n")
	result.WriteString("|------------|-------|------|----------|--------|----------|-------|--------|---------------|\n")

	for rows.Next() {
		var sessionID, cpuTime, memoryUsage, scheduledTime, elapsedTime int
		var loginName, hostName, programName, status string
		var lastRequestStart, lastRequestEnd sql.NullTime
		var reads, writes, logicalReads int64

		err := rows.Scan(&sessionID, &loginName, &hostName, &programName, &status,
			&cpuTime, &memoryUsage, &scheduledTime, &elapsedTime,
			&lastRequestStart, &lastRequestEnd, &reads, &writes, &logicalReads)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %d | %d | %d | %d |\n",
			sessionID, loginName, hostName, programName, status, cpuTime, reads, writes, logicalReads))
	}

	return result.String(), nil
}

// AnalyzeLocks analisa locks usando procedures nativas
func (s *SQLServerAnalyzer) AnalyzeLocks() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Locks\n\n")

	// Usar sp_lock para obter informações de locks
	result.WriteString("## Locks (sp_lock)\n\n")
	lockQuery := `EXEC sp_lock`

	rows, err := s.db.Query(lockQuery)
	if err != nil {
		// Se sp_lock não funcionar, usar DMV
		return s.analyzeLocksDMV()
	}
	defer rows.Close()

	result.WriteString("| SPID | DBID | ObjID | IndId | Type | Resource | Mode | Status |\n")
	result.WriteString("|------|------|-------|-------|------|----------|------|--------|\n")

	cols, err := rows.Columns()
	if err == nil {
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

			rowData := make([]string, len(cols))
			for i, val := range values {
				if val != nil {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = "NULL"
				}
			}

			if len(rowData) >= 8 {
				result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
					rowData[0], rowData[1], rowData[2], rowData[3], rowData[4],
					rowData[5], rowData[6], rowData[7]))
			}
		}
	}

	// Análise adicional usando DMVs
	result.WriteString("\n## Locks Detalhados (DMV)\n\n")
	dmvResult, err := s.analyzeLocksDMV()
	if err == nil {
		result.WriteString(dmvResult)
	}

	return result.String(), nil
}

// analyzeLocksDMV analisa locks usando DMVs
func (s *SQLServerAnalyzer) analyzeLocksDMV() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			l.request_session_id,
			l.resource_database_id,
			l.resource_associated_entity_id,
			l.resource_type,
			l.resource_description,
			l.request_mode,
			l.request_status,
			OBJECT_NAME(p.object_id) as object_name,
			p.index_id,
			i.name as index_name
		FROM sys.dm_tran_locks l
		LEFT JOIN sys.partitions p ON l.resource_associated_entity_id = p.hobt_id
		LEFT JOIN sys.indexes i ON p.object_id = i.object_id AND p.index_id = i.index_id
		WHERE l.resource_database_id = DB_ID()
		ORDER BY l.request_session_id
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar locks: %w", err)
	}
	defer rows.Close()

	result.WriteString("| Session ID | Database | Object | Index | Resource Type | Mode | Status |\n")
	result.WriteString("|------------|----------|--------|-------|---------------|------|--------|\n")

	for rows.Next() {
		var sessionID, dbID, entityID int
		var resourceType, resourceDesc, requestMode, requestStatus string
		var objectName, indexName sql.NullString
		var indexID sql.NullInt64

		err := rows.Scan(&sessionID, &dbID, &entityID, &resourceType, &resourceDesc,
			&requestMode, &requestStatus, &objectName, &indexID, &indexName)
		if err != nil {
			continue
		}

		objName := "N/A"
		if objectName.Valid {
			objName = objectName.String
		}

		idxName := "N/A"
		if indexName.Valid {
			idxName = indexName.String
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %s | %s | %s | %s | %s |\n",
			sessionID, dbID, objName, idxName, resourceType, requestMode, requestStatus))
	}

	return result.String(), nil
}

// AnalyzeActiveSessions analisa sessões ativas com queries em execução (tipo SQL Profiler)
func (s *SQLServerAnalyzer) AnalyzeActiveSessions() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Sessões Ativas (Tipo SQL Profiler)\n\n")

	query := `
		SELECT 
			s.session_id,
			s.login_name,
			s.host_name,
			s.program_name,
			s.status,
			r.command,
			r.status as request_status,
			r.start_time,
			r.cpu_time,
			r.total_elapsed_time,
			r.reads,
			r.writes,
			r.logical_reads,
			t.text as sql_text,
			p.query_plan
		FROM sys.dm_exec_sessions s
		INNER JOIN sys.dm_exec_requests r ON s.session_id = r.session_id
		OUTER APPLY sys.dm_exec_sql_text(r.sql_handle) t
		OUTER APPLY sys.dm_exec_query_plan(r.plan_handle) p
		WHERE s.is_user_process = 1
		ORDER BY r.start_time DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar sessões ativas: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Queries em Execução\n\n")
	result.WriteString("| Session ID | Login | Programa | Command | Status | CPU (ms) | Elapsed (ms) | SQL Text |\n")
	result.WriteString("|------------|-------|----------|---------|--------|----------|--------------|----------|\n")

	for rows.Next() {
		var sessionID, cpuTime, elapsedTime int
		var loginName, hostName, programName, status, command, requestStatus string
		var startTime sql.NullTime
		var reads, writes, logicalReads int64
		var sqlText sql.NullString
		var queryPlan sql.NullString

		err := rows.Scan(&sessionID, &loginName, &hostName, &programName, &status,
			&command, &requestStatus, &startTime, &cpuTime, &elapsedTime,
			&reads, &writes, &logicalReads, &sqlText, &queryPlan)
		if err != nil {
			continue
		}

		sqlTextShort := "N/A"
		if sqlText.Valid {
			if len(sqlText.String) > 100 {
				sqlTextShort = sqlText.String[:100] + "..."
			} else {
				sqlTextShort = sqlText.String
			}
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %d | %d | %s |\n",
			sessionID, loginName, programName, command, requestStatus, cpuTime, elapsedTime, sqlTextShort))
	}

	// Adicionar análise de bloqueios entre sessões
	result.WriteString("\n## Bloqueios entre Sessões\n\n")
	blockingResult, err := s.analyzeBlockingSessions()
	if err == nil {
		result.WriteString(blockingResult)
	}

	return result.String(), nil
}

// analyzeBlockingSessions analisa sessões bloqueadas e bloqueadoras
func (s *SQLServerAnalyzer) analyzeBlockingSessions() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			blocking.session_id AS blocking_session_id,
			blocking.login_name AS blocking_login,
			blocking.program_name AS blocking_program,
			blocked.session_id AS blocked_session_id,
			blocked.login_name AS blocked_login,
			blocked.program_name AS blocked_program,
			blocked.wait_type,
			blocked.wait_time,
			blocked.wait_resource,
			t.text AS blocked_sql_text
		FROM sys.dm_exec_sessions blocking
		INNER JOIN sys.dm_exec_requests blocked ON blocking.session_id = blocked.blocking_session_id
		OUTER APPLY sys.dm_exec_sql_text(blocked.sql_handle) t
		WHERE blocking.is_user_process = 1
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar bloqueios: %w", err)
	}
	defer rows.Close()

	result.WriteString("| Bloqueando (Session) | Bloqueado (Session) | Wait Type | Wait Time (ms) | Wait Resource | SQL Text |\n")
	result.WriteString("|----------------------|---------------------|-----------|----------------|---------------|----------|\n")

	for rows.Next() {
		var blockingSessionID, blockedSessionID, waitTime int
		var blockingLogin, blockingProgram, blockedLogin, blockedProgram string
		var waitType, waitResource sql.NullString
		var blockedSQLText sql.NullString

		err := rows.Scan(&blockingSessionID, &blockingLogin, &blockingProgram,
			&blockedSessionID, &blockedLogin, &blockedProgram,
			&waitType, &waitTime, &waitResource, &blockedSQLText)
		if err != nil {
			continue
		}

		sqlText := "N/A"
		if blockedSQLText.Valid && len(blockedSQLText.String) > 50 {
			sqlText = blockedSQLText.String[:50] + "..."
		} else if blockedSQLText.Valid {
			sqlText = blockedSQLText.String
		}

		result.WriteString(fmt.Sprintf("| %d (%s) | %d (%s) | %s | %d | %s | %s |\n",
			blockingSessionID, blockingLogin, blockedSessionID, blockedLogin,
			getString(waitType), waitTime, getString(waitResource), sqlText))
	}

	return result.String(), nil
}

// AnalyzeRunningQueries analisa queries em execução com planos
func (s *SQLServerAnalyzer) AnalyzeRunningQueries() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Queries em Execução com Planos\n\n")

	query := `
		SELECT 
			r.session_id,
			r.request_id,
			r.start_time,
			r.status,
			r.command,
			r.cpu_time,
			r.total_elapsed_time,
			r.reads,
			r.writes,
			r.logical_reads,
			SUBSTRING(t.text, (r.statement_start_offset/2)+1,
				((CASE r.statement_end_offset
					WHEN -1 THEN DATALENGTH(t.text)
					ELSE r.statement_end_offset
				END - r.statement_start_offset)/2) + 1) AS statement_text,
			p.query_plan,
			s.login_name,
			s.program_name,
			s.host_name
		FROM sys.dm_exec_requests r
		CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) t
		OUTER APPLY sys.dm_exec_query_plan(r.plan_handle) p
		INNER JOIN sys.dm_exec_sessions s ON r.session_id = s.session_id
		WHERE r.status IN ('running', 'runnable', 'suspended')
		ORDER BY r.total_elapsed_time DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar queries em execução: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Queries Ativas\n\n")

	count := 0
	for rows.Next() {
		var sessionID, requestID, cpuTime, elapsedTime int
		var startTime time.Time
		var status, command string
		var reads, writes, logicalReads int64
		var statementText, queryPlan sql.NullString
		var loginName, programName, hostName string

		err := rows.Scan(&sessionID, &requestID, &startTime, &status, &command,
			&cpuTime, &elapsedTime, &reads, &writes, &logicalReads,
			&statementText, &queryPlan, &loginName, &programName, &hostName)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("### Query %d (Session: %d)\n\n", requestID, sessionID))
		result.WriteString(fmt.Sprintf("- **Login:** %s\n", loginName))
		result.WriteString(fmt.Sprintf("- **Programa:** %s\n", programName))
		result.WriteString(fmt.Sprintf("- **Host:** %s\n", hostName))
		result.WriteString(fmt.Sprintf("- **Status:** %s\n", status))
		result.WriteString(fmt.Sprintf("- **Command:** %s\n", command))
		result.WriteString(fmt.Sprintf("- **Início:** %s\n", startTime.Format("2006-01-02 15:04:05")))
		result.WriteString(fmt.Sprintf("- **CPU Time:** %d ms\n", cpuTime))
		result.WriteString(fmt.Sprintf("- **Elapsed Time:** %d ms\n", elapsedTime))
		result.WriteString(fmt.Sprintf("- **Reads:** %d | **Writes:** %d | **Logical Reads:** %d\n\n", reads, writes, logicalReads))

		if statementText.Valid {
			result.WriteString("**SQL Statement:**\n```sql\n")
			if len(statementText.String) > 500 {
				result.WriteString(statementText.String[:500] + "...\n")
			} else {
				result.WriteString(statementText.String + "\n")
			}
			result.WriteString("```\n\n")
		}

		if queryPlan.Valid {
			result.WriteString("**Query Plan:**\n```xml\n")
			if len(queryPlan.String) > 1000 {
				result.WriteString(queryPlan.String[:1000] + "...\n")
			} else {
				result.WriteString(queryPlan.String + "\n")
			}
			result.WriteString("```\n\n")
		}

		result.WriteString("---\n\n")

		count++
		if count >= 10 {
			result.WriteString(fmt.Sprintf("\n... e mais queries (limitado a 10)\n"))
			break
		}
	}

	return result.String(), nil
}

func getString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "N/A"
}

// analyzeIOStats analisa estatísticas de I/O
func (s *SQLServerAnalyzer) analyzeIOStats(beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de I/O Statistics\n\n")

	query := `
		SELECT 
			database_id,
			file_id,
			io_stall_read_ms,
			io_stall_write_ms,
			io_stall,
			num_of_reads,
			num_of_writes,
			bytes_read,
			bytes_written
		FROM sys.dm_io_virtual_file_stats(NULL, NULL)
		ORDER BY io_stall DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar I/O stats: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Estatísticas de I/O por Arquivo\n\n")
	result.WriteString("| Database ID | File ID | Reads | Writes | Stall (ms) | Bytes Read | Bytes Written |\n")
	result.WriteString("|-------------|---------|-------|--------|------------|------------|--------------|\n")

	for rows.Next() {
		var dbID, fileID, numReads, numWrites int64
		var ioStallRead, ioStallWrite, ioStall, bytesRead, bytesWritten int64

		err := rows.Scan(&dbID, &fileID, &ioStallRead, &ioStallWrite, &ioStall,
			&numReads, &numWrites, &bytesRead, &bytesWritten)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %d | %d | %d | %d | %d |\n",
			dbID, fileID, numReads, numWrites, ioStall, bytesRead, bytesWritten))
	}

	return result.String(), nil
}

// analyzeGeneralDMV análise genérica de DMV
func (s *SQLServerAnalyzer) analyzeGeneralDMV(dmvType string, beginTime, endTime time.Time) (string, error) {
	return fmt.Sprintf("Análise de DMV tipo '%s' ainda não implementada", dmvType), nil
}

// AnalyzeInstance analisa a instância SQL Server geral
func (s *SQLServerAnalyzer) AnalyzeInstance() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise da Instância SQL Server\n\n")

	// Informações gerais da instância
	result.WriteString("## Informações da Instância\n\n")
	instanceQuery := `
		SELECT 
			@@SERVERNAME AS server_name,
			@@VERSION AS version,
			@@SERVICENAME AS service_name,
			@@MAX_CONNECTIONS AS max_connections,
			DB_NAME() AS current_database,
			GETDATE() AS current_time
	`

	var serverName, version, serviceName, currentDB string
	var maxConnections int
	var currentTime time.Time

	err := s.db.QueryRow(instanceQuery).Scan(&serverName, &version, &serviceName, &maxConnections, &currentDB, &currentTime)
	if err == nil {
		result.WriteString(fmt.Sprintf("- **Servidor:** %s\n", serverName))
		result.WriteString(fmt.Sprintf("- **Versão:** %s\n", version))
		result.WriteString(fmt.Sprintf("- **Serviço:** %s\n", serviceName))
		result.WriteString(fmt.Sprintf("- **Max Conexões:** %d\n", maxConnections))
		result.WriteString(fmt.Sprintf("- **Database Atual:** %s\n", currentDB))
		result.WriteString(fmt.Sprintf("- **Hora Atual:** %s\n\n", currentTime.Format("2006-01-02 15:04:05")))
	}

	// Estatísticas de memória
	result.WriteString("## Memória da Instância\n\n")
	memoryQuery := `
		SELECT 
			physical_memory_kb / 1024 AS physical_memory_mb,
			committed_kb / 1024 AS committed_mb,
			committed_target_kb / 1024 AS committed_target_mb,
			available_physical_memory_kb / 1024 AS available_mb
		FROM sys.dm_os_sys_info
	`

	var physicalMB, committedMB, targetMB, availableMB int64
	err = s.db.QueryRow(memoryQuery).Scan(&physicalMB, &committedMB, &targetMB, &availableMB)
	if err == nil {
		result.WriteString(fmt.Sprintf("- **Memória Física:** %d MB\n", physicalMB))
		result.WriteString(fmt.Sprintf("- **Memória Comprometida:** %d MB\n", committedMB))
		result.WriteString(fmt.Sprintf("- **Target de Memória:** %d MB\n", targetMB))
		result.WriteString(fmt.Sprintf("- **Memória Disponível:** %d MB\n\n", availableMB))
	}

	// Estatísticas de CPU
	result.WriteString("## Estatísticas de CPU\n\n")
	cpuQuery := `
		SELECT 
			cpu_count,
			hyperthread_ratio,
			physical_memory_kb / 1024 AS memory_mb,
			sqlserver_start_time
		FROM sys.dm_os_sys_info
	`

	var cpuCount, hyperthreadRatio int
	var memoryMB int64
	var startTime time.Time
	err = s.db.QueryRow(cpuQuery).Scan(&cpuCount, &hyperthreadRatio, &memoryMB, &startTime)
	if err == nil {
		result.WriteString(fmt.Sprintf("- **CPUs:** %d\n", cpuCount))
		result.WriteString(fmt.Sprintf("- **Hyperthread Ratio:** %d\n", hyperthreadRatio))
		result.WriteString(fmt.Sprintf("- **Memória:** %d MB\n", memoryMB))
		result.WriteString(fmt.Sprintf("- **Início do Servidor:** %s\n\n", startTime.Format("2006-01-02 15:04:05")))
	}

	// Conexões ativas
	result.WriteString("## Conexões Ativas\n\n")
	connQuery := `
		SELECT 
			COUNT(*) AS total_connections,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) AS running,
			SUM(CASE WHEN status = 'sleeping' THEN 1 ELSE 0 END) AS sleeping,
			SUM(CASE WHEN status = 'runnable' THEN 1 ELSE 0 END) AS runnable
		FROM sys.dm_exec_sessions
		WHERE is_user_process = 1
	`

	var total, running, sleeping, runnable int
	err = s.db.QueryRow(connQuery).Scan(&total, &running, &sleeping, &runnable)
	if err == nil {
		result.WriteString(fmt.Sprintf("- **Total de Conexões:** %d\n", total))
		result.WriteString(fmt.Sprintf("- **Running:** %d\n", running))
		result.WriteString(fmt.Sprintf("- **Sleeping:** %d\n", sleeping))
		result.WriteString(fmt.Sprintf("- **Runnable:** %d\n\n", runnable))
	}

	return result.String(), nil
}

// AnalyzeDatabases analisa todos os databases na instância
func (s *SQLServerAnalyzer) AnalyzeDatabases() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Databases SQL Server\n\n")

	// Listar todos os databases
	query := `
		SELECT 
			database_id,
			name,
			state_desc,
			recovery_model_desc,
			compatibility_level,
			collation_name,
			create_date,
			user_access_desc
		FROM sys.databases
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar databases: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Databases Disponíveis\n\n")
	result.WriteString("| ID | Nome | Estado | Recovery Model | Compatibility | User Access |\n")
	result.WriteString("|----|------|--------|----------------|---------------|-------------|\n")

	for rows.Next() {
		var dbID, compatLevel int
		var name, state, recoveryModel, collation, userAccess string
		var createDate time.Time

		err := rows.Scan(&dbID, &name, &state, &recoveryModel, &compatLevel, &collation, &createDate, &userAccess)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %d | %s |\n",
			dbID, name, state, recoveryModel, compatLevel, userAccess))
	}

	// Tamanho dos databases
	result.WriteString("\n## Tamanho dos Databases\n\n")
	sizeQuery := `
		SELECT 
			d.name,
			SUM(mf.size) * 8 / 1024 AS size_mb,
			SUM(mf.size) * 8 / 1024 / 1024 AS size_gb
		FROM sys.master_files mf
		INNER JOIN sys.databases d ON mf.database_id = d.database_id
		GROUP BY d.name
		ORDER BY size_mb DESC
	`

	sizeRows, err := s.db.Query(sizeQuery)
	if err == nil {
		defer sizeRows.Close()
		result.WriteString("| Database | Tamanho (MB) | Tamanho (GB) |\n")
		result.WriteString("|----------|--------------|--------------|\n")

		for sizeRows.Next() {
			var name string
			var sizeMB, sizeGB float64

			if err := sizeRows.Scan(&name, &sizeMB, &sizeGB); err == nil {
				result.WriteString(fmt.Sprintf("| %s | %.2f | %.2f |\n", name, sizeMB, sizeGB))
			}
		}
	}

	// Estatísticas de I/O por database
	result.WriteString("\n## Estatísticas de I/O por Database\n\n")
	ioQuery := `
		SELECT 
			DB_NAME(database_id) AS database_name,
			SUM(num_of_reads) AS total_reads,
			SUM(num_of_writes) AS total_writes,
			SUM(io_stall_read_ms) AS read_stall_ms,
			SUM(io_stall_write_ms) AS write_stall_ms
		FROM sys.dm_io_virtual_file_stats(NULL, NULL)
		GROUP BY database_id
		ORDER BY total_reads DESC
	`

	ioRows, err := s.db.Query(ioQuery)
	if err == nil {
		defer ioRows.Close()
		result.WriteString("| Database | Reads | Writes | Read Stall (ms) | Write Stall (ms) |\n")
		result.WriteString("|----------|-------|--------|-----------------|------------------|\n")

		for ioRows.Next() {
			var dbName sql.NullString
			var reads, writes, readStall, writeStall int64

			if err := ioRows.Scan(&dbName, &reads, &writes, &readStall, &writeStall); err == nil {
				name := "N/A"
				if dbName.Valid {
					name = dbName.String
				}
				result.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d |\n",
					name, reads, writes, readStall, writeStall))
			}
		}
	}

	return result.String(), nil
}

// AnalyzeSpecificDatabase analisa um database específico
func (s *SQLServerAnalyzer) AnalyzeSpecificDatabase(databaseName string) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Análise Detalhada do Database: %s\n\n", databaseName))

	// Informações do database
	infoQuery := fmt.Sprintf(`
		USE [%s]
		SELECT 
			DB_NAME() AS database_name,
			DB_ID() AS database_id,
			DATABASEPROPERTYEX(DB_NAME(), 'Recovery') AS recovery_model,
			DATABASEPROPERTYEX(DB_NAME(), 'Status') AS status,
			DATABASEPROPERTYEX(DB_NAME(), 'Updateability') AS updateability
	`, databaseName)

	var dbName string
	var dbID int
	var recoveryModel, status, updateability sql.NullString

	err := s.db.QueryRow(infoQuery).Scan(&dbName, &dbID, &recoveryModel, &status, &updateability)
	if err != nil {
		return "", fmt.Errorf("erro ao obter informações do database: %w", err)
	}

	result.WriteString(fmt.Sprintf("**Database ID:** %d\n", dbID))
	result.WriteString(fmt.Sprintf("**Recovery Model:** %s\n", getString(recoveryModel)))
	result.WriteString(fmt.Sprintf("**Status:** %s\n", getString(status)))
	result.WriteString(fmt.Sprintf("**Updateability:** %s\n\n", getString(updateability)))

	// Tabelas no database
	result.WriteString("## Tabelas no Database\n\n")
	tableQuery := fmt.Sprintf(`
		USE [%s]
		SELECT TOP 20
			SCHEMA_NAME(schema_id) AS schema_name,
			name AS table_name,
			OBJECT_ID(name) AS object_id,
			create_date,
			modify_date
		FROM sys.tables
		ORDER BY name
	`, databaseName)

	tableRows, err := s.db.Query(tableQuery)
	if err == nil {
		defer tableRows.Close()
		result.WriteString("| Schema | Tabela | Object ID | Criada | Modificada |\n")
		result.WriteString("|--------|--------|-----------|--------|------------|\n")

		for tableRows.Next() {
			var schema, table string
			var objectID int
			var createDate, modifyDate time.Time

			if err := tableRows.Scan(&schema, &table, &objectID, &createDate, &modifyDate); err == nil {
				result.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s |\n",
					schema, table, objectID, createDate.Format("2006-01-02"), modifyDate.Format("2006-01-02")))
			}
		}
	}

	// Estatísticas de espaço
	result.WriteString("\n## Uso de Espaço\n\n")
	spaceQuery := fmt.Sprintf(`
		USE [%s]
		SELECT 
			SUM(size) * 8 / 1024 AS total_size_mb,
			SUM(CASE WHEN type_desc = 'ROWS' THEN size ELSE 0 END) * 8 / 1024 AS data_size_mb,
			SUM(CASE WHEN type_desc = 'LOG' THEN size ELSE 0 END) * 8 / 1024 AS log_size_mb
		FROM sys.database_files
	`, databaseName)

	var totalMB, dataMB, logMB float64
	err = s.db.QueryRow(spaceQuery).Scan(&totalMB, &dataMB, &logMB)
	if err == nil {
		result.WriteString(fmt.Sprintf("- **Tamanho Total:** %.2f MB\n", totalMB))
		result.WriteString(fmt.Sprintf("- **Tamanho de Dados:** %.2f MB\n", dataMB))
		result.WriteString(fmt.Sprintf("- **Tamanho de Log:** %.2f MB\n\n", logMB))
	}

	return result.String(), nil
}

// GetExecutionPlan obtém plano de execução
func (s *SQLServerAnalyzer) GetExecutionPlan(sqlHandle string) (string, error) {
	var result strings.Builder

	result.WriteString("# Plano de Execução SQL Server\n\n")

	query := `
		SELECT 
			p.query_plan
		FROM sys.dm_exec_query_plan(CONVERT(varbinary(64), :sql_handle, 1)) p
	`

	var queryPlan string
	err := s.db.QueryRow(query, sqlHandle).Scan(&queryPlan)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar plano de execução: %w", err)
	}

	result.WriteString("## XML do Plano de Execução\n\n")
	result.WriteString("```xml\n")
	result.WriteString(queryPlan)
	result.WriteString("\n```\n")

	return result.String(), nil
}

