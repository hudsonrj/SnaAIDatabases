package dbtypes

// DatabaseType representa o tipo de banco de dados
type DatabaseType string

const (
	DatabaseTypeOracle    DatabaseType = "oracle"
	DatabaseTypeSQLServer DatabaseType = "sqlserver"
	DatabaseTypeMySQL     DatabaseType = "mysql"
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMongoDB   DatabaseType = "mongodb"
)

// ConnectionConfig representa a configuração de conexão com o banco
type ConnectionConfig struct {
	Type            DatabaseType `json:"type"`
	Host            string       `json:"host"`
	Port            int          `json:"port"`
	Database        string       `json:"database"`
	Username        string       `json:"username"`
	Password        string       `json:"password"`
	IsRemote        bool         `json:"is_remote"`
	JDBCURL         string       `json:"jdbc_url,omitempty"`
	ConnectionString string      `json:"connection_string,omitempty"`
}

