package cloud

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GCPConfig armazena configuração do GCP
type GCPConfig struct {
	ProjectID string `json:"project_id"`
	Region    string `json:"region"`
	Zone      string `json:"zone"`
}

// GCPDatabase representa um banco de dados no GCP
type GCPDatabase struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"` // Cloud SQL, Firestore, BigQuery
	Status       string            `json:"status"`
	Tier         string            `json:"tier"`
	Region       string            `json:"region"`
	Endpoint     string            `json:"endpoint"`
	Cost         float64           `json:"cost"`
	Tags         map[string]string `json:"tags"`
	BackupEnabled bool             `json:"backup_enabled"`
}

// GCPClient é o cliente para GCP
type GCPClient struct {
	config GCPConfig
}

// NewGCPClient cria um novo cliente GCP
func NewGCPClient(config GCPConfig) *GCPClient {
	return &GCPClient{config: config}
}

// ListDatabases lista todos os bancos de dados Cloud SQL
func (c *GCPClient) ListDatabases() ([]GCPDatabase, error) {
	var databases []GCPDatabase

	// Usar gcloud CLI para listar instâncias Cloud SQL
	cmd := exec.Command("gcloud", "sql", "instances", "list",
		"--format", "json",
		"--project", c.config.ProjectID)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar instâncias Cloud SQL: %w", err)
	}

	var instances []struct {
		Name             string `json:"name"`
		DatabaseVersion string `json:"databaseVersion"`
		State            string `json:"state"`
		Region           string `json:"region"`
		Settings         struct {
			Tier            string `json:"tier"`
			BackupConfiguration struct {
				Enabled bool `json:"enabled"`
			} `json:"backupConfiguration"`
		} `json:"settings"`
		IpAddresses []struct {
			IpAddress string `json:"ipAddress"`
		} `json:"ipAddresses"`
	}

	if err := json.Unmarshal(output, &instances); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta GCP: %w", err)
	}

	for _, instance := range instances {
		endpoint := ""
		if len(instance.IpAddresses) > 0 {
			endpoint = instance.IpAddresses[0].IpAddress
		}

		cost := c.estimateCost(instance.Settings.Tier)

		databases = append(databases, GCPDatabase{
			Name:          instance.Name,
			Type:          instance.DatabaseVersion,
			Status:        instance.State,
			Tier:          instance.Settings.Tier,
			Region:        instance.Region,
			Endpoint:      endpoint,
			Cost:          cost,
			BackupEnabled: instance.Settings.BackupConfiguration.Enabled,
		})
	}

	return databases, nil
}

func (c *GCPClient) estimateCost(tier string) float64 {
	// Valores aproximados mensais (pode ser melhorado com GCP Pricing API)
	baseCost := 0.0
	if strings.Contains(tier, "db-f1-micro") {
		baseCost = 7.67
	} else if strings.Contains(tier, "db-g1-small") {
		baseCost = 25.56
	} else if strings.Contains(tier, "db-n1-standard-1") {
		baseCost = 51.11
	} else if strings.Contains(tier, "db-n1-standard-2") {
		baseCost = 102.22
	} else if strings.Contains(tier, "db-n1-highmem-2") {
		baseCost = 122.67
	} else {
		baseCost = 50.0 // Default
	}
	return baseCost
}

// GetDatabaseInsights obtém insights sobre um banco de dados
func (c *GCPClient) GetDatabaseInsights(name string) (string, error) {
	insights := fmt.Sprintf("Insights para %s:\n- Métricas disponíveis no Cloud Monitoring\n- Recomendações de otimização disponíveis", name)
	return insights, nil
}

