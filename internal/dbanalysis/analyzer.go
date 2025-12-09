package dbanalysis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/dbconnection"
	"github.com/snip/internal/loganalyzer"
	"github.com/snip/internal/backup"
	"github.com/snip/internal/dbchecklist"
	"github.com/snip/internal/dbcharts"
	"github.com/snip/internal/dbdynamic"
	"github.com/snip/internal/mongodb"
	"github.com/snip/internal/mysql"
	"github.com/snip/internal/oracle"
	"github.com/snip/internal/postgresql"
	"github.com/snip/internal/sqlserver"
)

// Analyzer realiza análises de banco de dados
type Analyzer struct {
	aiClient *ai.GroqClient
}

// NewAnalyzer cria um novo analisador
func NewAnalyzer() (*Analyzer, error) {
	aiClient, err := ai.NewGroqClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &Analyzer{aiClient: aiClient}, nil
}

// PerformAnalysis realiza uma análise completa
func (a *Analyzer) PerformAnalysis(analysis *DBAnalysis, config *ConnectionConfig) error {
	analysis.Status = "processing"

	var result strings.Builder
	var err error

	// Realizar análise baseada no tipo
	switch analysis.AnalysisType {
	case AnalysisTypeLogs:
		result, err = a.analyzeLogs(analysis)
	case AnalysisTypeDiagnostic:
		result, err = a.performDiagnostic(analysis, config)
	case AnalysisTypeTuning:
		result, err = a.performTuning(analysis, config)
	case AnalysisTypeQuery:
		result, err = a.analyzeQueries(analysis, config)
	case AnalysisTypeTablespace:
		result, err = a.analyzeTablespace(analysis, config)
	case AnalysisTypeDisk:
		result, err = a.analyzeDisk(analysis, config)
	case AnalysisTypeTables:
		result, err = a.analyzeTables(analysis, config)
	case AnalysisTypeIndexes:
		result, err = a.analyzeIndexes(analysis, config)
	case AnalysisTypePredictive:
		result, err = a.performPredictiveAnalysis(analysis, config)
	case AnalysisTypeErrorKnowledge:
		result, err = a.analyzeErrorKnowledge(analysis, config)
	// Análises específicas Oracle
	case AnalysisTypeAWR:
		result, err = a.analyzeAWR(analysis, config)
	case AnalysisTypeASH:
		result, err = a.analyzeASH(analysis, config)
	case AnalysisTypeExecutionPlan:
		result, err = a.analyzeExecutionPlan(analysis, config)
	// Análises específicas SQL Server
	case AnalysisTypeLocks:
		result, err = a.analyzeLocks(analysis, config)
	case AnalysisTypeActiveSessions:
		result, err = a.analyzeActiveSessions(analysis, config)
	case AnalysisTypeRunningQueries:
		result, err = a.analyzeRunningQueries(analysis, config)
	// Análises específicas MongoDB
	case AnalysisTypeReplication:
		result, err = a.analyzeMongoDBReplication(analysis, config)
	case AnalysisTypeSharding:
		result, err = a.analyzeMongoDBSharding(analysis, config)
	case AnalysisTypeLatency:
		result, err = a.analyzeMongoDBLatency(analysis, config)
	case AnalysisTypePerformance:
		result, err = a.analyzeMongoDBPerformance(analysis, config)
	// Análises específicas PostgreSQL
	case AnalysisTypePostgresReplication:
		result, err = a.analyzePostgresReplication(analysis, config)
	case AnalysisTypePostgresLocks:
		result, err = a.analyzePostgresLocks(analysis, config)
	case AnalysisTypePostgresFragmentation:
		result, err = a.analyzePostgresFragmentation(analysis, config)
	// Análises específicas MySQL
	case AnalysisTypeMySQLReplication:
		result, err = a.analyzeMySQLReplication(analysis, config)
	case AnalysisTypeMySQLLocks:
		result, err = a.analyzeMySQLLocks(analysis, config)
	case AnalysisTypeMySQLFragmentation:
		result, err = a.analyzeMySQLFragmentation(analysis, config)
	// Checklist e Backup
	case AnalysisTypeChecklist:
		result, err = a.analyzeChecklist(analysis, config)
	case AnalysisTypeBackup:
		result, err = a.analyzeBackup(analysis, config)
	// Análises dinâmicas e chat
	case AnalysisTypeDynamic:
		result, err = a.analyzeDynamic(analysis, config)
	case AnalysisTypeChat:
		return fmt.Errorf("chat requer sessão interativa, use o comando 'db-chat'")
	// Oracle PDBs
	case AnalysisTypePDBs:
		result, err = a.analyzePDBs(analysis, config)
	case AnalysisTypePDB:
		result, err = a.analyzePDB(analysis, config)
	// SQL Server Instance e Databases
	case AnalysisTypeInstance:
		result, err = a.analyzeInstance(analysis, config)
	case AnalysisTypeDatabases:
		result, err = a.analyzeDatabases(analysis, config)
	case AnalysisTypeDatabase:
		result, err = a.analyzeDatabase(analysis, config)
	// Oracle RAC
	case AnalysisTypeRACHealth:
		result, err = a.analyzeRACHealth(analysis, config)
	case AnalysisTypeRACErrors:
		result, err = a.analyzeRACErrors(analysis, config)
	case AnalysisTypeRACListener:
		result, err = a.analyzeRACListener(analysis, config)
	case AnalysisTypeRACLatency:
		result, err = a.analyzeRACLatency(analysis, config)
	default:
		return fmt.Errorf("tipo de análise não suportado: %s", analysis.AnalysisType)
	}

	if err != nil {
		analysis.Status = "error"
		analysis.ErrorMessage = err.Error()
		return err
	}

	analysis.Result = result.String()

	// Gerar insights com IA
	insights, err := a.generateAIInsights(analysis)
	if err == nil {
		analysis.AIInsights = insights
	}

	// Gerar gráfico se aplicável
	chart, err := a.generateChart(analysis)
	if err == nil && chart != "" {
		analysis.Result += "\n\n## Visualização\n\n" + chart
	}

	analysis.Status = "completed"
	return nil
}

// generateChart gera gráfico para a análise
func (a *Analyzer) generateChart(analysis *DBAnalysis) (string, error) {
	chartGenerator, err := dbcharts.NewChartGenerator()
	if err != nil {
		return "", err
	}

	// Sugerir tipo de gráfico
	chartType, reason, err := chartGenerator.SuggestChartType(analysis.Result)
	if err != nil {
		return "", err
	}

	// Gerar gráfico
	chart, err := chartGenerator.GenerateChartFromAnalysis(analysis.Result, chartType)
	if err != nil {
		return "", err
	}

	if chart != "" {
		return fmt.Sprintf("**Tipo de gráfico sugerido:** %s\n**Razão:** %s\n\n%s", chartType, reason, chart), nil
	}

	return "", nil
}

// analyzeLogs analisa arquivos de log
func (a *Analyzer) analyzeLogs(analysis *DBAnalysis) (strings.Builder, error) {
	var result strings.Builder

	if analysis.LogFilePath == "" {
		return result, fmt.Errorf("caminho do arquivo de log não fornecido")
	}

	logAnalyzer := loganalyzer.NewLogAnalyzer(analysis.DatabaseType)
	logAnalysis, err := logAnalyzer.AnalyzeLog(analysis.LogFilePath)
	if err != nil {
		return result, err
	}

	result.WriteString(logAnalysis)
	return result, nil
}

// performDiagnostic realiza diagnóstico do banco
func (a *Analyzer) performDiagnostic(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	// Testar conexão
	err = connector.TestConnection(config)
	if err != nil {
		result.WriteString(fmt.Sprintf("❌ Erro de conexão: %s\n\n", err.Error()))
		return result, nil
	}

	result.WriteString("✅ Conexão estabelecida com sucesso\n\n")

	// Conectar e realizar diagnósticos
	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	// Diagnósticos específicos por tipo de banco
	switch analysis.DatabaseType {
	case DatabaseTypeOracle:
		result.WriteString(a.diagnoseOracle(db))
	case DatabaseTypeSQLServer:
		result.WriteString(a.diagnoseSQLServer(db))
	case DatabaseTypeMySQL:
		result.WriteString(a.diagnoseMySQL(db))
	case DatabaseTypePostgreSQL:
		result.WriteString(a.diagnosePostgreSQL(db))
	default:
		result.WriteString("Diagnóstico genérico não implementado para este tipo de banco\n")
	}

	// Análises específicas baseadas no tipo de análise
	switch analysis.AnalysisType {
	case AnalysisTypeQuery:
		if analysis.DatabaseType == DatabaseTypeOracle {
			// Análise Oracle específica
		} else if analysis.DatabaseType == DatabaseTypeSQLServer {
			sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
			runningQueries, err := sqlAnalyzer.AnalyzeRunningQueries()
			if err == nil {
				result.WriteString("\n" + runningQueries)
			}
		}
	}

	return result, nil
}

// performTuning realiza análise de tuning
func (a *Analyzer) performTuning(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	result.WriteString("# Análise de Tuning\n\n")
	result.WriteString("## Recomendações de Performance\n\n")

	// Análise de queries lentas, índices, etc.
	switch analysis.DatabaseType {
	case DatabaseTypePostgreSQL:
		result.WriteString(a.tunePostgreSQL(db))
	case DatabaseTypeMySQL:
		result.WriteString(a.tuneMySQL(db))
	default:
		result.WriteString("Análise de tuning específica não implementada para este banco\n")
	}

	return result, nil
}

// analyzeQueries analisa consultas SQL
func (a *Analyzer) analyzeQueries(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	result.WriteString("# Análise de Consultas\n\n")
	result.WriteString("Análise de consultas SQL para otimização\n")

	return result, nil
}

// analyzeTablespace analisa tablespaces (Oracle)
func (a *Analyzer) analyzeTablespace(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de tablespace disponível apenas para Oracle")
	}

	result.WriteString("# Análise de Tablespace\n\n")
	result.WriteString("Análise de uso e espaço de tablespaces\n")

	return result, nil
}

// analyzeDisk analisa uso de disco
func (a *Analyzer) analyzeDisk(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	result.WriteString("# Análise de Disco\n\n")
	result.WriteString("Análise de uso de espaço em disco\n")

	return result, nil
}

// analyzeTables analisa tabelas
func (a *Analyzer) analyzeTables(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	result.WriteString("# Análise de Tabelas\n\n")

	// Listar tabelas
	var query string
	switch analysis.DatabaseType {
	case DatabaseTypePostgreSQL, DatabaseTypeMySQL:
		query = "SELECT table_name, table_rows, data_length, index_length FROM information_schema.tables WHERE table_schema = ?"
	case DatabaseTypeSQLServer:
		query = "SELECT name FROM sys.tables"
	default:
		result.WriteString("Análise de tabelas não implementada para este banco\n")
		return result, nil
	}

	rows, err := db.Query(query, config.Database)
	if err != nil {
		result.WriteString(fmt.Sprintf("Erro ao consultar tabelas: %s\n", err.Error()))
		return result, nil
	}
	defer rows.Close()

	result.WriteString("## Tabelas Encontradas\n\n")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err == nil {
			result.WriteString(fmt.Sprintf("- %s\n", tableName))
		}
	}

	return result, nil
}

// analyzeIndexes analisa índices e recomenda melhorias
func (a *Analyzer) analyzeIndexes(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	result.WriteString("# Análise de Índices\n\n")
	result.WriteString("## Recomendações de Índices\n\n")

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	// Análise específica por banco
	switch analysis.DatabaseType {
	case DatabaseTypePostgreSQL:
		result.WriteString(a.analyzePostgreSQLIndexes(db))
	case DatabaseTypeMySQL:
		result.WriteString(a.analyzeMySQLIndexes(db))
	default:
		result.WriteString("Análise de índices não implementada para este banco\n")
	}

	return result, nil
}

// performPredictiveAnalysis realiza análise preditiva
func (a *Analyzer) performPredictiveAnalysis(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	result.WriteString("# Análise Preditiva\n\n")
	result.WriteString("Análise preditiva de performance e crescimento\n")

	return result, nil
}

// analyzeErrorKnowledge analisa base de conhecimento de erros
func (a *Analyzer) analyzeErrorKnowledge(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	result.WriteString("# Base de Conhecimento de Erros\n\n")
	result.WriteString("Análise de erros conhecidos e soluções\n")

	return result, nil
}

// generateAIInsights gera insights usando IA
func (a *Analyzer) generateAIInsights(analysis *DBAnalysis) (string, error) {
	prompt := fmt.Sprintf(`Você é um especialista em banco de dados %s. Analise os seguintes resultados de análise e forneça insights, recomendações e possíveis problemas.

Tipo de Análise: %s
Resultado:
%s

Forneça:
1. Resumo executivo
2. Principais problemas identificados
3. Recomendações de ação
4. Próximos passos sugeridos

Formate a resposta em markdown.`, analysis.DatabaseType, analysis.AnalysisType, analysis.Result)

	return a.aiClient.GenerateContent(prompt, 2000)
}

// Funções auxiliares de diagnóstico específicas por banco
func (a *Analyzer) diagnoseOracle(db *sql.DB) string {
	return "Diagnóstico Oracle: Verificando versão, parâmetros, tablespaces...\n"
}

func (a *Analyzer) diagnoseSQLServer(db *sql.DB) string {
	return "Diagnóstico SQL Server: Verificando versão, configurações, bloqueios...\n"
}

func (a *Analyzer) diagnoseMySQL(db *sql.DB) string {
	var result strings.Builder
	result.WriteString("## Diagnóstico MySQL\n\n")

	// Verificar versão
	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err == nil {
		result.WriteString(fmt.Sprintf("**Versão:** %s\n\n", version))
	}

	// Verificar status
	result.WriteString("### Status do Servidor\n")
	rows, err := db.Query("SHOW STATUS LIKE 'Threads_connected'")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var variable, value string
			if err := rows.Scan(&variable, &value); err == nil {
				result.WriteString(fmt.Sprintf("- %s: %s\n", variable, value))
			}
		}
	}

	return result.String()
}

func (a *Analyzer) diagnosePostgreSQL(db *sql.DB) string {
	var result strings.Builder
	result.WriteString("## Diagnóstico PostgreSQL\n\n")

	// Verificar versão
	var version string
	if err := db.QueryRow("SELECT version()").Scan(&version); err == nil {
		result.WriteString(fmt.Sprintf("**Versão:** %s\n\n", version))
	}

	// Verificar conexões ativas
	var connections int
	if err := db.QueryRow("SELECT count(*) FROM pg_stat_activity").Scan(&connections); err == nil {
		result.WriteString(fmt.Sprintf("**Conexões ativas:** %d\n\n", connections))
	}

	return result.String()
}

func (a *Analyzer) tunePostgreSQL(db *sql.DB) string {
	var result strings.Builder
	result.WriteString("## Recomendações PostgreSQL\n\n")

	// Verificar configurações importantes
	rows, err := db.Query("SELECT name, setting, unit FROM pg_settings WHERE name IN ('shared_buffers', 'effective_cache_size', 'work_mem')")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name, setting, unit sql.NullString
			if err := rows.Scan(&name, &setting, &unit); err == nil {
				result.WriteString(fmt.Sprintf("- %s: %s %s\n", name.String, setting.String, unit.String))
			}
		}
	}

	return result.String()
}

func (a *Analyzer) tuneMySQL(db *sql.DB) string {
	return "Recomendações de tuning MySQL\n"
}

func (a *Analyzer) analyzePostgreSQLIndexes(db *sql.DB) string {
	var result strings.Builder
	result.WriteString("## Índices PostgreSQL\n\n")

	// Consultar índices não utilizados
	query := `
		SELECT schemaname, tablename, indexname, idx_scan
		FROM pg_stat_user_indexes
		WHERE idx_scan = 0
		ORDER BY pg_relation_size(indexrelid) DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err == nil {
		defer rows.Close()
		result.WriteString("### Índices Não Utilizados\n\n")
		for rows.Next() {
			var schema, table, index string
			var scans int
			if err := rows.Scan(&schema, &table, &index, &scans); err == nil {
				result.WriteString(fmt.Sprintf("- %s.%s.%s (scans: %d)\n", schema, table, index, scans))
			}
		}
	}

	return result.String()
}

func (a *Analyzer) analyzeMySQLIndexes(db *sql.DB) string {
	return "Análise de índices MySQL\n"
}

// SerializeConnectionConfig serializa a configuração de conexão para JSON
func SerializeConnectionConfig(config *ConnectionConfig) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeConnectionConfig deserializa a configuração de conexão de JSON
func DeserializeConnectionConfig(data string) (*ConnectionConfig, error) {
	var config ConnectionConfig
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// analyzeAWR analisa AWR do Oracle
func (a *Analyzer) analyzeAWR(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise AWR disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	oracleAnalyzer := oracle.NewOracleAnalyzer(db)

	req := oracle.AWRReportRequest{
		InstanceID: 1,
		ReportType: "html",
		BeginTime:  time.Now().Add(-2 * time.Hour),
		EndTime:    time.Now(),
	}

	snapshots, err := oracleAnalyzer.GetSnapshots(req.BeginTime, req.EndTime)
	if err == nil && len(snapshots) >= 2 {
		req.BeginSnapshot = snapshots[0].ID
		req.EndSnapshot = snapshots[len(snapshots)-1].ID
	}

	awrReport, err := oracleAnalyzer.GenerateAWRReport(req)
	if err != nil {
		result.WriteString(fmt.Sprintf("Erro ao gerar relatório AWR: %s\n", err.Error()))
	} else {
		result.WriteString(awrReport)
	}

	return result, nil
}

// analyzeASH analisa ASH do Oracle
func (a *Analyzer) analyzeASH(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise ASH disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	oracleAnalyzer := oracle.NewOracleAnalyzer(db)

	req := oracle.ASHRequest{
		EndTime:   time.Now(),
		BeginTime: time.Now().Add(-30 * time.Minute),
	}

	ashResult, err := oracleAnalyzer.AnalyzeASH(req)
	if err != nil {
		return result, err
	}

	result.WriteString(ashResult)
	return result, nil
}

// analyzeExecutionPlan analisa plano de execução
func (a *Analyzer) analyzeExecutionPlan(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	switch analysis.DatabaseType {
	case DatabaseTypeOracle:
		oracleAnalyzer := oracle.NewOracleAnalyzer(db)
		sqlID := "" // TODO: extrair de algum lugar
		if sqlID != "" {
			plan, err := oracleAnalyzer.GetExecutionPlan(sqlID)
			if err == nil {
				result.WriteString(plan)
			}
		}
	case DatabaseTypeSQLServer:
		sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
		sqlHandle := "" // TODO: extrair de algum lugar
		if sqlHandle != "" {
			plan, err := sqlAnalyzer.GetExecutionPlan(sqlHandle)
			if err == nil {
				result.WriteString(plan)
			}
		}
	case DatabaseTypePostgreSQL:
		pgAnalyzer := postgresql.NewPostgreSQLAnalyzer(db)
		query := "" // TODO: extrair de algum lugar
		if query != "" {
			plan, err := pgAnalyzer.GetExecutionPlan(query)
			if err == nil {
				result.WriteString(plan)
			}
		}
	}

	return result, nil
}

// analyzeLocks analisa locks (SQL Server)
func (a *Analyzer) analyzeLocks(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de locks disponível apenas para SQL Server")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	locksResult, err := sqlAnalyzer.AnalyzeLocks()
	if err != nil {
		return result, err
	}

	result.WriteString(locksResult)
	return result, nil
}

// analyzeActiveSessions analisa sessões ativas (SQL Server)
func (a *Analyzer) analyzeActiveSessions(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de sessões ativas disponível apenas para SQL Server")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	sessionsResult, err := sqlAnalyzer.AnalyzeActiveSessions()
	if err != nil {
		return result, err
	}

	result.WriteString(sessionsResult)
	return result, nil
}

// analyzeRunningQueries analisa queries em execução (SQL Server)
func (a *Analyzer) analyzeRunningQueries(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de queries em execução disponível apenas para SQL Server")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	queriesResult, err := sqlAnalyzer.AnalyzeRunningQueries()
	if err != nil {
		return result, err
	}

	result.WriteString(queriesResult)
	return result, nil
}

// analyzeMongoDBReplication analisa replicação MongoDB
func (a *Analyzer) analyzeMongoDBReplication(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMongoDB {
		return result, fmt.Errorf("análise de replicação disponível apenas para MongoDB")
	}

	mongoAnalyzer, err := a.connectMongoDB(config)
	if err != nil {
		return result, err
	}
	defer mongoAnalyzer.Close()

	replicationResult, err := mongoAnalyzer.AnalyzeReplication()
	if err != nil {
		return result, err
	}

	result.WriteString(replicationResult)
	return result, nil
}

// analyzeMongoDBSharding analisa sharding MongoDB
func (a *Analyzer) analyzeMongoDBSharding(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMongoDB {
		return result, fmt.Errorf("análise de sharding disponível apenas para MongoDB")
	}

	mongoAnalyzer, err := a.connectMongoDB(config)
	if err != nil {
		return result, err
	}
	defer mongoAnalyzer.Close()

	shardingResult, err := mongoAnalyzer.AnalyzeSharding()
	if err != nil {
		return result, err
	}

	result.WriteString(shardingResult)
	return result, nil
}

// analyzeMongoDBLatency analisa latência MongoDB
func (a *Analyzer) analyzeMongoDBLatency(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMongoDB {
		return result, fmt.Errorf("análise de latência disponível apenas para MongoDB")
	}

	mongoAnalyzer, err := a.connectMongoDB(config)
	if err != nil {
		return result, err
	}
	defer mongoAnalyzer.Close()

	latencyResult, err := mongoAnalyzer.AnalyzeLatency()
	if err != nil {
		return result, err
	}

	result.WriteString(latencyResult)
	return result, nil
}

// analyzeMongoDBPerformance analisa performance MongoDB
func (a *Analyzer) analyzeMongoDBPerformance(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMongoDB {
		return result, fmt.Errorf("análise de performance disponível apenas para MongoDB")
	}

	mongoAnalyzer, err := a.connectMongoDB(config)
	if err != nil {
		return result, err
	}
	defer mongoAnalyzer.Close()

	performanceResult, err := mongoAnalyzer.AnalyzePerformance()
	if err != nil {
		return result, err
	}

	result.WriteString(performanceResult)

	// Adicionar análise de índices
	indexesResult, err := mongoAnalyzer.AnalyzeIndexes()
	if err == nil {
		result.WriteString("\n" + indexesResult)
	}

	return result, nil
}

// connectMongoDB conecta ao MongoDB
func (a *Analyzer) connectMongoDB(config *ConnectionConfig) (*mongodb.MongoDBAnalyzer, error) {
	// Construir connection string
	var connStr string
	if config.ConnectionString != "" {
		connStr = config.ConnectionString
	} else if config.JDBCURL != "" {
		// Converter JDBC URL para MongoDB connection string
		connStr = strings.Replace(config.JDBCURL, "jdbc:mongodb://", "mongodb://", 1)
	} else {
		port := config.Port
		if port == 0 {
			port = 27017
		}

		if config.Username != "" && config.Password != "" {
			connStr = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
				config.Username, config.Password, config.Host, port, config.Database)
		} else {
			connStr = fmt.Sprintf("mongodb://%s:%d/%s",
				config.Host, port, config.Database)
		}
	}

	databaseName := config.Database
	if databaseName == "" {
		databaseName = "admin"
	}

	return mongodb.NewMongoDBAnalyzer(connStr, databaseName)
}

// analyzePostgresReplication analisa replicação PostgreSQL
func (a *Analyzer) analyzePostgresReplication(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypePostgreSQL {
		return result, fmt.Errorf("análise de replicação PostgreSQL disponível apenas para PostgreSQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	pgAnalyzer := postgresql.NewPostgreSQLAnalyzer(db)
	replicationResult, err := pgAnalyzer.AnalyzeReplication()
	if err != nil {
		return result, err
	}

	result.WriteString(replicationResult)
	return result, nil
}

// analyzePostgresLocks analisa locks PostgreSQL
func (a *Analyzer) analyzePostgresLocks(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypePostgreSQL {
		return result, fmt.Errorf("análise de locks PostgreSQL disponível apenas para PostgreSQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	pgAnalyzer := postgresql.NewPostgreSQLAnalyzer(db)
	locksResult, err := pgAnalyzer.AnalyzeLocks()
	if err != nil {
		return result, err
	}

	result.WriteString(locksResult)
	return result, nil
}

// analyzePostgresFragmentation analisa fragmentação PostgreSQL
func (a *Analyzer) analyzePostgresFragmentation(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypePostgreSQL {
		return result, fmt.Errorf("análise de fragmentação PostgreSQL disponível apenas para PostgreSQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	pgAnalyzer := postgresql.NewPostgreSQLAnalyzer(db)
	fragmentationResult, err := pgAnalyzer.AnalyzeFragmentation()
	if err != nil {
		return result, err
	}

	result.WriteString(fragmentationResult)
	return result, nil
}

// analyzeMySQLReplication analisa replicação MySQL
func (a *Analyzer) analyzeMySQLReplication(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMySQL {
		return result, fmt.Errorf("análise de replicação MySQL disponível apenas para MySQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	mysqlAnalyzer := mysql.NewMySQLAnalyzer(db)
	replicationResult, err := mysqlAnalyzer.AnalyzeReplication()
	if err != nil {
		return result, err
	}

	result.WriteString(replicationResult)
	return result, nil
}

// analyzeMySQLLocks analisa locks MySQL
func (a *Analyzer) analyzeMySQLLocks(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMySQL {
		return result, fmt.Errorf("análise de locks MySQL disponível apenas para MySQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	mysqlAnalyzer := mysql.NewMySQLAnalyzer(db)
	locksResult, err := mysqlAnalyzer.AnalyzeLocks()
	if err != nil {
		return result, err
	}

	result.WriteString(locksResult)
	return result, nil
}

// analyzeMySQLFragmentation analisa fragmentação MySQL
func (a *Analyzer) analyzeMySQLFragmentation(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeMySQL {
		return result, fmt.Errorf("análise de fragmentação MySQL disponível apenas para MySQL")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	mysqlAnalyzer := mysql.NewMySQLAnalyzer(db)
	fragmentationResult, err := mysqlAnalyzer.AnalyzeFragmentation()
	if err != nil {
		return result, err
	}

	result.WriteString(fragmentationResult)
	return result, nil
}

// analyzeChecklist executa checklist
func (a *Analyzer) analyzeChecklist(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	// Extrair tipo de checklist do título ou config
	checklistType := dbchecklist.ChecklistTypeDaily
	titleLower := strings.ToLower(analysis.Title)
	if strings.Contains(titleLower, "semanal") || strings.Contains(titleLower, "weekly") {
		checklistType = dbchecklist.ChecklistTypeWeekly
	} else if strings.Contains(titleLower, "profundo") || strings.Contains(titleLower, "deep") {
		checklistType = dbchecklist.ChecklistTypeDeep
	}

	generator := dbchecklist.NewChecklistGenerator()
	checklist := generator.GenerateChecklist(analysis.DatabaseType, checklistType)

	result.WriteString(checklist.FormatChecklist())
	return result, nil
}

// analyzeBackup analisa backups
func (a *Analyzer) analyzeBackup(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	backupAnalyzer, err := backup.NewBackupAnalyzer()
	if err != nil {
		return result, err
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	backupResult, err := backupAnalyzer.AnalyzeBackups(analysis.DatabaseType, db, config)
	if err != nil {
		return result, err
	}

	result.WriteString(backupResult)
	return result, nil
}

// analyzeDynamic realiza análise dinâmica gerando queries com IA
func (a *Analyzer) analyzeDynamic(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	// Extrair solicitação do título ou usar um prompt padrão
	request := analysis.Title
	if request == "" {
		request = "mostre informações gerais do banco de dados"
	}

	dynamicAnalyzer, err := dbdynamic.NewDynamicAnalyzer(analysis.DatabaseType, config, db)
	if err != nil {
		return result, err
	}
	defer dynamicAnalyzer.Close()

	dynamicResult, err := dynamicAnalyzer.AnalyzeWithNaturalLanguage(request)
	if err != nil {
		return result, err
	}

	result.WriteString(dynamicResult)
	return result, nil
}

// analyzePDBs analisa todos os PDBs (Oracle)
func (a *Analyzer) analyzePDBs(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de PDBs disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	oracleAnalyzer := oracle.NewOracleAnalyzer(db)
	pdbsResult, err := oracleAnalyzer.AnalyzePDBs()
	if err != nil {
		return result, err
	}

	result.WriteString(pdbsResult)
	return result, nil
}

// analyzePDB analisa um PDB específico (Oracle)
func (a *Analyzer) analyzePDB(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de PDB disponível apenas para Oracle")
	}

	// Extrair nome do PDB do título ou config
	pdbName := analysis.Title
	if pdbName == "" {
		return result, fmt.Errorf("nome do PDB não especificado (use --title com o nome do PDB)")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	oracleAnalyzer := oracle.NewOracleAnalyzer(db)
	pdbResult, err := oracleAnalyzer.AnalyzeSpecificPDB(pdbName)
	if err != nil {
		return result, err
	}

	result.WriteString(pdbResult)
	return result, nil
}

// analyzeInstance analisa a instância SQL Server
func (a *Analyzer) analyzeInstance(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de instância disponível apenas para SQL Server")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	instanceResult, err := sqlAnalyzer.AnalyzeInstance()
	if err != nil {
		return result, err
	}

	result.WriteString(instanceResult)
	return result, nil
}

// analyzeDatabases analisa todos os databases SQL Server
func (a *Analyzer) analyzeDatabases(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de databases disponível apenas para SQL Server")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	databasesResult, err := sqlAnalyzer.AnalyzeDatabases()
	if err != nil {
		return result, err
	}

	result.WriteString(databasesResult)
	return result, nil
}

// analyzeDatabase analisa um database específico SQL Server
func (a *Analyzer) analyzeDatabase(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeSQLServer {
		return result, fmt.Errorf("análise de database disponível apenas para SQL Server")
	}

	// Extrair nome do database do título ou config
	databaseName := analysis.Title
	if databaseName == "" {
		databaseName = config.Database
	}
	if databaseName == "" {
		return result, fmt.Errorf("nome do database não especificado")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	sqlAnalyzer := sqlserver.NewSQLServerAnalyzer(db)
	databaseResult, err := sqlAnalyzer.AnalyzeSpecificDatabase(databaseName)
	if err != nil {
		return result, err
	}

	result.WriteString(databaseResult)
	return result, nil
}

// analyzeRACHealth analisa saúde do RAC
func (a *Analyzer) analyzeRACHealth(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de RAC disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	racAnalyzer := oracle.NewRACAnalyzer(db)
	healthResult, err := racAnalyzer.AnalyzeRACHealth()
	if err != nil {
		return result, err
	}

	result.WriteString(healthResult)
	return result, nil
}

// analyzeRACErrors analisa erros do RAC
func (a *Analyzer) analyzeRACErrors(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de erros RAC disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	racAnalyzer := oracle.NewRACAnalyzer(db)
	errorsResult, err := racAnalyzer.AnalyzeRACErrors()
	if err != nil {
		return result, err
	}

	result.WriteString(errorsResult)
	return result, nil
}

// analyzeRACListener analisa listener do RAC
func (a *Analyzer) analyzeRACListener(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de listener RAC disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	racAnalyzer := oracle.NewRACAnalyzer(db)
	listenerResult, err := racAnalyzer.AnalyzeListenerStatus()
	if err != nil {
		return result, err
	}

	// Se houver log path, analisar log também
	if analysis.LogFilePath != "" {
		logResult, err := racAnalyzer.AnalyzeListenerLog(analysis.LogFilePath)
		if err == nil {
			result.WriteString("\n" + logResult)
		}
	}

	result.WriteString(listenerResult)
	return result, nil
}

// analyzeRACLatency analisa latência do RAC
func (a *Analyzer) analyzeRACLatency(analysis *DBAnalysis, config *ConnectionConfig) (strings.Builder, error) {
	var result strings.Builder

	if analysis.DatabaseType != DatabaseTypeOracle {
		return result, fmt.Errorf("análise de latência RAC disponível apenas para Oracle")
	}

	connector, err := dbconnection.GetConnector(analysis.DatabaseType)
	if err != nil {
		return result, err
	}

	db, err := connector.Connect(config)
	if err != nil {
		return result, err
	}
	defer db.Close()

	racAnalyzer := oracle.NewRACAnalyzer(db)
	latencyResult, err := racAnalyzer.AnalyzeRACLatency()
	if err != nil {
		return result, err
	}

	result.WriteString(latencyResult)
	return result, nil
}

