package dbconnection

import (
	"database/sql"
	"fmt"
	"os/user"
	"strings"

	"github.com/snip/internal/dbtypes"
)

// DBConnector interface para diferentes tipos de banco de dados
type DBConnector interface {
	Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error)
	GetConnectionString(config *dbtypes.ConnectionConfig) (string, error)
	TestConnection(config *dbtypes.ConnectionConfig) error
}

// GetConnector retorna o connector apropriado para o tipo de banco
func GetConnector(dbType dbtypes.DatabaseType) (DBConnector, error) {
	switch dbType {
	case dbtypes.DatabaseTypeMySQL:
		return &MySQLConnector{}, nil
	case dbtypes.DatabaseTypePostgreSQL:
		return &PostgreSQLConnector{}, nil
	case dbtypes.DatabaseTypeOracle:
		return &OracleConnector{}, nil
	case dbtypes.DatabaseTypeSQLServer:
		return &SQLServerConnector{}, nil
	case dbtypes.DatabaseTypeMongoDB:
		return &MongoDBConnector{}, nil
	default:
		return nil, fmt.Errorf("tipo de banco de dados não suportado: %s", dbType)
	}
}

// MySQLConnector implementa conexão MySQL
type MySQLConnector struct{}

func (m *MySQLConnector) Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error) {
	_, err := m.GetConnectionString(config)
	if err != nil {
		return nil, err
	}
	// Nota: requer driver MySQL como "github.com/go-sql-driver/mysql"
	// db, err := sql.Open("mysql", connStr)
	// return db, err
	return nil, fmt.Errorf("driver MySQL não instalado. Execute: go get github.com/go-sql-driver/mysql")
}

func (m *MySQLConnector) GetConnectionString(config *dbtypes.ConnectionConfig) (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.JDBCURL != "" {
		// Converter JDBC URL para formato Go
		return convertJDBCToGo(config.JDBCURL, "mysql"), nil
	}
	port := config.Port
	if port == 0 {
		port = 3306
	}

	// Verificar se é conexão local e se pode usar autenticação sem senha
	isLocal := isLocalConnection(config.Host)
	usePasswordless := isLocal && config.Password == "" && canUsePasswordlessAuth(config.Username, "mysql")

	if usePasswordless {
		// MySQL local pode usar socket Unix ou autenticação sem senha
		// Tentar primeiro com socket Unix se disponível
		if config.Username != "" {
			return fmt.Sprintf("%s@unix(/tmp/mysql.sock)/%s?parseTime=true",
				config.Username, config.Database), nil
		}
		// Fallback: tentar sem senha
		return fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true&allowNativePasswordPassThrough=true",
			config.Username, config.Host, port, config.Database), nil
	}

	// Com senha ou conexão remota
	if config.Password != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			config.Username, config.Password, config.Host, port, config.Database), nil
	}

	// Tentar sem senha (pode funcionar se configurado)
	return fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username, config.Host, port, config.Database), nil
}

func (m *MySQLConnector) TestConnection(config *dbtypes.ConnectionConfig) error {
	db, err := m.Connect(config)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// PostgreSQLConnector implementa conexão PostgreSQL
type PostgreSQLConnector struct{}

func (p *PostgreSQLConnector) Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error) {
	_, err := p.GetConnectionString(config)
	if err != nil {
		return nil, err
	}
	// Nota: requer driver PostgreSQL como "github.com/lib/pq"
	// db, err := sql.Open("postgres", connStr)
	// return db, err
	return nil, fmt.Errorf("driver PostgreSQL não instalado. Execute: go get github.com/lib/pq")
}

func (p *PostgreSQLConnector) GetConnectionString(config *dbtypes.ConnectionConfig) (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.JDBCURL != "" {
		return convertJDBCToGo(config.JDBCURL, "postgres"), nil
	}
	port := config.Port
	if port == 0 {
		port = 5432
	}

	// Verificar se é conexão local e se pode usar autenticação sem senha
	isLocal := isLocalConnection(config.Host)
	usePasswordless := isLocal && config.Password == "" && canUsePasswordlessAuth(config.Username, "postgresql")

	var connStr strings.Builder
	connStr.WriteString(fmt.Sprintf("host=%s port=%d", config.Host, port))
	
	if config.Username != "" {
		connStr.WriteString(fmt.Sprintf(" user=%s", config.Username))
	}
	
	if usePasswordless {
		// PostgreSQL local pode usar peer authentication ou trust
		connStr.WriteString(" sslmode=disable")
	} else if config.Password != "" {
		connStr.WriteString(fmt.Sprintf(" password=%s sslmode=disable", config.Password))
	} else {
		// Tentar sem senha mesmo assim (pode funcionar com peer/trust)
		connStr.WriteString(" sslmode=disable")
	}
	
	if config.Database != "" {
		connStr.WriteString(fmt.Sprintf(" dbname=%s", config.Database))
	}

	return connStr.String(), nil
}

func (p *PostgreSQLConnector) TestConnection(config *dbtypes.ConnectionConfig) error {
	db, err := p.Connect(config)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// OracleConnector implementa conexão Oracle
type OracleConnector struct{}

func (o *OracleConnector) Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error) {
	_, err := o.GetConnectionString(config)
	if err != nil {
		return nil, err
	}
	// Nota: requer driver Oracle como "github.com/godror/godror" ou "github.com/sijms/go-ora"
	// db, err := sql.Open("godror", connStr)
	// return db, err
	return nil, fmt.Errorf("driver Oracle não instalado. Execute: go get github.com/godror/godror")
}

func (o *OracleConnector) GetConnectionString(config *dbtypes.ConnectionConfig) (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.JDBCURL != "" {
		return convertJDBCToGo(config.JDBCURL, "oracle"), nil
	}
	port := config.Port
	if port == 0 {
		port = 1521
	}

	// Verificar se é conexão local e se pode usar autenticação sem senha
	isLocal := isLocalConnection(config.Host)
	usePasswordless := isLocal && config.Password == "" && canUsePasswordlessAuth(config.Username, "oracle")

	if usePasswordless {
		// Oracle local pode usar OS authentication (/ as sysdba ou / as sysoper)
		// Ou conectar como o usuário do OS
		if config.Username == "" || config.Username == "sys" {
			// Tentar OS authentication
			return fmt.Sprintf("/@%s:%d/%s", config.Host, port, config.Database), nil
		}
		// Tentar conectar como usuário do OS
		return fmt.Sprintf("/@%s:%d/%s", config.Host, port, config.Database), nil
	}

	// Com senha ou conexão remota
	if config.Password != "" {
		return fmt.Sprintf("%s/%s@%s:%d/%s",
			config.Username, config.Password, config.Host, port, config.Database), nil
	}

	// Tentar sem senha (pode funcionar com OS auth)
	return fmt.Sprintf("/@%s:%d/%s", config.Host, port, config.Database), nil
}

func (o *OracleConnector) TestConnection(config *dbtypes.ConnectionConfig) error {
	db, err := o.Connect(config)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// SQLServerConnector implementa conexão SQL Server
type SQLServerConnector struct{}

func (s *SQLServerConnector) Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error) {
	_, err := s.GetConnectionString(config)
	if err != nil {
		return nil, err
	}
	// Nota: requer driver SQL Server como "github.com/denisenkom/go-mssqldb"
	// db, err := sql.Open("sqlserver", connStr)
	// return db, err
	return nil, fmt.Errorf("driver SQL Server não instalado. Execute: go get github.com/denisenkom/go-mssqldb")
}

func (s *SQLServerConnector) GetConnectionString(config *dbtypes.ConnectionConfig) (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.JDBCURL != "" {
		return convertJDBCToGo(config.JDBCURL, "sqlserver"), nil
	}
	port := config.Port
	if port == 0 {
		port = 1433
	}

	// Verificar se é conexão local e se pode usar autenticação sem senha
	isLocal := isLocalConnection(config.Host)
	usePasswordless := isLocal && config.Password == "" && canUsePasswordlessAuth(config.Username, "sqlserver")

	var connStr strings.Builder
	connStr.WriteString(fmt.Sprintf("server=%s;port=%d", config.Host, port))

	if usePasswordless {
		// SQL Server local pode usar Windows Authentication (Integrated Security)
		connStr.WriteString(";Integrated Security=true")
		if config.Database != "" {
			connStr.WriteString(fmt.Sprintf(";database=%s", config.Database))
		}
		return connStr.String(), nil
	}

	// Com senha ou conexão remota
	if config.Username != "" {
		connStr.WriteString(fmt.Sprintf(";user id=%s", config.Username))
	}
	if config.Password != "" {
		connStr.WriteString(fmt.Sprintf(";password=%s", config.Password))
	}
	if config.Database != "" {
		connStr.WriteString(fmt.Sprintf(";database=%s", config.Database))
	}

	return connStr.String(), nil
}

func (s *SQLServerConnector) TestConnection(config *dbtypes.ConnectionConfig) error {
	db, err := s.Connect(config)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// MongoDBConnector implementa conexão MongoDB
type MongoDBConnector struct{}

// isLocalConnection verifica se a conexão é local
func isLocalConnection(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return host == "localhost" || host == "127.0.0.1" || host == "::1" || host == ""
}

// canUsePasswordlessAuth verifica se pode usar autenticação sem senha
func canUsePasswordlessAuth(username, dbType string) bool {
	// Verificar se estamos no mesmo servidor do banco
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	// Para PostgreSQL: se o usuário do OS é o mesmo do banco ou postgres
	if dbType == "postgresql" {
		if username == "" || username == currentUser.Username || username == "postgres" {
			// Verificar se existe arquivo .pgpass ou se pg_hba.conf permite trust/peer
			return true
		}
	}

	// Para MySQL: se o usuário do OS é o mesmo do banco ou root
	if dbType == "mysql" {
		if username == "" || username == currentUser.Username || username == "root" {
			return true
		}
	}

	// Para Oracle: se o usuário do OS está no grupo dba ou oas
	if dbType == "oracle" {
		groups, err := currentUser.GroupIds()
		if err == nil {
			for _, gid := range groups {
				group, err := user.LookupGroupId(gid)
				if err == nil {
					if group.Name == "dba" || group.Name == "oinstall" || group.Name == "oracle" {
						return true
					}
				}
			}
		}
		// Também verificar se username é sys ou system
		if username == "sys" || username == "system" {
			return true
		}
	}

	// Para SQL Server: Windows Authentication sempre disponível localmente
	if dbType == "sqlserver" {
		// Windows Authentication funciona localmente
		return true
	}

	return false
}

func (m *MongoDBConnector) Connect(config *dbtypes.ConnectionConfig) (*sql.DB, error) {
	// MongoDB não usa SQL, então retornamos erro
	// Para MongoDB, seria necessário usar o driver oficial: "go.mongodb.org/mongo-driver"
	return nil, fmt.Errorf("MongoDB não usa SQL. Use o driver oficial: go.mongodb.org/mongo-driver")
}

func (m *MongoDBConnector) GetConnectionString(config *dbtypes.ConnectionConfig) (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.JDBCURL != "" {
		return config.JDBCURL, nil
	}
	port := config.Port
	if port == 0 {
		port = 27017
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		config.Username, config.Password, config.Host, port, config.Database), nil
}

func (m *MongoDBConnector) TestConnection(config *dbtypes.ConnectionConfig) error {
	// MongoDB requer implementação específica com o driver oficial
	return fmt.Errorf("teste de conexão MongoDB requer driver oficial")
}

// convertJDBCToGo converte uma URL JDBC para formato de string de conexão Go
func convertJDBCToGo(jdbcURL, dbType string) string {
	// Implementação básica - pode ser expandida conforme necessário
	// Exemplo: jdbc:mysql://localhost:3306/dbname -> mysql://localhost:3306/dbname
	return jdbcURL
}

