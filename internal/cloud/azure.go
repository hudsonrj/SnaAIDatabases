package cloud

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// AzureConfig armazena configuração do Azure
type AzureConfig struct {
	SubscriptionID string `json:"subscription_id"`
	TenantID       string `json:"tenant_id"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	ResourceGroup  string `json:"resource_group"`
}

// AzureDatabase representa um banco de dados no Azure
type AzureDatabase struct {
	Name           string            `json:"name"`
	Type           string            `json:"type"` // SQL, PostgreSQL, MySQL, CosmosDB
	Status         string            `json:"status"`
	Tier           string            `json:"tier"`
	Size           string            `json:"size"`
	Location       string            `json:"location"`
	Endpoint       string            `json:"endpoint"`
	Cost           float64           `json:"cost"`
	Tags           map[string]string `json:"tags"`
	BackupEnabled  bool              `json:"backup_enabled"`
}

// AzureClient é o cliente para Azure
type AzureClient struct {
	config AzureConfig
}

// NewAzureClient cria um novo cliente Azure
func NewAzureClient(config AzureConfig) *AzureClient {
	return &AzureClient{config: config}
}

// ListDatabases lista todos os bancos de dados
func (c *AzureClient) ListDatabases() ([]AzureDatabase, error) {
	var databases []AzureDatabase

	// Usar Azure CLI para listar recursos
	cmd := exec.Command("az", "sql", "server", "list",
		"--output", "json",
		"--subscription", c.config.SubscriptionID)

	if c.config.ResourceGroup != "" {
		cmd.Args = append(cmd.Args, "--resource-group", c.config.ResourceGroup)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar bancos SQL: %w", err)
	}

	var servers []struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		State    string `json:"state"`
		Tags     map[string]interface{} `json:"tags"`
	}

	if err := json.Unmarshal(output, &servers); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta Azure: %w", err)
	}

	for _, server := range servers {
		// Obter databases do servidor
		dbCmd := exec.Command("az", "sql", "db", "list",
			"--server", server.Name,
			"--output", "json",
			"--subscription", c.config.SubscriptionID)

		dbOutput, err := dbCmd.Output()
		if err != nil {
			continue
		}

		var dbs []struct {
			Name           string `json:"name"`
			Status         string `json:"status"`
			ServiceLevelObjective string `json:"serviceLevelObjective"`
			CurrentServiceObjectiveName string `json:"currentServiceObjectiveName"`
		}

		if err := json.Unmarshal(dbOutput, &dbs); err != nil {
			continue
		}

		for _, db := range dbs {
			tags := make(map[string]string)
			for k, v := range server.Tags {
				if str, ok := v.(string); ok {
					tags[k] = str
				}
			}

			cost := c.estimateCost(db.ServiceLevelObjective, db.CurrentServiceObjectiveName)

			databases = append(databases, AzureDatabase{
				Name:          fmt.Sprintf("%s/%s", server.Name, db.Name),
				Type:          "SQL",
				Status:        db.Status,
				Tier:          db.ServiceLevelObjective,
				Size:          db.CurrentServiceObjectiveName,
				Location:      server.Location,
				Endpoint:      fmt.Sprintf("%s.database.windows.net", server.Name),
				Cost:          cost,
				Tags:          tags,
				BackupEnabled: true, // Azure SQL tem backup automático
			})
		}
	}

	// Listar PostgreSQL e MySQL também
	postgresServers, _ := c.listPostgreSQLServers()
	databases = append(databases, postgresServers...)

	mysqlServers, _ := c.listMySQLServers()
	databases = append(databases, mysqlServers...)

	return databases, nil
}

func (c *AzureClient) listPostgreSQLServers() ([]AzureDatabase, error) {
	var databases []AzureDatabase

	cmd := exec.Command("az", "postgres", "server", "list",
		"--output", "json",
		"--subscription", c.config.SubscriptionID)

	output, err := cmd.Output()
	if err != nil {
		return databases, nil // Não é erro crítico
	}

	var servers []struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		State    string `json:"userVisibleState"`
		Sku      struct {
			Name string `json:"name"`
			Tier string `json:"tier"`
		} `json:"sku"`
	}

	if err := json.Unmarshal(output, &servers); err != nil {
		return databases, nil
	}

	for _, server := range servers {
		cost := c.estimatePostgreSQLCost(server.Sku.Tier, server.Sku.Name)

		databases = append(databases, AzureDatabase{
			Name:     server.Name,
			Type:     "PostgreSQL",
			Status:   server.State,
			Tier:     server.Sku.Tier,
			Size:     server.Sku.Name,
			Location: server.Location,
			Endpoint: fmt.Sprintf("%s.postgres.database.azure.com", server.Name),
			Cost:     cost,
		})
	}

	return databases, nil
}

func (c *AzureClient) listMySQLServers() ([]AzureDatabase, error) {
	var databases []AzureDatabase

	cmd := exec.Command("az", "mysql", "server", "list",
		"--output", "json",
		"--subscription", c.config.SubscriptionID)

	output, err := cmd.Output()
	if err != nil {
		return databases, nil
	}

	var servers []struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		State    string `json:"userVisibleState"`
		Sku      struct {
			Name string `json:"name"`
			Tier string `json:"tier"`
		} `json:"sku"`
	}

	if err := json.Unmarshal(output, &servers); err != nil {
		return databases, nil
	}

	for _, server := range servers {
		cost := c.estimateMySQLCost(server.Sku.Tier, server.Sku.Name)

		databases = append(databases, AzureDatabase{
			Name:     server.Name,
			Type:     "MySQL",
			Status:   server.State,
			Tier:     server.Sku.Tier,
			Size:     server.Sku.Name,
			Location: server.Location,
			Endpoint: fmt.Sprintf("%s.mysql.database.azure.com", server.Name),
			Cost:     cost,
		})
	}

	return databases, nil
}

func (c *AzureClient) estimateCost(tier, size string) float64 {
	// Valores aproximados mensais (pode ser melhorado com Azure Pricing API)
	baseCost := 0.0
	if strings.Contains(tier, "Basic") {
		baseCost = 5.0
	} else if strings.Contains(tier, "Standard") {
		baseCost = 15.0
	} else if strings.Contains(tier, "Premium") {
		baseCost = 465.0
	}
	return baseCost
}

func (c *AzureClient) estimatePostgreSQLCost(tier, size string) float64 {
	if strings.Contains(tier, "Basic") {
		return 25.0
	} else if strings.Contains(tier, "GeneralPurpose") {
		return 200.0
	} else if strings.Contains(tier, "MemoryOptimized") {
		return 500.0
	}
	return 100.0
}

func (c *AzureClient) estimateMySQLCost(tier, size string) float64 {
	if strings.Contains(tier, "Basic") {
		return 25.0
	} else if strings.Contains(tier, "GeneralPurpose") {
		return 200.0
	} else if strings.Contains(tier, "MemoryOptimized") {
		return 500.0
	}
	return 100.0
}

// GetDatabaseInsights obtém insights sobre um banco de dados
func (c *AzureClient) GetDatabaseInsights(name string) (string, error) {
	insights := fmt.Sprintf("Insights para %s:\n- Métricas disponíveis no Azure Monitor\n- Recomendações de otimização disponíveis", name)
	return insights, nil
}

