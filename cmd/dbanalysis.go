package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/dbanalysis"
	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var (
	dbAnalysisTitle      string
	dbAnalysisType       string
	dbAnalysisAnalysisType string
	dbAnalysisOutputType  string
	dbAnalysisHost        string
	dbAnalysisPort        int
	dbAnalysisDatabase    string
	dbAnalysisUsername    string
	dbAnalysisPassword    string
	dbAnalysisIsRemote    bool
	dbAnalysisJDBCURL     string
	dbAnalysisConnString  string
	dbAnalysisLogPath      string
	dbAnalysisLimit        int
	dbAnalysisDBType       string
	dbAnalysisAnalysisTypeFilter string
	dbAnalysisWithChart    bool
	dbAnalysisChartType    string
	dbAnalysisGeneratePlan bool
	dbAnalysisGenerateProject bool
)

func init() {
	// Comando create
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisTitle, "title", "t", "", "Título da análise")
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisType, "db-type", "d", "", "Tipo de banco (oracle, sqlserver, mysql, postgresql, mongodb)")
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisAnalysisType, "analysis-type", "a", "", "Tipo de análise (diagnostic, tuning, query, tablespace, disk, tables, indexes, logs, predictive, error_knowledge, awr, ash, execution_plan, locks, active_sessions, running_queries, replication, sharding, latency, performance, postgres_replication, postgres_locks, postgres_fragmentation, mysql_replication, mysql_locks, mysql_fragmentation, checklist, backup, dynamic, pdbs, pdb, instance, databases, database, rac_health, rac_errors, rac_listener, rac_latency)")
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisOutputType, "output", "o", "markdown", "Tipo de saída (json, markdown, text, html)")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisHost, "host", "localhost", "Host do banco de dados")
	dbAnalysisCreateCmd.Flags().IntVar(&dbAnalysisPort, "port", 0, "Porta do banco de dados")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisDatabase, "database", "", "Nome do banco de dados")
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisUsername, "username", "u", "", "Usuário do banco de dados")
	dbAnalysisCreateCmd.Flags().StringVarP(&dbAnalysisPassword, "password", "p", "", "Senha do banco de dados (opcional para conexões locais)")
	dbAnalysisCreateCmd.Flags().BoolVar(&dbAnalysisIsRemote, "remote", false, "Conexão remota")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisJDBCURL, "jdbc-url", "", "URL JDBC completa")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisConnString, "conn-string", "", "String de conexão completa")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisLogPath, "log-path", "", "Caminho do arquivo de log (.log ou .xml)")
	dbAnalysisCreateCmd.Flags().BoolVar(&dbAnalysisWithChart, "with-chart", false, "Incluir gráfico na análise")
	dbAnalysisCreateCmd.Flags().StringVar(&dbAnalysisChartType, "chart-type", "", "Tipo de gráfico (line, bar, pie, area, table, ascii, html)")
	dbAnalysisCreateCmd.Flags().BoolVar(&dbAnalysisGeneratePlan, "generate-plan", false, "Gerar plano de manutenção após análise")
	dbAnalysisCreateCmd.Flags().BoolVar(&dbAnalysisGenerateProject, "generate-project", false, "Transformar análise em projeto")

	// Comando list
	dbAnalysisListCmd.Flags().IntVarP(&dbAnalysisLimit, "limit", "l", 0, "Limite de resultados")
	dbAnalysisListCmd.Flags().StringVar(&dbAnalysisDBType, "db-type", "", "Filtrar por tipo de banco")
	dbAnalysisListCmd.Flags().StringVar(&dbAnalysisAnalysisTypeFilter, "analysis-type", "", "Filtrar por tipo de análise")

	// Comando get
	dbAnalysisGetCmd.Flags().BoolP("verbose", "v", false, "Mostrar informações detalhadas")

	// Comando run
	dbAnalysisRunCmd.Flags().StringVarP(&dbAnalysisTitle, "title", "t", "", "Título da análise (opcional)")

	// Adicionar comandos ao root
	rootCmd.AddCommand(dbAnalysisCmd)
	dbAnalysisCmd.AddCommand(dbAnalysisCreateCmd)
	dbAnalysisCmd.AddCommand(dbAnalysisListCmd)
	dbAnalysisCmd.AddCommand(dbAnalysisGetCmd)
	dbAnalysisCmd.AddCommand(dbAnalysisDeleteCmd)
	dbAnalysisCmd.AddCommand(dbAnalysisRunCmd)
}

var dbAnalysisCmd = &cobra.Command{
	Use:   "db-analysis",
	Short: "Análise de banco de dados com IA",
	Long: `Comandos para análise de banco de dados usando IA.

Suporta análise de:
- Oracle, SQL Server, MySQL, PostgreSQL, MongoDB
- Logs (.log e .xml)
- Diagnóstico, tuning, consultas, tablespace, disco, tabelas, índices
- Análises preditivas e base de conhecimento de erros`,
}

var dbAnalysisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Criar uma nova análise de banco de dados",
	Long: `Cria uma nova análise de banco de dados.

Exemplos:
  snip db-analysis create --title "Análise PostgreSQL" --db-type postgresql --analysis-type diagnostic --host localhost --port 5432 --database mydb --username user --password pass
  snip db-analysis create --title "Análise de Logs Oracle" --db-type oracle --analysis-type logs --log-path /path/to/logfile.log
  snip db-analysis create --title "Tuning MySQL" --db-type mysql --analysis-type tuning --jdbc-url "jdbc:mysql://localhost:3306/db"
  snip db-analysis create --title "Checklist Diário PostgreSQL" --db-type postgresql --analysis-type checklist --host localhost --port 5432 --database mydb
  snip db-analysis create --title "Análise de Backups" --db-type sqlserver --analysis-type backup --host localhost --port 1433 --database mydb
  snip db-analysis create --title "Locks PostgreSQL" --db-type postgresql --analysis-type postgres_locks --host localhost --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			if dbAnalysisTitle == "" {
				return fmt.Errorf("título é obrigatório (use --title)")
			}
			if dbAnalysisType == "" {
				return fmt.Errorf("tipo de banco é obrigatório (use --db-type)")
			}
			if dbAnalysisAnalysisType == "" {
				return fmt.Errorf("tipo de análise é obrigatório (use --analysis-type)")
			}

			dbType := parseDatabaseType(dbAnalysisType)
			analysisType := parseAnalysisType(dbAnalysisAnalysisType)
			outputType := parseOutputType(dbAnalysisOutputType)

			config := &dbanalysis.ConnectionConfig{
				Type:            dbType,
				Host:            dbAnalysisHost,
				Port:            dbAnalysisPort,
				Database:        dbAnalysisDatabase,
				Username:        dbAnalysisUsername,
				Password:        dbAnalysisPassword,
				IsRemote:        dbAnalysisIsRemote,
				JDBCURL:         dbAnalysisJDBCURL,
				ConnectionString: dbAnalysisConnString,
			}

			return h.CreateAnalysis(dbAnalysisTitle, dbType, analysisType, outputType, config, dbAnalysisLogPath)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var dbAnalysisListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar análises de banco de dados",
	Long: `Lista todas as análises de banco de dados armazenadas.

Exemplos:
  snip db-analysis list
  snip db-analysis list --limit 10
  snip db-analysis list --db-type postgresql
  snip db-analysis list --analysis-type diagnostic`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			var dbType dbanalysis.DatabaseType
			var analysisType dbanalysis.AnalysisType

			if dbAnalysisDBType != "" {
				dbType = parseDatabaseType(dbAnalysisDBType)
			}
			if dbAnalysisAnalysisTypeFilter != "" {
				analysisType = parseAnalysisType(dbAnalysisAnalysisTypeFilter)
			}

			return h.ListAnalyses(dbAnalysisLimit, dbType, analysisType)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var dbAnalysisGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Obter detalhes de uma análise",
	Long: `Obtém os detalhes completos de uma análise de banco de dados.

Exemplos:
  snip db-analysis get 1
  snip db-analysis get 1 --verbose`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			return h.GetAnalysis(args[0], verbose)
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var dbAnalysisDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Deletar uma análise",
	Long: `Deleta uma análise de banco de dados.

Exemplo:
  snip db-analysis delete 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			return h.DeleteAnalysis(args[0])
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

var dbAnalysisRunCmd = &cobra.Command{
	Use:   "run [id]",
	Short: "Executar uma análise",
	Long: `Executa uma análise de banco de dados.

A análise será executada e os resultados serão armazenados para consulta posterior.

Exemplo:
  snip db-analysis run 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			return h.RunAnalysis(args[0])
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

// Funções auxiliares para parsing
func parseDatabaseType(s string) dbanalysis.DatabaseType {
	s = strings.ToLower(s)
	switch s {
	case "oracle":
		return dbanalysis.DatabaseTypeOracle
	case "sqlserver", "sql-server", "mssql":
		return dbanalysis.DatabaseTypeSQLServer
	case "mysql":
		return dbanalysis.DatabaseTypeMySQL
	case "postgresql", "postgres":
		return dbanalysis.DatabaseTypePostgreSQL
	case "mongodb", "mongo":
		return dbanalysis.DatabaseTypeMongoDB
	default:
		return dbanalysis.DatabaseType(s)
	}
}

func parseAnalysisType(s string) dbanalysis.AnalysisType {
	s = strings.ToLower(s)
	switch s {
	case "diagnostic":
		return dbanalysis.AnalysisTypeDiagnostic
	case "tuning":
		return dbanalysis.AnalysisTypeTuning
	case "query":
		return dbanalysis.AnalysisTypeQuery
	case "tablespace":
		return dbanalysis.AnalysisTypeTablespace
	case "disk":
		return dbanalysis.AnalysisTypeDisk
	case "tables":
		return dbanalysis.AnalysisTypeTables
	case "indexes", "index":
		return dbanalysis.AnalysisTypeIndexes
	case "logs", "log":
		return dbanalysis.AnalysisTypeLogs
	case "predictive":
		return dbanalysis.AnalysisTypePredictive
	case "error_knowledge", "error-knowledge", "errors":
		return dbanalysis.AnalysisTypeErrorKnowledge
	// Oracle
	case "awr":
		return dbanalysis.AnalysisTypeAWR
	case "ash":
		return dbanalysis.AnalysisTypeASH
	case "execution_plan", "execution-plan", "plan":
		return dbanalysis.AnalysisTypeExecutionPlan
	// SQL Server
	case "locks":
		return dbanalysis.AnalysisTypeLocks
	case "active_sessions", "active-sessions", "sessions":
		return dbanalysis.AnalysisTypeActiveSessions
	case "running_queries", "running-queries":
		return dbanalysis.AnalysisTypeRunningQueries
	// MongoDB
	case "replication":
		return dbanalysis.AnalysisTypeReplication
	case "sharding":
		return dbanalysis.AnalysisTypeSharding
	case "latency":
		return dbanalysis.AnalysisTypeLatency
	case "performance":
		return dbanalysis.AnalysisTypePerformance
	// PostgreSQL
	case "postgres_replication", "postgres-replication", "pg_replication":
		return dbanalysis.AnalysisTypePostgresReplication
	case "postgres_locks", "postgres-locks", "pg_locks":
		return dbanalysis.AnalysisTypePostgresLocks
	case "postgres_fragmentation", "postgres-fragmentation", "pg_fragmentation":
		return dbanalysis.AnalysisTypePostgresFragmentation
	// MySQL
	case "mysql_replication", "mysql-replication":
		return dbanalysis.AnalysisTypeMySQLReplication
	case "mysql_locks", "mysql-locks":
		return dbanalysis.AnalysisTypeMySQLLocks
	case "mysql_fragmentation", "mysql-fragmentation":
		return dbanalysis.AnalysisTypeMySQLFragmentation
	// Checklist e Backup
	case "checklist":
		return dbanalysis.AnalysisTypeChecklist
	case "backup", "backups":
		return dbanalysis.AnalysisTypeBackup
	// Análises dinâmicas
	case "dynamic", "dinamica":
		return dbanalysis.AnalysisTypeDynamic
	// Oracle PDBs
	case "pdbs", "pdb_list":
		return dbanalysis.AnalysisTypePDBs
	case "pdb":
		return dbanalysis.AnalysisTypePDB
	// SQL Server Instance e Databases
	case "instance", "instancia":
		return dbanalysis.AnalysisTypeInstance
	case "databases", "databases_list":
		return dbanalysis.AnalysisTypeDatabases
	case "database", "db":
		return dbanalysis.AnalysisTypeDatabase
	// Oracle RAC
	case "rac_health", "rac-health":
		return dbanalysis.AnalysisTypeRACHealth
	case "rac_errors", "rac-errors":
		return dbanalysis.AnalysisTypeRACErrors
	case "rac_listener", "rac-listener":
		return dbanalysis.AnalysisTypeRACListener
	case "rac_latency", "rac-latency":
		return dbanalysis.AnalysisTypeRACLatency
	default:
		return dbanalysis.AnalysisType(s)
	}
}

func parseOutputType(s string) dbanalysis.OutputType {
	s = strings.ToLower(s)
	switch s {
	case "json":
		return dbanalysis.OutputTypeJSON
	case "markdown", "md":
		return dbanalysis.OutputTypeMarkdown
	case "text", "txt":
		return dbanalysis.OutputTypeText
	case "html":
		return dbanalysis.OutputTypeHTML
	default:
		return dbanalysis.OutputTypeMarkdown
	}
}

