package oracle

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// RACAnalyzer realiza análises específicas do Oracle RAC
type RACAnalyzer struct {
	db *sql.DB
}

// NewRACAnalyzer cria um novo analisador RAC
func NewRACAnalyzer(db *sql.DB) *RACAnalyzer {
	return &RACAnalyzer{db: db}
}

// AnalyzeRACHealth analisa a saúde geral do cluster RAC
func (r *RACAnalyzer) AnalyzeRACHealth() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Saúde do Oracle RAC\n\n")

	// Verificar se é RAC
	isRAC, err := r.isRACDatabase()
	if err != nil {
		return "", fmt.Errorf("erro ao verificar se é RAC: %w", err)
	}

	if !isRAC {
		result.WriteString("⚠️ Este banco de dados não é um RAC (Real Application Cluster).\n")
		return result.String(), nil
	}

	result.WriteString("✅ Cluster RAC detectado\n\n")

	// Informações do cluster
	result.WriteString("## Informações do Cluster\n\n")
	clusterInfo, err := r.getClusterInfo()
	if err == nil {
		result.WriteString(clusterInfo)
	}

	// Status dos nós
	result.WriteString("\n## Status dos Nós do Cluster\n\n")
	nodesStatus, err := r.getNodesStatus()
	if err == nil {
		result.WriteString(nodesStatus)
	}

	// Serviços do cluster
	result.WriteString("\n## Serviços do Cluster\n\n")
	servicesStatus, err := r.getServicesStatus()
	if err == nil {
		result.WriteString(servicesStatus)
	}

	// Recursos do cluster
	result.WriteString("\n## Recursos do Cluster\n\n")
	resourcesStatus, err := r.getResourcesStatus()
	if err == nil {
		result.WriteString(resourcesStatus)
	}

	return result.String(), nil
}

// AnalyzeRACErrors analisa erros no cluster RAC
func (r *RACAnalyzer) AnalyzeRACErrors() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Erros no Oracle RAC\n\n")

	// Erros de instância
	result.WriteString("## Erros de Instância\n\n")
	instanceErrors, err := r.getInstanceErrors()
	if err == nil {
		result.WriteString(instanceErrors)
	}

	// Erros de clusterware
	result.WriteString("\n## Erros de Clusterware\n\n")
	clusterwareErrors, err := r.getClusterwareErrors()
	if err == nil {
		result.WriteString(clusterwareErrors)
	}

	// Erros de interconnect
	result.WriteString("\n## Erros de Interconnect\n\n")
	interconnectErrors, err := r.getInterconnectErrors()
	if err == nil {
		result.WriteString(interconnectErrors)
	}

	// Deadlocks entre instâncias
	result.WriteString("\n## Deadlocks entre Instâncias\n\n")
	deadlocks, err := r.getRACDeadlocks()
	if err == nil {
		result.WriteString(deadlocks)
	}

	return result.String(), nil
}

// AnalyzeListenerStatus analisa status do listener
func (r *RACAnalyzer) AnalyzeListenerStatus() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Status do Listener\n\n")

	// Status geral do listener
	result.WriteString("## Status do Listener\n\n")
	listenerStatus, err := r.getListenerStatus()
	if err == nil {
		result.WriteString(listenerStatus)
	}

	// Serviços registrados no listener
	result.WriteString("\n## Serviços Registrados\n\n")
	registeredServices, err := r.getRegisteredServices()
	if err == nil {
		result.WriteString(registeredServices)
	}

	// Conexões ativas no listener
	result.WriteString("\n## Conexões Ativas\n\n")
	activeConnections, err := r.getActiveConnections()
	if err == nil {
		result.WriteString(activeConnections)
	}

	// Erros do listener
	result.WriteString("\n## Erros do Listener\n\n")
	listenerErrors, err := r.getListenerErrors()
	if err == nil {
		result.WriteString(listenerErrors)
	}

	return result.String(), nil
}

// AnalyzeListenerLog analisa log do listener
func (r *RACAnalyzer) AnalyzeListenerLog(logPath string) (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Log do Listener\n\n")
	result.WriteString(fmt.Sprintf("**Arquivo:** %s\n\n", logPath))

	// Esta função seria implementada para ler e analisar o arquivo de log
	// Por enquanto, retorna instruções
	result.WriteString("⚠️ Análise de log do listener requer acesso ao arquivo de log.\n")
	result.WriteString("O arquivo geralmente está em: `$ORACLE_BASE/diag/tnslsnr/<hostname>/listener/trace/listener.log`\n\n")

	result.WriteString("Use o comando com --log-path para analisar o arquivo:\n")
	result.WriteString("```bash\n")
	result.WriteString("snip db-analysis create --title \"Análise Log Listener\" --db-type oracle --analysis-type logs --log-path /path/to/listener.log\n")
	result.WriteString("```\n")

	return result.String(), nil
}

// AnalyzeRACLatency analisa latência e tempo de resposta do RAC
func (r *RACAnalyzer) AnalyzeRACLatency() (string, error) {
	var result strings.Builder

	result.WriteString("# Análise de Latência e Tempo de Resposta do RAC\n\n")

	// Latência de interconnect
	result.WriteString("## Latência de Interconnect\n\n")
	interconnectLatency, err := r.getInterconnectLatency()
	if err == nil {
		result.WriteString(interconnectLatency)
	}

	// Tempo de resposta por instância
	result.WriteString("\n## Tempo de Resposta por Instância\n\n")
	responseTime, err := r.getInstanceResponseTime()
	if err == nil {
		result.WriteString(responseTime)
	}

	// Estatísticas de cache fusion
	result.WriteString("\n## Estatísticas de Cache Fusion\n\n")
	cacheFusion, err := r.getCacheFusionStats()
	if err == nil {
		result.WriteString(cacheFusion)
	}

	// Blocking entre instâncias
	result.WriteString("\n## Bloqueios entre Instâncias\n\n")
	blocking, err := r.getRACBlocking()
	if err == nil {
		result.WriteString(blocking)
	}

	return result.String(), nil
}

// isRACDatabase verifica se o banco é RAC
func (r *RACAnalyzer) isRACDatabase() (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM v$instance
		WHERE instance_name IS NOT NULL
	`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, err
	}

	// Verificar se há múltiplas instâncias ou se é RAC
	query = `
		SELECT value 
		FROM v$parameter 
		WHERE name = 'cluster_database'
	`

	var value string
	err = r.db.QueryRow(query).Scan(&value)
	if err != nil {
		// Tentar outra forma
		query = `
			SELECT COUNT(DISTINCT instance_number)
			FROM gv$instance
		`
		var instanceCount int
		err = r.db.QueryRow(query).Scan(&instanceCount)
		if err == nil && instanceCount > 1 {
			return true, nil
		}
		return false, nil
	}

	return strings.ToUpper(value) == "TRUE", nil
}

// getClusterInfo obtém informações do cluster
func (r *RACAnalyzer) getClusterInfo() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			instance_name,
			host_name,
			instance_number,
			version,
			status,
			database_status,
			instance_role
		FROM gv$instance
		ORDER BY instance_number
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Instância | Host | Número | Versão | Status | DB Status | Role |\n")
	result.WriteString("|-----------|------|--------|--------|--------|-----------|------|\n")

	for rows.Next() {
		var instanceName, hostName, version, status, dbStatus, role string
		var instanceNumber int

		err := rows.Scan(&instanceName, &hostName, &instanceNumber, &version, &status, &dbStatus, &role)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s | %s | %s |\n",
			instanceName, hostName, instanceNumber, version, status, dbStatus, role))
	}

	return result.String(), nil
}

// getNodesStatus obtém status dos nós
func (r *RACAnalyzer) getNodesStatus() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			instance_name,
			host_name,
			status,
			database_status,
			thread#,
			startup_time
		FROM gv$instance
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Instância | Host | Status | DB Status | Thread | Início |\n")
	result.WriteString("|---------|-----------|------|--------|-----------|--------|--------|\n")

	for rows.Next() {
		var instID, thread int
		var instanceName, hostName, status, dbStatus string
		var startupTime time.Time

		err := rows.Scan(&instID, &instanceName, &hostName, &status, &dbStatus, &thread, &startupTime)
		if err != nil {
			continue
		}

		statusIcon := "✅"
		if status != "OPEN" {
			statusIcon = "❌"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s %s | %s | %d | %s |\n",
			instID, instanceName, hostName, statusIcon, status, dbStatus, thread, startupTime.Format("2006-01-02 15:04:05")))
	}

	return result.String(), nil
}

// getServicesStatus obtém status dos serviços
func (r *RACAnalyzer) getServicesStatus() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			name,
			goal,
			clb_goal,
			aq_ha_notification,
			enabled,
			status
		FROM gv$services
		ORDER BY inst_id, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Serviço | Goal | CLB Goal | AQ HA | Habilitado | Status |\n")
	result.WriteString("|---------|---------|------|----------|-------|------------|--------|\n")

	for rows.Next() {
		var instID int
		var name, goal, clbGoal, aqHA, enabled, status string

		err := rows.Scan(&instID, &name, &goal, &clbGoal, &aqHA, &enabled, &status)
		if err != nil {
			continue
		}

		statusIcon := "✅"
		if status != "READY" {
			statusIcon = "⚠️"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s %s |\n",
			instID, name, goal, clbGoal, aqHA, enabled, statusIcon, status))
	}

	return result.String(), nil
}

// getResourcesStatus obtém status dos recursos do cluster
func (r *RACAnalyzer) getResourcesStatus() (string, error) {
	var result strings.Builder

	// Verificar recursos via v$cluster_interconnects
	query := `
		SELECT 
			inst_id,
			name,
			ip_address,
			is_public,
			is_standby
		FROM gv$cluster_interconnects
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		// Se não disponível, tentar outra query
		result.WriteString("Informações de recursos do cluster não disponíveis via SQL.\n")
		result.WriteString("Use `crsctl status resource -t` no servidor para verificar recursos do clusterware.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("## Interconnects\n\n")
	result.WriteString("| Inst ID | Nome | IP | Público | Standby |\n")
	result.WriteString("|---------|------|----|---------|---------|\n")

	for rows.Next() {
		var instID int
		var name, ipAddress string
		var isPublic, isStandby sql.NullString

		err := rows.Scan(&instID, &name, &ipAddress, &isPublic, &isStandby)
		if err != nil {
			continue
		}

		public := "Não"
		if isPublic.Valid && isPublic.String == "Y" {
			public = "Sim"
		}

		standby := "Não"
		if isStandby.Valid && isStandby.String == "Y" {
			standby = "Sim"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
			instID, name, ipAddress, public, standby))
	}

	return result.String(), nil
}

// getInstanceErrors obtém erros de instância
func (r *RACAnalyzer) getInstanceErrors() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			error_number,
			error_facility,
			error_message,
			error_timestamp,
			error_count
		FROM gv$diag_alert_ext
		WHERE error_number IS NOT NULL
		AND error_timestamp >= SYSDATE - 1
		ORDER BY error_timestamp DESC, inst_id
		FETCH FIRST 50 ROWS ONLY
	`

	rows, err := r.db.Query(query)
	if err != nil {
		// Se não disponível, usar outra abordagem
		result.WriteString("Verificando erros via v$instance...\n\n")
		return r.getBasicInstanceErrors()
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Código | Facilidade | Mensagem | Timestamp | Count |\n")
	result.WriteString("|---------|--------|------------|----------|-----------|-------|\n")

	count := 0
	for rows.Next() {
		var instID, errorNumber, errorCount int
		var errorFacility, errorMessage sql.NullString
		var errorTimestamp time.Time

		err := rows.Scan(&instID, &errorNumber, &errorFacility, &errorMessage, &errorTimestamp, &errorCount)
		if err != nil {
			continue
		}

		msg := "N/A"
		if errorMessage.Valid {
			if len(errorMessage.String) > 50 {
				msg = errorMessage.String[:50] + "..."
			} else {
				msg = errorMessage.String
			}
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %s | %s | %s | %d |\n",
			instID, errorNumber, getString(errorFacility), msg, errorTimestamp.Format("2006-01-02 15:04:05"), errorCount))

		count++
		if count >= 50 {
			break
		}
	}

	if count == 0 {
		result.WriteString("✅ Nenhum erro de instância encontrado nas últimas 24 horas.\n")
	}

	return result.String(), nil
}

// getBasicInstanceErrors obtém erros básicos de instância
func (r *RACAnalyzer) getBasicInstanceErrors() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			instance_name,
			status,
			database_status
		FROM gv$instance
		WHERE status != 'OPEN' OR database_status != 'ACTIVE'
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		result.WriteString("✅ Todas as instâncias estão com status OPEN e ACTIVE.\n")
		return result.String(), nil
	}

	result.WriteString("⚠️ Instâncias com problemas:\n\n")
	result.WriteString("| Inst ID | Instância | Status | DB Status |\n")
	result.WriteString("|---------|-----------|--------|-----------|\n")

	for {
		var instID int
		var instanceName, status, dbStatus string

		err := rows.Scan(&instID, &instanceName, &status, &dbStatus)
		if err != nil {
			break
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", instID, instanceName, status, dbStatus))

		if !rows.Next() {
			break
		}
	}

	return result.String(), nil
}

// getClusterwareErrors obtém erros do clusterware
func (r *RACAnalyzer) getClusterwareErrors() (string, error) {
	var result strings.Builder

	result.WriteString("⚠️ Erros de clusterware requerem acesso ao clusterware (crsctl, olsnodes, etc.).\n")
	result.WriteString("Execute no servidor:\n")
	result.WriteString("```bash\n")
	result.WriteString("crsctl check cluster\n")
	result.WriteString("crsctl status resource -t\n")
	result.WriteString("olsnodes -n\n")
	result.WriteString("```\n\n")

	// Tentar obter informações via views do banco
	query := `
		SELECT 
			inst_id,
			instance_name,
			host_name,
			status
		FROM gv$instance
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("Status das instâncias (indicador de saúde do clusterware):\n\n")
	result.WriteString("| Inst ID | Instância | Host | Status |\n")
	result.WriteString("|---------|-----------|------|--------|\n")

	for rows.Next() {
		var instID int
		var instanceName, hostName, status string

		err := rows.Scan(&instID, &instanceName, &hostName, &status)
		if err != nil {
			continue
		}

		statusIcon := "✅"
		if status != "OPEN" {
			statusIcon = "❌"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s %s |\n", instID, instanceName, hostName, statusIcon, status))
	}

	return result.String(), nil
}

// getInterconnectErrors obtém erros de interconnect
func (r *RACAnalyzer) getInterconnectErrors() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			name,
			ip_address,
			is_public,
			is_standby
		FROM gv$cluster_interconnects
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		result.WriteString("⚠️ Informações de interconnect não disponíveis.\n")
		result.WriteString("Verifique a conectividade de rede entre os nós do cluster.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Nome | IP | Público | Standby |\n")
	result.WriteString("|---------|------|----|---------|---------|\n")

	for rows.Next() {
		var instID int
		var name, ipAddress string
		var isPublic, isStandby sql.NullString

		err := rows.Scan(&instID, &name, &ipAddress, &isPublic, &isStandby)
		if err != nil {
			continue
		}

		public := "Não"
		if isPublic.Valid && isPublic.String == "Y" {
			public = "Sim"
		}

		standby := "Não"
		if isStandby.Valid && isStandby.String == "Y" {
			standby = "Sim"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
			instID, name, ipAddress, public, standby))
	}

	// Verificar estatísticas de interconnect
	result.WriteString("\n## Estatísticas de Interconnect\n\n")
	statsQuery := `
		SELECT 
			inst_id,
			network_name,
			bytes_sent,
			bytes_received,
			bytes_sent_delta,
			bytes_received_delta
		FROM gv$cluster_network
		ORDER BY inst_id
	`

	statsRows, err := r.db.Query(statsQuery)
	if err == nil {
		defer statsRows.Close()
		result.WriteString("| Inst ID | Network | Bytes Enviados | Bytes Recebidos |\n")
		result.WriteString("|---------|---------|----------------|-----------------|\n")

		for statsRows.Next() {
			var instID int
			var networkName string
			var bytesSent, bytesReceived, bytesSentDelta, bytesReceivedDelta sql.NullInt64

			err := statsRows.Scan(&instID, &networkName, &bytesSent, &bytesReceived, &bytesSentDelta, &bytesReceivedDelta)
			if err != nil {
				continue
			}

			result.WriteString(fmt.Sprintf("| %d | %s | %d | %d |\n",
				instID, networkName, getInt64(bytesSent), getInt64(bytesReceived)))
		}
	}

	return result.String(), nil
}

// getRACDeadlocks obtém deadlocks entre instâncias
func (r *RACAnalyzer) getRACDeadlocks() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			deadlock_id,
			deadlock_type,
			deadlock_graph
		FROM gv$deadlock_history
		WHERE deadlock_timestamp >= SYSDATE - 7
		ORDER BY deadlock_timestamp DESC
		FETCH FIRST 20 ROWS ONLY
	`

	rows, err := r.db.Query(query)
	if err != nil {
		result.WriteString("⚠️ Informações de deadlock não disponíveis ou nenhum deadlock recente.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Deadlock ID | Tipo | Detalhes |\n")
	result.WriteString("|---------|-------------|------|----------|\n")

	count := 0
	for rows.Next() {
		var instID, deadlockID int
		var deadlockType sql.NullString
		var deadlockGraph sql.NullString

		err := rows.Scan(&instID, &deadlockID, &deadlockType, &deadlockGraph)
		if err != nil {
			continue
		}

		details := "N/A"
		if deadlockGraph.Valid {
			if len(deadlockGraph.String) > 50 {
				details = deadlockGraph.String[:50] + "..."
			} else {
				details = deadlockGraph.String
			}
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %s | %s |\n",
			instID, deadlockID, getString(deadlockType), details))

		count++
	}

	if count == 0 {
		result.WriteString("✅ Nenhum deadlock encontrado nos últimos 7 dias.\n")
	}

	return result.String(), nil
}

// getListenerStatus obtém status do listener
func (r *RACAnalyzer) getListenerStatus() (string, error) {
	var result strings.Builder

	// Verificar serviços registrados (indica que listener está funcionando)
	query := `
		SELECT 
			inst_id,
			name,
			network_name,
			goal,
			status
		FROM gv$services
		WHERE name NOT LIKE 'SYS%'
		ORDER BY inst_id, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("⚠️ Status detalhado do listener requer execução de `lsnrctl status` no servidor.\n")
	result.WriteString("Verificando serviços registrados (indicador de listener ativo):\n\n")

	result.WriteString("| Inst ID | Serviço | Network | Goal | Status |\n")
	result.WriteString("|---------|---------|---------|------|--------|\n")

	activeListeners := make(map[int]bool)
	for rows.Next() {
		var instID int
		var name, networkName, goal, status string

		err := rows.Scan(&instID, &name, &networkName, &goal, &status)
		if err != nil {
			continue
		}

		if status == "READY" {
			activeListeners[instID] = true
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
			instID, name, networkName, goal, status))
	}

	result.WriteString("\n### Resumo\n\n")
	for instID, active := range activeListeners {
		if active {
			result.WriteString(fmt.Sprintf("✅ Instância %d: Listener parece estar ativo (serviços READY)\n", instID))
		} else {
			result.WriteString(fmt.Sprintf("⚠️ Instância %d: Listener pode ter problemas (nenhum serviço READY)\n", instID))
		}
	}

	result.WriteString("\n💡 Para verificação completa, execute no servidor:\n")
	result.WriteString("```bash\n")
	result.WriteString("lsnrctl status\n")
	result.WriteString("lsnrctl status LISTENER\n")
	result.WriteString("```\n")

	return result.String(), nil
}

// getRegisteredServices obtém serviços registrados no listener
func (r *RACAnalyzer) getRegisteredServices() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			name,
			network_name,
			goal,
			clb_goal,
			aq_ha_notification,
			enabled,
			status
		FROM gv$services
		ORDER BY inst_id, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Serviço | Network | Goal | Status | Habilitado |\n")
	result.WriteString("|---------|---------|---------|------|--------|------------|\n")

	for rows.Next() {
		var instID int
		var name, networkName, goal, clbGoal, aqHA, enabled, status string

		err := rows.Scan(&instID, &name, &networkName, &goal, &clbGoal, &aqHA, &enabled, &status)
		if err != nil {
			continue
		}

		statusIcon := "✅"
		if status != "READY" {
			statusIcon = "⚠️"
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s %s | %s |\n",
			instID, name, networkName, goal, statusIcon, status, enabled))
	}

	return result.String(), nil
}

// getActiveConnections obtém conexões ativas
func (r *RACAnalyzer) getActiveConnections() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			COUNT(*) as connection_count,
			program,
			status
		FROM gv$session
		WHERE username IS NOT NULL
		GROUP BY inst_id, program, status
		ORDER BY inst_id, connection_count DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Conexões | Programa | Status |\n")
	result.WriteString("|---------|----------|----------|--------|\n")

	for rows.Next() {
		var instID, connectionCount int
		var program, status sql.NullString

		err := rows.Scan(&instID, &connectionCount, &program, &status)
		if err != nil {
			continue
		}

		prog := "N/A"
		if program.Valid {
			if len(program.String) > 30 {
				prog = program.String[:30] + "..."
			} else {
				prog = program.String
			}
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %s | %s |\n",
			instID, connectionCount, prog, getString(status)))
	}

	return result.String(), nil
}

// getListenerErrors obtém erros do listener
func (r *RACAnalyzer) getListenerErrors() (string, error) {
	var result strings.Builder

	result.WriteString("⚠️ Erros detalhados do listener requerem análise do arquivo de log.\n")
	result.WriteString("Use o comando com --log-path para analisar o log do listener.\n\n")

	result.WriteString("Verificando problemas indiretos (sessões com erros):\n\n")

	query := `
		SELECT 
			inst_id,
			COUNT(*) as error_count,
			error#
		FROM gv$session
		WHERE error# IS NOT NULL AND error# != 0
		GROUP BY inst_id, error#
		ORDER BY inst_id, error_count DESC
		FETCH FIRST 20 ROWS ONLY
	`

	rows, err := r.db.Query(query)
	if err != nil {
		result.WriteString("Nenhuma informação de erro disponível via SQL.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Erros | Código de Erro |\n")
	result.WriteString("|---------|-------|----------------|\n")

	count := 0
	for rows.Next() {
		var instID, errorCount, errorCode int

		err := rows.Scan(&instID, &errorCount, &errorCode)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %d | %d |\n", instID, errorCount, errorCode))
		count++
	}

	if count == 0 {
		result.WriteString("✅ Nenhum erro de sessão detectado.\n")
	}

	return result.String(), nil
}

// getInterconnectLatency obtém latência de interconnect
func (r *RACAnalyzer) getInterconnectLatency() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			network_name,
			round_trip_ms,
			bytes_sent,
			bytes_received
		FROM gv$cluster_network
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		result.WriteString("⚠️ Informações de latência de interconnect não disponíveis via SQL.\n")
		result.WriteString("Use `ping` ou ferramentas de rede para medir latência entre nós.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Network | Round Trip (ms) | Bytes Enviados | Bytes Recebidos |\n")
	result.WriteString("|---------|---------|-----------------|----------------|-----------------|\n")

	for rows.Next() {
		var instID int
		var networkName string
		var roundTripMS sql.NullFloat64
		var bytesSent, bytesReceived sql.NullInt64

		err := rows.Scan(&instID, &networkName, &roundTripMS, &bytesSent, &bytesReceived)
		if err != nil {
			continue
		}

		latency := "N/A"
		if roundTripMS.Valid {
			latency = fmt.Sprintf("%.2f", roundTripMS.Float64)
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %d | %d |\n",
			instID, networkName, latency, getInt64(bytesSent), getInt64(bytesReceived)))
	}

	return result.String(), nil
}

// getInstanceResponseTime obtém tempo de resposta por instância
func (r *RACAnalyzer) getInstanceResponseTime() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			AVG(elapsed_time) / 1000000 as avg_elapsed_ms,
			MAX(elapsed_time) / 1000000 as max_elapsed_ms,
			MIN(elapsed_time) / 1000000 as min_elapsed_ms,
			COUNT(*) as query_count
		FROM gv$sqlstats
		WHERE elapsed_time > 0
		GROUP BY inst_id
		ORDER BY inst_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Avg (ms) | Max (ms) | Min (ms) | Queries |\n")
	result.WriteString("|---------|----------|----------|----------|---------|\n")

	for rows.Next() {
		var instID, queryCount int
		var avgElapsed, maxElapsed, minElapsed sql.NullFloat64

		err := rows.Scan(&instID, &avgElapsed, &maxElapsed, &minElapsed, &queryCount)
		if err != nil {
			continue
		}

		avg := "N/A"
		if avgElapsed.Valid {
			avg = fmt.Sprintf("%.2f", avgElapsed.Float64)
		}

		max := "N/A"
		if maxElapsed.Valid {
			max = fmt.Sprintf("%.2f", maxElapsed.Float64)
		}

		min := "N/A"
		if minElapsed.Valid {
			min = fmt.Sprintf("%.2f", minElapsed.Float64)
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %d |\n",
			instID, avg, max, min, queryCount))
	}

	return result.String(), nil
}

// getCacheFusionStats obtém estatísticas de cache fusion
func (r *RACAnalyzer) getCacheFusionStats() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			inst_id,
			stat_name,
			value
		FROM gv$sysstat
		WHERE stat_name LIKE '%cache fusion%' 
		   OR stat_name LIKE '%gc%'
		   OR stat_name LIKE '%global cache%'
		ORDER BY inst_id, stat_name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result.WriteString("| Inst ID | Estatística | Valor |\n")
	result.WriteString("|---------|-------------|-------|\n")

	for rows.Next() {
		var instID int
		var statName string
		var value sql.NullInt64

		err := rows.Scan(&instID, &statName, &value)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d | %s | %d |\n",
			instID, statName, getInt64(value)))
	}

	return result.String(), nil
}

// getRACBlocking obtém bloqueios entre instâncias
func (r *RACAnalyzer) getRACBlocking() (string, error) {
	var result strings.Builder

	query := `
		SELECT 
			blocking_inst_id,
			blocking_session,
			blocked_inst_id,
			blocked_session,
			blocking_status,
			blocking_type
		FROM gv$lock
		WHERE blocking_inst_id IS NOT NULL
		  AND blocking_inst_id != blocked_inst_id
		ORDER BY blocking_inst_id
		FETCH FIRST 20 ROWS ONLY
	`

	rows, err := r.db.Query(query)
	if err != nil {
		result.WriteString("⚠️ Informações de bloqueio entre instâncias não disponíveis.\n")
		return result.String(), nil
	}
	defer rows.Close()

	result.WriteString("| Bloqueando (Inst:Session) | Bloqueado (Inst:Session) | Status | Tipo |\n")
	result.WriteString("|---------------------------|--------------------------|--------|------|\n")

	count := 0
	for rows.Next() {
		var blockingInstID, blockingSession, blockedInstID, blockedSession int
		var blockingStatus, blockingType sql.NullString

		err := rows.Scan(&blockingInstID, &blockingSession, &blockedInstID, &blockedSession,
			&blockingStatus, &blockingType)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf("| %d:%d | %d:%d | %s | %s |\n",
			blockingInstID, blockingSession, blockedInstID, blockedSession,
			getString(blockingStatus), getString(blockingType)))

		count++
	}

	if count == 0 {
		result.WriteString("✅ Nenhum bloqueio entre instâncias detectado.\n")
	}

	return result.String(), nil
}

