package oracle

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// AWRReportRequest representa uma requisição de relatório AWR
type AWRReportRequest struct {
	DBID          int64
	InstanceID    int
	BeginSnapshot int
	EndSnapshot   int
	BeginTime     time.Time
	EndTime       time.Time
	ReportType    string // "html" ou "text"
}

// ASHRequest representa uma requisição de análise ASH
type ASHRequest struct {
	SQLID      string
	Serial     int
	SID        int
	BeginTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
}

// OracleAnalyzer realiza análises específicas do Oracle
type OracleAnalyzer struct {
	db *sql.DB
}

// NewOracleAnalyzer cria um novo analisador Oracle
func NewOracleAnalyzer(db *sql.DB) *OracleAnalyzer {
	return &OracleAnalyzer{db: db}
}

// GetSnapshots obtém snapshots disponíveis
func (o *OracleAnalyzer) GetSnapshots(beginTime, endTime time.Time) ([]Snapshot, error) {
	query := `
		SELECT snap_id, instance_number, begin_interval_time, end_interval_time
		FROM dba_hist_snapshot
		WHERE begin_interval_time >= :1 AND end_interval_time <= :2
		ORDER BY snap_id
	`

	rows, err := o.db.Query(query, beginTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []Snapshot
	for rows.Next() {
		var snap Snapshot
		err := rows.Scan(&snap.ID, &snap.InstanceNumber, &snap.BeginTime, &snap.EndTime)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snap)
	}

	return snapshots, nil
}

// GenerateAWRReport gera um relatório AWR
func (o *OracleAnalyzer) GenerateAWRReport(req AWRReportRequest) (string, error) {
	var report strings.Builder

	// Determinar snapshots se não fornecidos
	if req.BeginSnapshot == 0 || req.EndSnapshot == 0 {
		snapshots, err := o.GetSnapshots(req.BeginTime, req.EndTime)
		if err != nil {
			return "", err
		}
		if len(snapshots) < 2 {
			return "", fmt.Errorf("é necessário pelo menos 2 snapshots no período")
		}
		req.BeginSnapshot = snapshots[0].ID
		req.EndSnapshot = snapshots[len(snapshots)-1].ID
	}

	// Gerar relatório AWR usando DBMS_WORKLOAD_REPOSITORY
	query := `
		SELECT output
		FROM TABLE(DBMS_WORKLOAD_REPOSITORY.AWR_REPORT_HTML(
			:dbid, :inst_num, :bid, :eid
		))
	`

	if req.ReportType == "text" {
		query = `
			SELECT output
			FROM TABLE(DBMS_WORKLOAD_REPOSITORY.AWR_REPORT_TEXT(
				:dbid, :inst_num, :bid, :eid
			))
		`
	}

	var output string
	err := o.db.QueryRow(query, req.DBID, req.InstanceID, req.BeginSnapshot, req.EndSnapshot).Scan(&output)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar relatório AWR: %w", err)
	}

	report.WriteString(output)
	return report.String(), nil
}

// AnalyzeAWRReportFile analisa um arquivo de relatório AWR existente
func (o *OracleAnalyzer) AnalyzeAWRReportFile(filePath string) (string, error) {
	// Ler arquivo e analisar conteúdo
	// Esta função será implementada para ler e parsear arquivos AWR
	return "", fmt.Errorf("análise de arquivo AWR ainda não implementada")
}

// AnalyzeASH analisa Active Session History
func (o *OracleAnalyzer) AnalyzeASH(req ASHRequest) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise ASH (Active Session History)\n\n")

	// Query ASH básica
	query := `
		SELECT 
			sample_time,
			session_id,
			session_serial#,
			sql_id,
			event,
			wait_class,
			time_waited,
			session_state
		FROM v$active_session_history
		WHERE 1=1
	`

	args := []interface{}{}

	if req.SQLID != "" {
		query += " AND sql_id = :sqlid"
		args = append(args, req.SQLID)
	}

	if req.SID > 0 {
		query += " AND session_id = :sid"
		args = append(args, req.SID)
	}

	if req.Serial > 0 {
		query += " AND session_serial# = :serial"
		args = append(args, req.Serial)
	}

	if !req.BeginTime.IsZero() {
		query += " AND sample_time >= :begin_time"
		args = append(args, req.BeginTime)
	}

	if !req.EndTime.IsZero() {
		query += " AND sample_time <= :end_time"
		args = append(args, req.EndTime)
	}

	query += " ORDER BY sample_time DESC"

	rows, err := o.db.Query(query, args...)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar ASH: %w", err)
	}
	defer rows.Close()

	result.WriteString("## Sessões Ativas\n\n")
	count := 0
	for rows.Next() {
		var sampleTime time.Time
		var sessionID, sessionSerial, timeWaited int
		var sqlID, event, waitClass, sessionState sql.NullString

		err := rows.Scan(&sampleTime, &sessionID, &sessionSerial, &sqlID, &event, &waitClass, &timeWaited, &sessionState)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("### Sessão %d (Serial: %d)\n", sessionID, sessionSerial))
		result.WriteString(fmt.Sprintf("- **Tempo:** %s\n", sampleTime.Format("2006-01-02 15:04:05")))
		if sqlID.Valid {
			result.WriteString(fmt.Sprintf("- **SQL ID:** %s\n", sqlID.String))
		}
		if event.Valid {
			result.WriteString(fmt.Sprintf("- **Evento:** %s\n", event.String))
		}
		if waitClass.Valid {
			result.WriteString(fmt.Sprintf("- **Classe de Wait:** %s\n", waitClass.String))
		}
		result.WriteString(fmt.Sprintf("- **Tempo de Wait:** %d ms\n", timeWaited))
		result.WriteString("\n")

		count++
		if count >= 100 { // Limitar resultados
			result.WriteString(fmt.Sprintf("\n... e mais resultados (limitado a 100)\n"))
			break
		}
	}

	return result.String(), nil
}

// GetExecutionPlan obtém plano de execução de um SQL ID
func (o *OracleAnalyzer) GetExecutionPlan(sqlID string) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Plano de Execução - SQL ID: %s\n\n", sqlID))

	// Obter plano de execução do AWR
	query := `
		SELECT 
			plan_hash_value,
			timestamp,
			operation,
			options,
			object_name,
			cost,
			cardinality,
			bytes
		FROM dba_hist_sql_plan
		WHERE sql_id = :sqlid
		ORDER BY timestamp DESC, id
		FETCH FIRST 1 ROW ONLY
	`

	rows, err := o.db.Query(query, sqlID)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar plano de execução: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var planHashValue int64
		var timestamp time.Time
		var operation, options, objectName sql.NullString
		var cost, cardinality, bytes sql.NullInt64

		err := rows.Scan(&planHashValue, &timestamp, &operation, &options, &objectName, &cost, &cardinality, &bytes)
		if err != nil {
			return "", err
		}

		result.WriteString(fmt.Sprintf("**Plan Hash Value:** %d\n", planHashValue))
		result.WriteString(fmt.Sprintf("**Timestamp:** %s\n\n", timestamp.Format("2006-01-02 15:04:05")))

		// Buscar detalhes completos do plano
		detailQuery := `
			SELECT 
				id,
				operation,
				options,
				object_name,
				cost,
				cardinality,
				bytes,
				time
			FROM dba_hist_sql_plan
			WHERE sql_id = :sqlid AND plan_hash_value = :plan_hash
			ORDER BY id
		`

		detailRows, err := o.db.Query(detailQuery, sqlID, planHashValue)
		if err == nil {
			defer detailRows.Close()
			result.WriteString("## Detalhes do Plano\n\n")
			result.WriteString("| ID | Operação | Objeto | Custo | Cardinalidade |\n")
			result.WriteString("|----|----------|--------|-------|---------------|\n")

			for detailRows.Next() {
				var id int
				var op, opts, obj sql.NullString
				var c, card, b sql.NullInt64

				err := detailRows.Scan(&id, &op, &opts, &obj, &c, &card, &b, nil)
				if err != nil {
					continue
				}

				operation := op.String
				if opts.Valid {
					operation += " " + opts.String
				}

				result.WriteString(fmt.Sprintf("| %d | %s | %s | %d | %d |\n",
					id, operation, getString(obj), getInt64(c), getInt64(card)))
			}
		}
	}

	return result.String(), nil
}

// GetSQLStatistics obtém estatísticas de um SQL ID
func (o *OracleAnalyzer) GetSQLStatistics(sqlID string) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Estatísticas SQL - SQL ID: %s\n\n", sqlID))

	query := `
		SELECT 
			elapsed_time,
			cpu_time,
			buffer_gets,
			disk_reads,
			direct_writes,
			executions,
			rows_processed,
			first_load_time,
			last_load_time
		FROM v$sqlstats
		WHERE sql_id = :sqlid
	`

	var elapsedTime, cpuTime, bufferGets, diskReads, directWrites, executions, rowsProcessed int64
	var firstLoadTime, lastLoadTime time.Time

	err := o.db.QueryRow(query, sqlID).Scan(
		&elapsedTime, &cpuTime, &bufferGets, &diskReads, &directWrites,
		&executions, &rowsProcessed, &firstLoadTime, &lastLoadTime,
	)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar estatísticas SQL: %w", err)
	}

	result.WriteString("## Métricas de Performance\n\n")
	result.WriteString(fmt.Sprintf("- **Tempo Total:** %d ms\n", elapsedTime/1000000))
	result.WriteString(fmt.Sprintf("- **Tempo CPU:** %d ms\n", cpuTime/1000000))
	result.WriteString(fmt.Sprintf("- **Buffer Gets:** %d\n", bufferGets))
	result.WriteString(fmt.Sprintf("- **Disk Reads:** %d\n", diskReads))
	result.WriteString(fmt.Sprintf("- **Execuções:** %d\n", executions))
	result.WriteString(fmt.Sprintf("- **Linhas Processadas:** %d\n", rowsProcessed))
	result.WriteString(fmt.Sprintf("- **Primeira Execução:** %s\n", firstLoadTime.Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("- **Última Execução:** %s\n", lastLoadTime.Format("2006-01-02 15:04:05")))

	return result.String(), nil
}

// Snapshot representa um snapshot AWR
type Snapshot struct {
	ID             int
	InstanceNumber int
	BeginTime      time.Time
	EndTime        time.Time
}

// Funções auxiliares
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

// AnalyzePDBs analisa Pluggable Databases (PDBs) no Oracle
func (o *OracleAnalyzer) AnalyzePDBs() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Pluggable Databases (PDBs)\n\n")

	// Listar PDBs
	query := `
		SELECT 
			pdb_id,
			pdb_name,
			status,
			creation_scn,
			con_id,
			guid
		FROM cdb_pdbs
		ORDER BY pdb_name
	`

	rows, err := o.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar PDBs: %w", err)
	}
	defer rows.Close()

	result.WriteString("## PDBs Disponíveis\n\n")
	result.WriteString("| PDB ID | Nome | Status | CON ID | GUID |\n")
	result.WriteString("|--------|------|--------|--------|------|\n")

	for rows.Next() {
		var pdbID, conID int
		var pdbName, status, guid string
		var creationSCN sql.NullInt64

		err := rows.Scan(&pdbID, &pdbName, &status, &creationSCN, &conID, &guid)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %d | %s |\n",
			pdbID, pdbName, status, conID, guid))
	}

	// Análise de uso de espaço por PDB
	result.WriteString("\n## Uso de Espaço por PDB\n\n")
	spaceQuery := `
		SELECT 
			pdb_name,
			SUM(bytes)/1024/1024/1024 AS size_gb,
			SUM(bytes)/1024/1024 AS size_mb
		FROM cdb_data_files
		GROUP BY pdb_name
		ORDER BY size_gb DESC
	`

	spaceRows, err := o.db.Query(spaceQuery)
	if err == nil {
		defer spaceRows.Close()
		result.WriteString("| PDB | Tamanho (GB) | Tamanho (MB) |\n")
		result.WriteString("|-----|---------------|--------------|\n")

		for spaceRows.Next() {
			var pdbName sql.NullString
			var sizeGB, sizeMB float64

			if err := spaceRows.Scan(&pdbName, &sizeGB, &sizeMB); err == nil {
				name := "N/A"
				if pdbName.Valid {
					name = pdbName.String
				}
				result.WriteString(fmt.Sprintf("| %s | %.2f | %.2f |\n", name, sizeGB, sizeMB))
			}
		}
	}

	// Análise de sessões por PDB
	result.WriteString("\n## Sessões por PDB\n\n")
	sessionQuery := `
		SELECT 
			con_id,
			COUNT(*) AS session_count,
			SUM(CASE WHEN status = 'ACTIVE' THEN 1 ELSE 0 END) AS active_sessions
		FROM v$session
		WHERE con_id > 0
		GROUP BY con_id
		ORDER BY session_count DESC
	`

	sessionRows, err := o.db.Query(sessionQuery)
	if err == nil {
		defer sessionRows.Close()
		result.WriteString("| CON ID | Total Sessões | Sessões Ativas |\n")
		result.WriteString("|--------|---------------|----------------|\n")

		for sessionRows.Next() {
			var conID, totalSessions, activeSessions int

			if err := sessionRows.Scan(&conID, &totalSessions, &activeSessions); err == nil {
				result.WriteString(fmt.Sprintf("| %d | %d | %d |\n", conID, totalSessions, activeSessions))
			}
		}
	}

	// Análise de performance por PDB
	result.WriteString("\n## Métricas de Performance por PDB\n\n")
	perfQuery := `
		SELECT 
			con_id,
			SUM(physical_reads) AS total_physical_reads,
			SUM(logical_reads) AS total_logical_reads,
			SUM(executions) AS total_executions
		FROM v$sqlstats
		WHERE con_id > 0
		GROUP BY con_id
		ORDER BY total_physical_reads DESC
	`

	perfRows, err := o.db.Query(perfQuery)
	if err == nil {
		defer perfRows.Close()
		result.WriteString("| CON ID | Physical Reads | Logical Reads | Execuções |\n")
		result.WriteString("|--------|----------------|---------------|------------|\n")

		for perfRows.Next() {
			var conID int
			var physicalReads, logicalReads, executions sql.NullInt64

			if err := perfRows.Scan(&conID, &physicalReads, &logicalReads, &executions); err == nil {
				result.WriteString(fmt.Sprintf("| %d | %d | %d | %d |\n",
					conID, getInt64(physicalReads), getInt64(logicalReads), getInt64(executions)))
			}
		}
	}

	return result.String(), nil
}

// AnalyzeSpecificPDB analisa um PDB específico
func (o *OracleAnalyzer) AnalyzeSpecificPDB(pdbName string) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Análise Detalhada do PDB: %s\n\n", pdbName))

	// Conectar ao PDB específico (requer alteração de contexto)
	// Por enquanto, analisamos usando queries com filtro por con_id
	query := `
		SELECT 
			con_id,
			pdb_name,
			status
		FROM cdb_pdbs
		WHERE pdb_name = :pdb_name
	`

	var conID int
	var status string
	err := o.db.QueryRow(query, pdbName).Scan(&conID, &pdbName, &status)
	if err != nil {
		return "", fmt.Errorf("PDB não encontrado: %w", err)
	}

	result.WriteString(fmt.Sprintf("**CON ID:** %d\n", conID))
	result.WriteString(fmt.Sprintf("**Status:** %s\n\n", status))

	// Análise de tabelas no PDB
	result.WriteString("## Tabelas no PDB\n\n")
	tableQuery := `
		SELECT 
			owner,
			table_name,
			num_rows,
			blocks,
			avg_row_len
		FROM cdb_tables
		WHERE con_id = :con_id
		AND owner NOT IN ('SYS', 'SYSTEM')
		ORDER BY num_rows DESC
		FETCH FIRST 20 ROWS ONLY
	`

	tableRows, err := o.db.Query(tableQuery, conID)
	if err == nil {
		defer tableRows.Close()
		result.WriteString("| Owner | Tabela | Linhas | Blocos | Avg Row Len |\n")
		result.WriteString("|-------|--------|--------|--------|-------------|\n")

		for tableRows.Next() {
			var owner, tableName string
			var numRows, blocks, avgRowLen sql.NullInt64

			if err := tableRows.Scan(&owner, &tableName, &numRows, &blocks, &avgRowLen); err == nil {
				result.WriteString(fmt.Sprintf("| %s | %s | %d | %d | %d |\n",
					owner, tableName, getInt64(numRows), getInt64(blocks), getInt64(avgRowLen)))
			}
		}
	}

	return result.String(), nil
}

