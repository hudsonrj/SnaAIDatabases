package backup

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/ai"
	"github.com/snip/internal/dbtypes"
)

// BackupAnalyzer analisa backups de diferentes bancos
type BackupAnalyzer struct {
	aiClient ai.AIClient
}

// NewBackupAnalyzer cria um novo analisador de backups
func NewBackupAnalyzer() (*BackupAnalyzer, error) {
	aiClient, err := ai.NewAIClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &BackupAnalyzer{aiClient: aiClient}, nil
}

// BackupInfo representa informações sobre um backup
type BackupInfo struct {
	DatabaseType    dbtypes.DatabaseType
	BackupType      string // full, incremental, differential, etc
	BackupDate      time.Time
	BackupSize      int64
	BackupDuration  time.Duration
	BackupLocation  string
	Status          string // success, failed, in_progress
	BackupMethod    string // RMAN, mysqldump, pg_dump, etc
	RetentionDays   int
	Compressed      bool
	Encrypted       bool
}

// AnalyzeBackups analisa backups de um banco de dados
func (b *BackupAnalyzer) AnalyzeBackups(dbType dbtypes.DatabaseType, db *sql.DB, config *dbtypes.ConnectionConfig) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# Análise de Backups - %s\n\n", dbType))

	var backups []BackupInfo
	var err error

	switch dbType {
	case dbtypes.DatabaseTypeOracle:
		backups, err = b.analyzeOracleBackups(db)
	case dbtypes.DatabaseTypeSQLServer:
		backups, err = b.analyzeSQLServerBackups(db)
	case dbtypes.DatabaseTypeMySQL:
		backups, err = b.analyzeMySQLBackups(db)
	case dbtypes.DatabaseTypePostgreSQL:
		backups, err = b.analyzePostgreSQLBackups(db)
	case dbtypes.DatabaseTypeMongoDB:
		backups, err = b.analyzeMongoDBBackups(config)
	default:
		return "", fmt.Errorf("tipo de banco não suportado para análise de backups")
	}

	if err != nil {
		return "", err
	}

	if len(backups) == 0 {
		result.WriteString("⚠️ Nenhum backup encontrado ou não foi possível obter informações.\n")
		return result.String(), nil
	}

	// Exibir informações dos backups
	result.WriteString("## Backups Encontrados\n\n")
	result.WriteString("| Data | Tipo | Tamanho | Duração | Status | Localização |\n")
	result.WriteString("|------|------|---------|---------|--------|-------------|\n")

	for _, backup := range backups {
		sizeStr := formatSize(backup.BackupSize)
		durationStr := backup.BackupDuration.String()
		if backup.BackupDuration == 0 {
			durationStr = "N/A"
		}

		result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			backup.BackupDate.Format("2006-01-02 15:04:05"),
			backup.BackupType,
			sizeStr,
			durationStr,
			backup.Status,
			backup.BackupLocation))
	}

	// Análise com IA
	result.WriteString("\n## Análise e Recomendações (IA)\n\n")
	aiAnalysis, err := b.generateAIAnalysis(backups, dbType)
	if err == nil {
		result.WriteString(aiAnalysis)
	}

	return result.String(), nil
}

// analyzeOracleBackups analisa backups Oracle
func (b *BackupAnalyzer) analyzeOracleBackups(db *sql.DB) ([]BackupInfo, error) {
	var backups []BackupInfo

	// Consultar RMAN backups
	query := `
		SELECT 
			recid,
			stamp,
			start_time,
			completion_time,
			elapsed_seconds,
			status,
			backup_type,
			controlfile_included,
			output_device_type
		FROM v$rman_backup_job_details
		ORDER BY start_time DESC
		FETCH FIRST 20 ROWS ONLY
	`

	rows, err := db.Query(query)
	if err != nil {
		// Pode não ter permissões ou RMAN não configurado
		return backups, nil
	}
	defer rows.Close()

	for rows.Next() {
		var recid, stamp, elapsedSeconds int64
		var startTime, completionTime time.Time
		var status, backupType, controlfileIncluded, outputDeviceType sql.NullString

		err := rows.Scan(&recid, &stamp, &startTime, &completionTime, &elapsedSeconds,
			&status, &backupType, &controlfileIncluded, &outputDeviceType)
		if err != nil {
			continue
		}

		backup := BackupInfo{
			DatabaseType:   dbtypes.DatabaseTypeOracle,
			BackupType:     getString(backupType),
			BackupDate:     startTime,
			BackupDuration: time.Duration(elapsedSeconds) * time.Second,
			Status:         getString(status),
			BackupMethod:   "RMAN",
			BackupLocation: getString(outputDeviceType),
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// analyzeSQLServerBackups analisa backups SQL Server
func (b *BackupAnalyzer) analyzeSQLServerBackups(db *sql.DB) ([]BackupInfo, error) {
	var backups []BackupInfo

	query := `
		SELECT 
			database_name,
			backup_start_date,
			backup_finish_date,
			type,
			backup_size,
			physical_device_name,
			compressed,
			encrypted
		FROM msdb.dbo.backupset bs
		INNER JOIN msdb.dbo.backupmediafamily bmf ON bs.media_set_id = bmf.media_set_id
		ORDER BY backup_start_date DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return backups, nil
	}
	defer rows.Close()

	for rows.Next() {
		var databaseName, backupType, physicalDeviceName sql.NullString
		var backupStartDate, backupFinishDate sql.NullTime
		var backupSize sql.NullInt64
		var compressed, encrypted bool

		err := rows.Scan(&databaseName, &backupStartDate, &backupFinishDate, &backupType,
			&backupSize, &physicalDeviceName, &compressed, &encrypted)
		if err != nil {
			continue
		}

		var duration time.Duration
		if backupStartDate.Valid && backupFinishDate.Valid {
			duration = backupFinishDate.Time.Sub(backupStartDate.Time)
		}

		backup := BackupInfo{
			DatabaseType:   dbtypes.DatabaseTypeSQLServer,
			BackupType:     getString(backupType),
			BackupDate:     getTime(backupStartDate),
			BackupSize:     getInt64(backupSize),
			BackupDuration:  duration,
			Status:         "success",
			BackupMethod:   "SQL Server Backup",
			BackupLocation: getString(physicalDeviceName),
			Compressed:     compressed,
			Encrypted:      encrypted,
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// analyzeMySQLBackups analisa backups MySQL
func (b *BackupAnalyzer) analyzeMySQLBackups(db *sql.DB) ([]BackupInfo, error) {
	var backups []BackupInfo

	// MySQL não tem tabela nativa de backups, verificar logs ou arquivos
	// Por enquanto, retornar vazio - seria necessário verificar sistema de arquivos
	// ou logs de backup externos

	return backups, nil
}

// analyzePostgreSQLBackups analisa backups PostgreSQL
func (b *BackupAnalyzer) analyzePostgreSQLBackups(db *sql.DB) ([]BackupInfo, error) {
	var backups []BackupInfo

	// PostgreSQL não tem tabela nativa de backups
	// Verificar pg_backup_history ou arquivos WAL
	// Por enquanto, retornar vazio

	return backups, nil
}

// analyzeMongoDBBackups analisa backups MongoDB
func (b *BackupAnalyzer) analyzeMongoDBBackups(config *dbtypes.ConnectionConfig) ([]BackupInfo, error) {
	var backups []BackupInfo

	// MongoDB backups são geralmente feitos via mongodump ou ferramentas externas
	// Seria necessário verificar logs ou sistema de arquivos

	return backups, nil
}

// generateAIAnalysis gera análise de backups usando IA
func (b *BackupAnalyzer) generateAIAnalysis(backups []BackupInfo, dbType dbtypes.DatabaseType) (string, error) {
	if len(backups) == 0 {
		return "Nenhum backup encontrado para análise.", nil
	}

	// Preparar informações para IA
	var backupSummary strings.Builder
	backupSummary.WriteString(fmt.Sprintf("Tipo de Banco: %s\n\n", dbType))
	backupSummary.WriteString("Backups encontrados:\n")

	latestBackup := backups[0]
	for _, backup := range backups {
		if backup.BackupDate.After(latestBackup.BackupDate) {
			latestBackup = backup
		}
		backupSummary.WriteString(fmt.Sprintf("- Data: %s, Tipo: %s, Status: %s, Tamanho: %s\n",
			backup.BackupDate.Format("2006-01-02 15:04:05"),
			backup.BackupType,
			backup.Status,
			formatSize(backup.BackupSize)))
	}

	backupSummary.WriteString(fmt.Sprintf("\nÚltimo backup: %s\n", latestBackup.BackupDate.Format("2006-01-02 15:04:05")))

	now := time.Now()
	hoursSinceBackup := now.Sub(latestBackup.BackupDate).Hours()

	prompt := fmt.Sprintf(`Você é um DBA experiente. Analise as seguintes informações de backup e forneça:

1. Avaliação do status dos backups
2. Tempo estimado de retorno (RTO) baseado nos backups disponíveis
3. Recomendações de melhoria
4. Alertas sobre problemas potenciais

%s

Horas desde o último backup: %.1f

Forneça uma análise detalhada em português, sendo claro e objetivo.`, backupSummary.String(), hoursSinceBackup)

	return b.aiClient.GenerateContent(prompt, 1500)
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

func getTime(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

func formatSize(bytes int64) string {
	if bytes == 0 {
		return "N/A"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

