package dbanalysis

import (
	"time"

	"github.com/snip/internal/dbtypes"
)

// DatabaseType é um alias para dbtypes.DatabaseType
type DatabaseType = dbtypes.DatabaseType

// Constantes de tipo de banco
const (
	DatabaseTypeOracle    = dbtypes.DatabaseTypeOracle
	DatabaseTypeSQLServer = dbtypes.DatabaseTypeSQLServer
	DatabaseTypeMySQL     = dbtypes.DatabaseTypeMySQL
	DatabaseTypePostgreSQL = dbtypes.DatabaseTypePostgreSQL
	DatabaseTypeMongoDB   = dbtypes.DatabaseTypeMongoDB
)

// AnalysisType representa o tipo de análise a ser realizada
type AnalysisType string

const (
	AnalysisTypeDiagnostic      AnalysisType = "diagnostic"
	AnalysisTypeTuning          AnalysisType = "tuning"
	AnalysisTypeQuery           AnalysisType = "query"
	AnalysisTypeTablespace      AnalysisType = "tablespace"
	AnalysisTypeDisk            AnalysisType = "disk"
	AnalysisTypeTables          AnalysisType = "tables"
	AnalysisTypeIndexes         AnalysisType = "indexes"
	AnalysisTypeLogs            AnalysisType = "logs"
	AnalysisTypePredictive      AnalysisType = "predictive"
	AnalysisTypeErrorKnowledge  AnalysisType = "error_knowledge"
	// Análises específicas Oracle
	AnalysisTypeAWR             AnalysisType = "awr"
	AnalysisTypeASH             AnalysisType = "ash"
	AnalysisTypeExecutionPlan   AnalysisType = "execution_plan"
	// Análises específicas SQL Server
	AnalysisTypeLocks           AnalysisType = "locks"
	AnalysisTypeActiveSessions  AnalysisType = "active_sessions"
	AnalysisTypeRunningQueries  AnalysisType = "running_queries"
	// Análises específicas MongoDB
	AnalysisTypeReplication     AnalysisType = "replication"
	AnalysisTypeSharding        AnalysisType = "sharding"
	AnalysisTypeLatency         AnalysisType = "latency"
	AnalysisTypePerformance     AnalysisType = "performance"
	// Análises específicas PostgreSQL
	AnalysisTypePostgresReplication AnalysisType = "postgres_replication"
	AnalysisTypePostgresLocks       AnalysisType = "postgres_locks"
	AnalysisTypePostgresFragmentation AnalysisType = "postgres_fragmentation"
	// Análises específicas MySQL
	AnalysisTypeMySQLReplication  AnalysisType = "mysql_replication"
	AnalysisTypeMySQLLocks        AnalysisType = "mysql_locks"
	AnalysisTypeMySQLFragmentation AnalysisType = "mysql_fragmentation"
	// Checklist e Backup
	AnalysisTypeChecklist         AnalysisType = "checklist"
	AnalysisTypeBackup            AnalysisType = "backup"
	// Análises dinâmicas e chat
	AnalysisTypeDynamic           AnalysisType = "dynamic"
	AnalysisTypeChat              AnalysisType = "chat"
	// Oracle PDBs
	AnalysisTypePDBs              AnalysisType = "pdbs"
	AnalysisTypePDB                AnalysisType = "pdb"
	// SQL Server Instance e Databases
	AnalysisTypeInstance           AnalysisType = "instance"
	AnalysisTypeDatabases          AnalysisType = "databases"
	AnalysisTypeDatabase           AnalysisType = "database"
	// Oracle RAC
	AnalysisTypeRACHealth          AnalysisType = "rac_health"
	AnalysisTypeRACErrors          AnalysisType = "rac_errors"
	AnalysisTypeRACListener        AnalysisType = "rac_listener"
	AnalysisTypeRACLatency         AnalysisType = "rac_latency"
)

// OutputType representa o formato de saída
type OutputType string

const (
	OutputTypeJSON  OutputType = "json"
	OutputTypeMarkdown OutputType = "markdown"
	OutputTypeText  OutputType = "text"
	OutputTypeHTML  OutputType = "html"
)

// ConnectionConfig é um alias para dbtypes.ConnectionConfig
type ConnectionConfig = dbtypes.ConnectionConfig

// DBAnalysis representa uma análise de banco de dados
type DBAnalysis struct {
	ID              int          `json:"id"`
	Title           string       `json:"title"`
	DatabaseType    DatabaseType `json:"database_type"`
	AnalysisType    AnalysisType `json:"analysis_type"`
	ConnectionConfig string      `json:"connection_config"` // JSON string
	LogFilePath     string       `json:"log_file_path,omitempty"`
	OutputType      OutputType   `json:"output_type"`
	Result          string       `json:"result"` // Resultado da análise
	AIInsights      string       `json:"ai_insights,omitempty"` // Insights gerados pela IA
	Status          string       `json:"status"` // pending, completed, error
	ErrorMessage    string       `json:"error_message,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

// NewDBAnalysis cria uma nova análise de banco de dados
func NewDBAnalysis(title string, dbType DatabaseType, analysisType AnalysisType, outputType OutputType) *DBAnalysis {
	now := time.Now()
	return &DBAnalysis{
		Title:        title,
		DatabaseType: dbType,
		AnalysisType: analysisType,
		OutputType:   outputType,
		Status:       "pending",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

