package postgresql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// PostgreSQLAnalyzer realiza análises específicas do PostgreSQL
type PostgreSQLAnalyzer struct {
	db *sql.DB
}

// NewPostgreSQLAnalyzer cria um novo analisador PostgreSQL
func NewPostgreSQLAnalyzer(db *sql.DB) *PostgreSQLAnalyzer {
	return &PostgreSQLAnalyzer{db: db}
}

// AnalyzePgStat analisa pg_stat views (equivalente ao AWR do Oracle)
func (p *PostgreSQLAnalyzer) AnalyzePgStat(statType string, beginTime, endTime time.Time) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Análise pg_stat - %s\n\n", statType))

	switch statType {
	case "database":
		return p.analyzeDatabaseStats()
	case "table":
		return p.analyzeTableStats()
	case "index":
		return p.analyzeIndexStats()
	case "query":
		return p.analyzeQueryStats()
	case "activity":
		return p.analyzeActivity()
	default:
		return p.analyzeGeneralStats(statType)
	}
}

// analyzeDatabaseStats analisa estatísticas de banco de dados
func (p *PostgreSQLAnalyzer) analyzeDatabaseStats() (string, error) {
	var result strings.Builder

	result.WriteString("# Estatísticas de Banco de Dados\n\n")

	query := `
		SELECT 
			datname,
			numbackends,
			xact_commit,
			xact_rollback,
			blks_read,
			blks_hit,
			tup_returned,
			tup_fetched,
			tup_inserted,
			tup_updated,
			tup_deleted
		FROM pg_stat_database
		WHERE datname NOT IN ('template0', 'template1', 'postgres')
		ORDER BY blks_hit DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar stats de banco: %w", err)
	}
	defer rows.Close()

	result.WriteString("| Database | Backends | Commits | Rollbacks | Blks Read | Blks Hit |\n")
	result.WriteString("|----------|----------|---------|-----------|-----------|----------|\n")

	for rows.Next() {
		var datname string
		var numbackends, xactCommit, xactRollback, blksRead, blksHit int64
		var tupReturned, tupFetched, tupInserted, tupUpdated, tupDeleted int64

		err := rows.Scan(&datname, &numbackends, &xactCommit, &xactRollback,
			&blksRead, &blksHit, &tupReturned, &tupFetched,
			&tupInserted, &tupUpdated, &tupDeleted)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d | %d |\n",
			datname, numbackends, xactCommit, xactRollback, blksRead, blksHit))
	}

	return result.String(), nil
}

// analyzeTableStats analisa estatísticas de tabelas
func (p *PostgreSQLAnalyzer) analyzeTableStats() (string, error) {
	var result strings.Builder

	result.WriteString("# Estatísticas de Tabelas\n\n")

	query := `
		SELECT 
			schemaname,
			tablename,
			seq_scan,
			seq_tup_read,
			idx_scan,
			idx_tup_fetch,
			n_tup_ins,
			n_tup_upd,
			n_tup_del,
			n_live_tup,
			n_dead_tup,
			last_vacuum,
			last_autovacuum,
			last_analyze,
			last_autoanalyze
		FROM pg_stat_user_tables
		ORDER BY seq_scan + idx_scan DESC
		LIMIT 20
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar stats de tabelas: %w", err)
	}
	defer rows.Close()

	result.WriteString("| Schema | Tabela | Seq Scan | Index Scan | Inserts | Updates | Deletes | Live Tuples |\n")
	result.WriteString("|--------|--------|----------|------------|---------|---------|---------|-------------|\n")

	for rows.Next() {
		var schemaname, tablename string
		var seqScan, seqTupRead, idxScan, idxTupFetch int64
		var nTupIns, nTupUpd, nTupDel, nLiveTup, nDeadTup int64
		var lastVacuum, lastAutovacuum, lastAnalyze, lastAutoanalyze sql.NullTime

		err := rows.Scan(&schemaname, &tablename, &seqScan, &seqTupRead, &idxScan, &idxTupFetch,
			&nTupIns, &nTupUpd, &nTupDel, &nLiveTup, &nDeadTup,
			&lastVacuum, &lastAutovacuum, &lastAnalyze, &lastAutoanalyze)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %d | %d | %d | %d | %d | %d |\n",
			schemaname, tablename, seqScan, idxScan, nTupIns, nTupUpd, nTupDel, nLiveTup))
	}

	return result.String(), nil
}

// analyzeIndexStats analisa estatísticas de índices
func (p *PostgreSQLAnalyzer) analyzeIndexStats() (string, error) {
	var result strings.Builder

	result.WriteString("# Estatísticas de Índices\n\n")

	query := `
		SELECT 
			schemaname,
			tablename,
			indexname,
			idx_scan,
			idx_tup_read,
			idx_tup_fetch
		FROM pg_stat_user_indexes
		ORDER BY idx_scan ASC
		LIMIT 20
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar stats de índices: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Índices Não Utilizados\n\n")
	result.WriteString("| Schema | Tabela | Índice | Scans | Tuples Read | Tuples Fetched |\n")
	result.WriteString("|--------|--------|--------|-------|-------------|----------------|\n")

	for rows.Next() {
		var schemaname, tablename, indexname string
		var idxScan, idxTupRead, idxTupFetch int64

		err := rows.Scan(&schemaname, &tablename, &indexname, &idxScan, &idxTupRead, &idxTupFetch)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d | %d |\n",
			schemaname, tablename, indexname, idxScan, idxTupRead, idxTupFetch))
	}

	return result.String(), nil
}

// analyzeQueryStats analisa estatísticas de queries
func (p *PostgreSQLAnalyzer) analyzeQueryStats() (string, error) {
	var result strings.Builder

	result.WriteString("# Estatísticas de Queries\n\n")

	query := `
		SELECT 
			query,
			calls,
			total_exec_time,
			mean_exec_time,
			max_exec_time,
			min_exec_time,
			stddev_exec_time
		FROM pg_stat_statements
		ORDER BY total_exec_time DESC
		LIMIT 20
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar pg_stat_statements (pode não estar habilitado): %w", err)
	}
	defer rows.Close()

	result.WriteString("## Top Queries por Tempo Total\n\n")
	result.WriteString("| Query | Chamadas | Tempo Total (ms) | Tempo Médio (ms) | Tempo Máx (ms) |\n")
	result.WriteString("|-------|----------|------------------|------------------|----------------|\n")

	for rows.Next() {
		var queryText string
		var calls int64
		var totalExecTime, meanExecTime, maxExecTime, minExecTime, stddevExecTime float64

		err := rows.Scan(&queryText, &calls, &totalExecTime, &meanExecTime,
			&maxExecTime, &minExecTime, &stddevExecTime)
		if err != nil {
			continue
		}

		// Truncar query se muito longa
		if len(queryText) > 100 {
			queryText = queryText[:100] + "..."
		}

		result.WriteString(fmt.Sprintf("| %s | %d | %.2f | %.2f | %.2f |\n",
			queryText, calls, totalExecTime, meanExecTime, maxExecTime))
	}

	return result.String(), nil
}

// analyzeActivity analisa atividade atual
func (p *PostgreSQLAnalyzer) analyzeActivity() (string, error) {
	var result strings.Builder

	result.WriteString("# Atividade Atual\n\n")

	query := `
		SELECT 
			pid,
			usename,
			application_name,
			client_addr,
			state,
			query_start,
			state_change,
			wait_event_type,
			wait_event,
			query
		FROM pg_stat_activity
		WHERE pid != pg_backend_pid()
		ORDER BY query_start
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar atividade: %w", err)
	}
	defer rows.Close()

	result.WriteString("| PID | Usuário | Aplicação | Estado | Wait Event | Query |\n")
	result.WriteString("|-----|---------|-----------|--------|------------|-------|\n")

	for rows.Next() {
		var pid int
		var usename, applicationName sql.NullString
		var clientAddr sql.NullString
		var state string
		var queryStart, stateChange sql.NullTime
		var waitEventType, waitEvent sql.NullString
		var query sql.NullString

		err := rows.Scan(&pid, &usename, &applicationName, &clientAddr, &state,
			&queryStart, &stateChange, &waitEventType, &waitEvent, &query)
		if err != nil {
			continue
		}

		queryText := "N/A"
		if query.Valid && len(query.String) > 50 {
			queryText = query.String[:50] + "..."
		} else if query.Valid {
			queryText = query.String
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
			pid, getString(usename), getString(applicationName), state,
			getString(waitEvent), queryText))
	}

	return result.String(), nil
}

// analyzeGeneralStats análise genérica
func (p *PostgreSQLAnalyzer) analyzeGeneralStats(statType string) (string, error) {
	return fmt.Sprintf("Análise de pg_stat tipo '%s' ainda não implementada", statType), nil
}

// GetExecutionPlan obtém plano de execução
func (p *PostgreSQLAnalyzer) GetExecutionPlan(query string) (string, error) {
	var result strings.Builder

	result.WriteString("# Plano de Execução PostgreSQL\n\n")

	explainQuery := "EXPLAIN (FORMAT JSON, ANALYZE, BUFFERS) " + query

	var planJSON string
	err := p.db.QueryRow(explainQuery).Scan(&planJSON)
	if err != nil {
		return "", fmt.Errorf("erro ao obter plano de execução: %w", err)
	}

	result.WriteString("## JSON do Plano de Execução\n\n")
	result.WriteString("```json\n")
	result.WriteString(planJSON)
	result.WriteString("\n```\n")

	return result.String(), nil
}

func getString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "N/A"
}

// AnalyzeReplication analisa replicação PostgreSQL
func (p *PostgreSQLAnalyzer) AnalyzeReplication() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Replicação PostgreSQL\n\n")

	// Verificar se é um servidor de streaming replication
	query := `
		SELECT 
			client_addr,
			state,
			sent_lsn,
			write_lsn,
			flush_lsn,
			replay_lsn,
			sync_priority,
			sync_state,
			pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn) as sent_lag,
			pg_wal_lsn_diff(pg_current_wal_lsn(), write_lsn) as write_lag,
			pg_wal_lsn_diff(pg_current_wal_lsn(), flush_lsn) as flush_lag,
			pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) as replay_lag
		FROM pg_stat_replication
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar replicação: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Replicas de Streaming\n\n")
	result.WriteString("| Cliente | Estado | Sync State | Sent Lag | Write Lag | Flush Lag | Replay Lag |\n")
	result.WriteString("|---------|--------|------------|----------|-----------|-----------|------------|\n")

	for rows.Next() {
		var clientAddr, state, sentLSN, writeLSN, flushLSN, replayLSN sql.NullString
		var syncPriority, syncState sql.NullString
		var sentLag, writeLag, flushLag, replayLag sql.NullInt64

		err := rows.Scan(&clientAddr, &state, &sentLSN, &writeLSN, &flushLSN, &replayLSN,
			&syncPriority, &syncState, &sentLag, &writeLag, &flushLag, &replayLag)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d | %d | %d |\n",
			getString(clientAddr), getString(state), getString(syncState),
			getInt64(sentLag), getInt64(writeLag), getInt64(flushLag), getInt64(replayLag)))
	}

	// Verificar slots de replicação
	result.WriteString("\n## Slots de Replicação\n\n")
	slotQuery := `
		SELECT 
			slot_name,
			plugin,
			slot_type,
			datoid,
			database,
			temporary,
			active,
			active_pid,
			xmin,
			catalog_xmin,
			restart_lsn,
			confirmed_flush_lsn
		FROM pg_replication_slots
	`

	slotRows, err := p.db.Query(slotQuery)
	if err == nil {
		defer slotRows.Close()

		result.WriteString("| Slot Name | Plugin | Type | Database | Active | Restart LSN |\n")
		result.WriteString("|-----------|--------|------|----------|--------|-------------|\n")

		for slotRows.Next() {
			var slotName, plugin, slotType, database sql.NullString
			var datoid sql.NullInt64
			var temporary, active bool
			var activePid sql.NullInt64
			var xmin, catalogXmin sql.NullInt64
			var restartLSN, confirmedFlushLSN sql.NullString

			err := slotRows.Scan(&slotName, &plugin, &slotType, &datoid, &database,
				&temporary, &active, &activePid, &xmin, &catalogXmin,
				&restartLSN, &confirmedFlushLSN)
			if err != nil {
				continue
			}

			result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %v | %s |\n",
				getString(slotName), getString(plugin), getString(slotType),
				getString(database), active, getString(restartLSN)))
		}
	}

	return result.String(), nil
}

// AnalyzeLocks analisa locks e bloqueios PostgreSQL
func (p *PostgreSQLAnalyzer) AnalyzeLocks() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Locks e Bloqueios PostgreSQL\n\n")

	// Locks ativos
	query := `
		SELECT 
			l.locktype,
			l.database,
			l.relation::regclass,
			l.page,
			l.tuple,
			l.virtualxid,
			l.transactionid,
			l.classid,
			l.objid,
			l.objsubid,
			l.virtualtransaction,
			l.pid,
			l.mode,
			l.granted,
			a.usename,
			a.query,
			a.query_start,
			age(now(), a.query_start) AS age
		FROM pg_locks l
		LEFT JOIN pg_stat_activity a ON l.pid = a.pid
		WHERE l.database = (SELECT oid FROM pg_database WHERE datname = current_database())
		ORDER BY a.query_start
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar locks: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Locks Ativos\n\n")
	result.WriteString("| PID | Usuário | Tipo | Relação | Modo | Concedido | Query | Idade |\n")
	result.WriteString("|-----|---------|------|---------|------|-----------|-------|-------|\n")

	for rows.Next() {
		var lockType, relation, virtualXID, virtualTransaction, mode sql.NullString
		var database, page, tuple, transactionID, classID, objID, objSubID, pid sql.NullInt64
		var granted bool
		var usename, query sql.NullString
		var queryStart sql.NullTime
		var age sql.NullString

		err := rows.Scan(&lockType, &database, &relation, &page, &tuple, &virtualXID,
			&transactionID, &classID, &objID, &objSubID, &virtualTransaction, &pid,
			&mode, &granted, &usename, &query, &queryStart, &age)
		if err != nil {
			continue
		}

		queryText := "N/A"
		if query.Valid && len(query.String) > 50 {
			queryText = query.String[:50] + "..."
		} else if query.Valid {
			queryText = query.String
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %v | %s | %s |\n",
			getInt64(pid), getString(usename), getString(lockType), getString(relation),
			getString(mode), granted, queryText, getString(age)))
	}

	// Bloqueios (deadlocks)
	result.WriteString("\n## Bloqueios Detectados\n\n")
	blockingQuery := `
		SELECT 
			blocked_locks.pid AS blocked_pid,
			blocked_activity.usename AS blocked_user,
			blocking_locks.pid AS blocking_pid,
			blocking_activity.usename AS blocking_user,
			blocked_activity.query AS blocked_statement,
			blocking_activity.query AS blocking_statement
		FROM pg_catalog.pg_locks blocked_locks
		JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
		JOIN pg_catalog.pg_locks blocking_locks 
			ON blocking_locks.locktype = blocked_locks.locktype
			AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
			AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
			AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
			AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
			AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
			AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
			AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
			AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
			AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
			AND blocking_locks.pid != blocked_locks.pid
		JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
		WHERE NOT blocked_locks.granted
	`

	blockingRows, err := p.db.Query(blockingQuery)
	if err == nil {
		defer blockingRows.Close()

		result.WriteString("| Bloqueado (PID) | Bloqueador (PID) | Query Bloqueada | Query Bloqueadora |\n")
		result.WriteString("|-----------------|------------------|-----------------|-------------------|\n")

		for blockingRows.Next() {
			var blockedPID, blockingPID int64
			var blockedUser, blockingUser, blockedStatement, blockingStatement sql.NullString

			err := blockingRows.Scan(&blockedPID, &blockedUser, &blockingPID, &blockingUser,
				&blockedStatement, &blockingStatement)
			if err != nil {
				continue
			}

			blockedStmt := "N/A"
			if blockedStatement.Valid && len(blockedStatement.String) > 50 {
				blockedStmt = blockedStatement.String[:50] + "..."
			} else if blockedStatement.Valid {
				blockedStmt = blockedStatement.String
			}

			blockingStmt := "N/A"
			if blockingStatement.Valid && len(blockingStatement.String) > 50 {
				blockingStmt = blockingStatement.String[:50] + "..."
			} else if blockingStatement.Valid {
				blockingStmt = blockingStatement.String
			}

			result.WriteString(fmt.Sprintf("| %d (%s) | %d (%s) | %s | %s |\n",
				blockedPID, getString(blockedUser), blockingPID, getString(blockingUser),
				blockedStmt, blockingStmt))
		}
	}

	return result.String(), nil
}

// AnalyzeFragmentation analisa fragmentação PostgreSQL
func (p *PostgreSQLAnalyzer) AnalyzeFragmentation() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Fragmentação PostgreSQL\n\n")

	query := `
		SELECT 
			schemaname,
			tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
			pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
			pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) AS indexes_size,
			pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename) - COALESCE(pg_indexes_size(schemaname||'.'||tablename), 0) AS bloat_size,
			n_dead_tup,
			n_live_tup,
			CASE 
				WHEN n_live_tup > 0 
				THEN ROUND(100.0 * n_dead_tup / (n_live_tup + n_dead_tup), 2)
				ELSE 0
			END AS dead_tuple_percent,
			last_vacuum,
			last_autovacuum,
			last_analyze,
			last_autoanalyze
		FROM pg_stat_user_tables
		WHERE n_dead_tup > 0
		ORDER BY bloat_size DESC
		LIMIT 20
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar fragmentação: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Tabelas com Fragmentação\n\n")
	result.WriteString("| Schema | Tabela | Tamanho Total | Tamanho Tabela | Dead Tuples | % Dead | Último Vacuum |\n")
	result.WriteString("|--------|--------|---------------|----------------|-------------|--------|---------------|\n")

	for rows.Next() {
		var schemaname, tablename, totalSize, tableSize, indexesSize sql.NullString
		var bloatSize, nDeadTup, nLiveTup sql.NullInt64
		var deadTuplePercent sql.NullFloat64
		var lastVacuum, lastAutovacuum, lastAnalyze, lastAutoanalyze sql.NullTime

		err := rows.Scan(&schemaname, &tablename, &totalSize, &tableSize, &indexesSize,
			&bloatSize, &nDeadTup, &nLiveTup, &deadTuplePercent,
			&lastVacuum, &lastAutovacuum, &lastAnalyze, &lastAutoanalyze)
		if err != nil {
			continue
		}

		lastVacuumStr := "Nunca"
		if lastVacuum.Valid {
			lastVacuumStr = lastVacuum.Time.Format("2006-01-02 15:04:05")
		} else if lastAutovacuum.Valid {
			lastVacuumStr = "Auto: " + lastAutovacuum.Time.Format("2006-01-02 15:04:05")
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d | %.2f%% | %s |\n",
			getString(schemaname), getString(tablename), getString(totalSize),
			getString(tableSize), getInt64(nDeadTup), getFloat64(deadTuplePercent), lastVacuumStr))
	}

	return result.String(), nil
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

