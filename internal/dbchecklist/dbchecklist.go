package dbchecklist

import (
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/dbtypes"
)

// ChecklistType representa o tipo de checklist
type ChecklistType string

const (
	ChecklistTypeDaily    ChecklistType = "daily"
	ChecklistTypeWeekly   ChecklistType = "weekly"
	ChecklistTypeDeep     ChecklistType = "deep"
)

// ChecklistItem representa um item de checklist
type ChecklistItem struct {
	ID          string
	Title       string
	Description string
	Category    string
	Priority    string // high, medium, low
	Status      string // pending, completed, skipped, failed
	Result      string
	CheckedAt   time.Time
}

// Checklist representa um checklist completo
type Checklist struct {
	ID           int
	DatabaseType dbtypes.DatabaseType
	ChecklistType ChecklistType
	Items        []ChecklistItem
	CreatedAt    time.Time
	CompletedAt  time.Time
	Status       string // pending, in_progress, completed
}

// ChecklistGenerator gera checklists para diferentes bancos
type ChecklistGenerator struct{}

// NewChecklistGenerator cria um novo gerador de checklist
func NewChecklistGenerator() *ChecklistGenerator {
	return &ChecklistGenerator{}
}

// GenerateChecklist gera um checklist baseado no tipo de banco e nível
func (g *ChecklistGenerator) GenerateChecklist(dbType dbtypes.DatabaseType, checklistType ChecklistType) *Checklist {
	checklist := &Checklist{
		DatabaseType:  dbType,
		ChecklistType: checklistType,
		Items:         []ChecklistItem{},
		CreatedAt:     time.Now(),
		Status:        "pending",
	}

	switch dbType {
	case dbtypes.DatabaseTypeOracle:
		checklist.Items = g.generateOracleChecklist(checklistType)
	case dbtypes.DatabaseTypeSQLServer:
		checklist.Items = g.generateSQLServerChecklist(checklistType)
	case dbtypes.DatabaseTypeMySQL:
		checklist.Items = g.generateMySQLChecklist(checklistType)
	case dbtypes.DatabaseTypePostgreSQL:
		checklist.Items = g.generatePostgreSQLChecklist(checklistType)
	case dbtypes.DatabaseTypeMongoDB:
		checklist.Items = g.generateMongoDBChecklist(checklistType)
	}

	return checklist
}

// generateOracleChecklist gera checklist para Oracle
func (g *ChecklistGenerator) generateOracleChecklist(checklistType ChecklistType) []ChecklistItem {
	items := []ChecklistItem{}

	// Itens diários
	dailyItems := []ChecklistItem{
		{ID: "ora_daily_1", Title: "Verificar alertas do banco", Category: "Monitoring", Priority: "high"},
		{ID: "ora_daily_2", Title: "Verificar espaço em tablespaces", Category: "Storage", Priority: "high"},
		{ID: "ora_daily_3", Title: "Verificar processos ativos", Category: "Performance", Priority: "medium"},
		{ID: "ora_daily_4", Title: "Verificar logs de erro", Category: "Logs", Priority: "high"},
		{ID: "ora_daily_5", Title: "Verificar status de backup", Category: "Backup", Priority: "high"},
	}

	// Itens semanais
	weeklyItems := []ChecklistItem{
		{ID: "ora_weekly_1", Title: "Analisar AWR reports", Category: "Performance", Priority: "high"},
		{ID: "ora_weekly_2", Title: "Verificar fragmentação de tabelas", Category: "Maintenance", Priority: "medium"},
		{ID: "ora_weekly_3", Title: "Revisar estatísticas de objetos", Category: "Optimization", Priority: "medium"},
		{ID: "ora_weekly_4", Title: "Verificar índices não utilizados", Category: "Optimization", Priority: "low"},
		{ID: "ora_weekly_5", Title: "Analisar ASH para queries lentas", Category: "Performance", Priority: "high"},
		{ID: "ora_weekly_6", Title: "Verificar integridade de dados", Category: "Integrity", Priority: "high"},
		{ID: "ora_weekly_7", Title: "Revisar configurações de segurança", Category: "Security", Priority: "medium"},
	}

	// Itens profundos
	deepItems := []ChecklistItem{
		{ID: "ora_deep_1", Title: "Análise completa de performance (AWR)", Category: "Performance", Priority: "high"},
		{ID: "ora_deep_2", Title: "Auditoria de segurança completa", Category: "Security", Priority: "high"},
		{ID: "ora_deep_3", Title: "Análise de capacidade e crescimento", Category: "Capacity", Priority: "high"},
		{ID: "ora_deep_4", Title: "Revisão completa de índices", Category: "Optimization", Priority: "high"},
		{ID: "ora_deep_5", Title: "Análise de fragmentação completa", Category: "Maintenance", Priority: "medium"},
		{ID: "ora_deep_6", Title: "Revisão de parâmetros de instância", Category: "Configuration", Priority: "high"},
		{ID: "ora_deep_7", Title: "Teste de restore de backup", Category: "Backup", Priority: "high"},
		{ID: "ora_deep_8", Title: "Análise de replicação/DataGuard", Category: "High Availability", Priority: "high"},
		{ID: "ora_deep_9", Title: "Revisão de permissões e roles", Category: "Security", Priority: "high"},
		{ID: "ora_deep_10", Title: "Análise de latência de rede", Category: "Network", Priority: "medium"},
	}

	switch checklistType {
	case ChecklistTypeDaily:
		items = dailyItems
	case ChecklistTypeWeekly:
		items = append(dailyItems, weeklyItems...)
	case ChecklistTypeDeep:
		items = append(append(dailyItems, weeklyItems...), deepItems...)
	}

	return items
}

// generateSQLServerChecklist gera checklist para SQL Server
func (g *ChecklistGenerator) generateSQLServerChecklist(checklistType ChecklistType) []ChecklistItem {
	items := []ChecklistItem{}

	dailyItems := []ChecklistItem{
		{ID: "mssql_daily_1", Title: "Verificar jobs do SQL Agent", Category: "Jobs", Priority: "high"},
		{ID: "mssql_daily_2", Title: "Verificar espaço em arquivos de dados", Category: "Storage", Priority: "high"},
		{ID: "mssql_daily_3", Title: "Verificar locks e bloqueios", Category: "Performance", Priority: "high"},
		{ID: "mssql_daily_4", Title: "Verificar logs de erro", Category: "Logs", Priority: "high"},
		{ID: "mssql_daily_5", Title: "Verificar status de backup", Category: "Backup", Priority: "high"},
	}

	weeklyItems := []ChecklistItem{
		{ID: "mssql_weekly_1", Title: "Analisar DMVs de performance", Category: "Performance", Priority: "high"},
		{ID: "mssql_weekly_2", Title: "Verificar fragmentação de índices", Category: "Maintenance", Priority: "medium"},
		{ID: "mssql_weekly_3", Title: "Revisar estatísticas de tabelas", Category: "Optimization", Priority: "medium"},
		{ID: "mssql_weekly_4", Title: "Verificar índices não utilizados", Category: "Optimization", Priority: "low"},
		{ID: "mssql_weekly_5", Title: "Analisar queries lentas", Category: "Performance", Priority: "high"},
		{ID: "mssql_weekly_6", Title: "Verificar integridade de dados (DBCC)", Category: "Integrity", Priority: "high"},
		{ID: "mssql_weekly_7", Title: "Revisar configurações de segurança", Category: "Security", Priority: "medium"},
	}

	deepItems := []ChecklistItem{
		{ID: "mssql_deep_1", Title: "Análise completa de performance (DMVs)", Category: "Performance", Priority: "high"},
		{ID: "mssql_deep_2", Title: "Auditoria de segurança completa", Category: "Security", Priority: "high"},
		{ID: "mssql_deep_3", Title: "Análise de capacidade e crescimento", Category: "Capacity", Priority: "high"},
		{ID: "mssql_deep_4", Title: "Revisão completa de índices", Category: "Optimization", Priority: "high"},
		{ID: "mssql_deep_5", Title: "Análise de fragmentação completa", Category: "Maintenance", Priority: "medium"},
		{ID: "mssql_deep_6", Title: "Revisão de configurações do servidor", Category: "Configuration", Priority: "high"},
		{ID: "mssql_deep_7", Title: "Teste de restore de backup", Category: "Backup", Priority: "high"},
		{ID: "mssql_deep_8", Title: "Análise de AlwaysOn/Replicação", Category: "High Availability", Priority: "high"},
		{ID: "mssql_deep_9", Title: "Revisão de permissões e roles", Category: "Security", Priority: "high"},
		{ID: "mssql_deep_10", Title: "Análise de latência de I/O", Category: "Performance", Priority: "medium"},
	}

	switch checklistType {
	case ChecklistTypeDaily:
		items = dailyItems
	case ChecklistTypeWeekly:
		items = append(dailyItems, weeklyItems...)
	case ChecklistTypeDeep:
		items = append(append(dailyItems, weeklyItems...), deepItems...)
	}

	return items
}

// generateMySQLChecklist gera checklist para MySQL
func (g *ChecklistGenerator) generateMySQLChecklist(checklistType ChecklistType) []ChecklistItem {
	items := []ChecklistItem{}

	dailyItems := []ChecklistItem{
		{ID: "mysql_daily_1", Title: "Verificar processos lentos", Category: "Performance", Priority: "high"},
		{ID: "mysql_daily_2", Title: "Verificar espaço em disco", Category: "Storage", Priority: "high"},
		{ID: "mysql_daily_3", Title: "Verificar locks e bloqueios", Category: "Performance", Priority: "high"},
		{ID: "mysql_daily_4", Title: "Verificar logs de erro", Category: "Logs", Priority: "high"},
		{ID: "mysql_daily_5", Title: "Verificar status de backup", Category: "Backup", Priority: "high"},
		{ID: "mysql_daily_6", Title: "Verificar status de replicação", Category: "Replication", Priority: "high"},
	}

	weeklyItems := []ChecklistItem{
		{ID: "mysql_weekly_1", Title: "Analisar slow query log", Category: "Performance", Priority: "high"},
		{ID: "mysql_weekly_2", Title: "Verificar fragmentação de tabelas", Category: "Maintenance", Priority: "medium"},
		{ID: "mysql_weekly_3", Title: "Revisar estatísticas de tabelas", Category: "Optimization", Priority: "medium"},
		{ID: "mysql_weekly_4", Title: "Verificar índices não utilizados", Category: "Optimization", Priority: "low"},
		{ID: "mysql_weekly_5", Title: "Analisar queries lentas", Category: "Performance", Priority: "high"},
		{ID: "mysql_weekly_6", Title: "Verificar integridade de dados", Category: "Integrity", Priority: "high"},
		{ID: "mysql_weekly_7", Title: "Revisar configurações de segurança", Category: "Security", Priority: "medium"},
		{ID: "mysql_weekly_8", Title: "Analisar lag de replicação", Category: "Replication", Priority: "high"},
	}

	deepItems := []ChecklistItem{
		{ID: "mysql_deep_1", Title: "Análise completa de performance", Category: "Performance", Priority: "high"},
		{ID: "mysql_deep_2", Title: "Auditoria de segurança completa", Category: "Security", Priority: "high"},
		{ID: "mysql_deep_3", Title: "Análise de capacidade e crescimento", Category: "Capacity", Priority: "high"},
		{ID: "mysql_deep_4", Title: "Revisão completa de índices", Category: "Optimization", Priority: "high"},
		{ID: "mysql_deep_5", Title: "Análise de fragmentação completa", Category: "Maintenance", Priority: "medium"},
		{ID: "mysql_deep_6", Title: "Revisão de configurações do servidor", Category: "Configuration", Priority: "high"},
		{ID: "mysql_deep_7", Title: "Teste de restore de backup", Category: "Backup", Priority: "high"},
		{ID: "mysql_deep_8", Title: "Análise completa de replicação", Category: "Replication", Priority: "high"},
		{ID: "mysql_deep_9", Title: "Revisão de permissões e usuários", Category: "Security", Priority: "high"},
		{ID: "mysql_deep_10", Title: "Análise de binlog e logs", Category: "Logs", Priority: "medium"},
	}

	switch checklistType {
	case ChecklistTypeDaily:
		items = dailyItems
	case ChecklistTypeWeekly:
		items = append(dailyItems, weeklyItems...)
	case ChecklistTypeDeep:
		items = append(append(dailyItems, weeklyItems...), deepItems...)
	}

	return items
}

// generatePostgreSQLChecklist gera checklist para PostgreSQL
func (g *ChecklistGenerator) generatePostgreSQLChecklist(checklistType ChecklistType) []ChecklistItem {
	items := []ChecklistItem{}

	dailyItems := []ChecklistItem{
		{ID: "pg_daily_1", Title: "Verificar processos ativos", Category: "Performance", Priority: "high"},
		{ID: "pg_daily_2", Title: "Verificar espaço em disco", Category: "Storage", Priority: "high"},
		{ID: "pg_daily_3", Title: "Verificar locks e bloqueios", Category: "Performance", Priority: "high"},
		{ID: "pg_daily_4", Title: "Verificar logs de erro", Category: "Logs", Priority: "high"},
		{ID: "pg_daily_5", Title: "Verificar status de backup", Category: "Backup", Priority: "high"},
		{ID: "pg_daily_6", Title: "Verificar status de replicação", Category: "Replication", Priority: "high"},
	}

	weeklyItems := []ChecklistItem{
		{ID: "pg_weekly_1", Title: "Analisar pg_stat_statements", Category: "Performance", Priority: "high"},
		{ID: "pg_weekly_2", Title: "Verificar fragmentação (VACUUM)", Category: "Maintenance", Priority: "medium"},
		{ID: "pg_weekly_3", Title: "Revisar estatísticas (ANALYZE)", Category: "Optimization", Priority: "medium"},
		{ID: "pg_weekly_4", Title: "Verificar índices não utilizados", Category: "Optimization", Priority: "low"},
		{ID: "pg_weekly_5", Title: "Analisar queries lentas", Category: "Performance", Priority: "high"},
		{ID: "pg_weekly_6", Title: "Verificar integridade de dados", Category: "Integrity", Priority: "high"},
		{ID: "pg_weekly_7", Title: "Revisar configurações de segurança", Category: "Security", Priority: "medium"},
		{ID: "pg_weekly_8", Title: "Analisar lag de replicação", Category: "Replication", Priority: "high"},
	}

	deepItems := []ChecklistItem{
		{ID: "pg_deep_1", Title: "Análise completa de performance", Category: "Performance", Priority: "high"},
		{ID: "pg_deep_2", Title: "Auditoria de segurança completa", Category: "Security", Priority: "high"},
		{ID: "pg_deep_3", Title: "Análise de capacidade e crescimento", Category: "Capacity", Priority: "high"},
		{ID: "pg_deep_4", Title: "Revisão completa de índices", Category: "Optimization", Priority: "high"},
		{ID: "pg_deep_5", Title: "Análise de fragmentação completa", Category: "Maintenance", Priority: "medium"},
		{ID: "pg_deep_6", Title: "Revisão de configurações (postgresql.conf)", Category: "Configuration", Priority: "high"},
		{ID: "pg_deep_7", Title: "Teste de restore de backup", Category: "Backup", Priority: "high"},
		{ID: "pg_deep_8", Title: "Análise completa de replicação", Category: "Replication", Priority: "high"},
		{ID: "pg_deep_9", Title: "Revisão de permissões e roles", Category: "Security", Priority: "high"},
		{ID: "pg_deep_10", Title: "Análise de WAL e logs", Category: "Logs", Priority: "medium"},
	}

	switch checklistType {
	case ChecklistTypeDaily:
		items = dailyItems
	case ChecklistTypeWeekly:
		items = append(dailyItems, weeklyItems...)
	case ChecklistTypeDeep:
		items = append(append(dailyItems, weeklyItems...), deepItems...)
	}

	return items
}

// generateMongoDBChecklist gera checklist para MongoDB
func (g *ChecklistGenerator) generateMongoDBChecklist(checklistType ChecklistType) []ChecklistItem {
	items := []ChecklistItem{}

	dailyItems := []ChecklistItem{
		{ID: "mongo_daily_1", Title: "Verificar status de replicação", Category: "Replication", Priority: "high"},
		{ID: "mongo_daily_2", Title: "Verificar espaço em disco", Category: "Storage", Priority: "high"},
		{ID: "mongo_daily_3", Title: "Verificar conexões ativas", Category: "Performance", Priority: "medium"},
		{ID: "mongo_daily_4", Title: "Verificar logs de erro", Category: "Logs", Priority: "high"},
		{ID: "mongo_daily_5", Title: "Verificar status de backup", Category: "Backup", Priority: "high"},
		{ID: "mongo_daily_6", Title: "Verificar status de sharding", Category: "Sharding", Priority: "high"},
	}

	weeklyItems := []ChecklistItem{
		{ID: "mongo_weekly_1", Title: "Analisar queries lentas", Category: "Performance", Priority: "high"},
		{ID: "mongo_weekly_2", Title: "Verificar índices não utilizados", Category: "Optimization", Priority: "low"},
		{ID: "mongo_weekly_3", Title: "Revisar estatísticas de coleções", Category: "Optimization", Priority: "medium"},
		{ID: "mongo_weekly_4", Title: "Analisar latência de operações", Category: "Performance", Priority: "high"},
		{ID: "mongo_weekly_5", Title: "Verificar integridade de dados", Category: "Integrity", Priority: "high"},
		{ID: "mongo_weekly_6", Title: "Revisar configurações de segurança", Category: "Security", Priority: "medium"},
		{ID: "mongo_weekly_7", Title: "Analisar lag de replicação", Category: "Replication", Priority: "high"},
		{ID: "mongo_weekly_8", Title: "Verificar distribuição de chunks", Category: "Sharding", Priority: "medium"},
	}

	deepItems := []ChecklistItem{
		{ID: "mongo_deep_1", Title: "Análise completa de performance", Category: "Performance", Priority: "high"},
		{ID: "mongo_deep_2", Title: "Auditoria de segurança completa", Category: "Security", Priority: "high"},
		{ID: "mongo_deep_3", Title: "Análise de capacidade e crescimento", Category: "Capacity", Priority: "high"},
		{ID: "mongo_deep_4", Title: "Revisão completa de índices", Category: "Optimization", Priority: "high"},
		{ID: "mongo_deep_5", Title: "Análise completa de sharding", Category: "Sharding", Priority: "high"},
		{ID: "mongo_deep_6", Title: "Revisão de configurações do servidor", Category: "Configuration", Priority: "high"},
		{ID: "mongo_deep_7", Title: "Teste de restore de backup", Category: "Backup", Priority: "high"},
		{ID: "mongo_deep_8", Title: "Análise completa de replicação", Category: "Replication", Priority: "high"},
		{ID: "mongo_deep_9", Title: "Revisão de permissões e roles", Category: "Security", Priority: "high"},
		{ID: "mongo_deep_10", Title: "Análise de Oplog e logs", Category: "Logs", Priority: "medium"},
	}

	switch checklistType {
	case ChecklistTypeDaily:
		items = dailyItems
	case ChecklistTypeWeekly:
		items = append(dailyItems, weeklyItems...)
	case ChecklistTypeDeep:
		items = append(append(dailyItems, weeklyItems...), deepItems...)
	}

	return items
}

// FormatChecklist formata checklist para exibição
func (c *Checklist) FormatChecklist() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Checklist %s - %s\n\n", c.ChecklistType, c.DatabaseType))
	result.WriteString(fmt.Sprintf("**Criado em:** %s\n", c.CreatedAt.Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("**Status:** %s\n\n", c.Status))

	// Agrupar por categoria
	categories := make(map[string][]ChecklistItem)
	for _, item := range c.Items {
		categories[item.Category] = append(categories[item.Category], item)
	}

	for category, items := range categories {
		result.WriteString(fmt.Sprintf("## %s\n\n", category))
		result.WriteString("| Item | Prioridade | Status |\n")
		result.WriteString("|------|------------|--------|\n")

		for _, item := range items {
			result.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				item.Title, item.Priority, item.Status))
		}
		result.WriteString("\n")
	}

	return result.String()
}
