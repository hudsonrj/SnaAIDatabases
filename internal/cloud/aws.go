package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// AWSConfig armazena configuração do AWS
type AWSConfig struct {
	Region      string `json:"region"`
	Profile     string `json:"profile"`
	AccessKeyID string `json:"access_key_id"`
	SecretKey   string `json:"secret_key"`
}

// AWSDatabase representa um banco de dados na AWS
type AWSDatabase struct {
	Identifier     string            `json:"identifier"`
	Engine         string            `json:"engine"`
	EngineVersion  string            `json:"engine_version"`
	Status         string            `json:"status"`
	InstanceClass  string            `json:"instance_class"`
	Storage        int               `json:"storage"`
	MultiAZ        bool              `json:"multi_az"`
	Endpoint       string            `json:"endpoint"`
	Port           int               `json:"port"`
	Cost           float64           `json:"cost"`
	Tags           map[string]string `json:"tags"`
	BackupRetention int              `json:"backup_retention"`
}

// AWSClient é o cliente para AWS
type AWSClient struct {
	config AWSConfig
}

// NewAWSClient cria um novo cliente AWS
func NewAWSClient(config AWSConfig) *AWSClient {
	return &AWSClient{config: config}
}

// ListDatabases lista todos os bancos de dados RDS
func (c *AWSClient) ListDatabases() ([]AWSDatabase, error) {
	var databases []AWSDatabase

	// Usar AWS CLI para listar instâncias RDS
	cmd := exec.Command("aws", "rds", "describe-db-instances",
		"--output", "json",
		"--region", c.config.Region)

	if c.config.Profile != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("AWS_PROFILE=%s", c.config.Profile))
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar bancos RDS: %w", err)
	}

	var result struct {
		DBInstances []struct {
			DBInstanceIdentifier string `json:"DBInstanceIdentifier"`
			Engine               string `json:"Engine"`
			EngineVersion        string `json:"EngineVersion"`
			DBInstanceStatus     string `json:"DBInstanceStatus"`
			DBInstanceClass      string `json:"DBInstanceClass"`
			AllocatedStorage     int    `json:"AllocatedStorage"`
			MultiAZ              bool   `json:"MultiAZ"`
			Endpoint             struct {
				Address string `json:"Address"`
				Port    int    `json:"Port"`
			} `json:"Endpoint"`
			BackupRetentionPeriod int `json:"BackupRetentionPeriod"`
			TagList               []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"TagList"`
		} `json:"DBInstances"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta AWS: %w", err)
	}

	for _, db := range result.DBInstances {
		tags := make(map[string]string)
		for _, tag := range db.TagList {
			tags[tag.Key] = tag.Value
		}

		// Estimar custo (simplificado - pode ser melhorado com Pricing API)
		cost := c.estimateCost(db.DBInstanceClass, db.MultiAZ, db.AllocatedStorage)

		databases = append(databases, AWSDatabase{
			Identifier:      db.DBInstanceIdentifier,
			Engine:          db.Engine,
			EngineVersion:   db.EngineVersion,
			Status:          db.DBInstanceStatus,
			InstanceClass:   db.DBInstanceClass,
			Storage:         db.AllocatedStorage,
			MultiAZ:         db.MultiAZ,
			Endpoint:        db.Endpoint.Address,
			Port:            db.Endpoint.Port,
			Cost:            cost,
			Tags:            tags,
			BackupRetention: db.BackupRetentionPeriod,
		})
	}

	return databases, nil
}

// estimateCost estima o custo mensal (simplificado)
func (c *AWSClient) estimateCost(instanceClass string, multiAZ bool, storage int) float64 {
	// Valores aproximados por hora (pode ser melhorado com Pricing API)
	baseCost := 0.0
	if strings.Contains(instanceClass, "db.t3.micro") {
		baseCost = 0.017
	} else if strings.Contains(instanceClass, "db.t3.small") {
		baseCost = 0.034
	} else if strings.Contains(instanceClass, "db.t3.medium") {
		baseCost = 0.068
	} else if strings.Contains(instanceClass, "db.r5.large") {
		baseCost = 0.24
	} else if strings.Contains(instanceClass, "db.r5.xlarge") {
		baseCost = 0.48
	} else {
		baseCost = 0.1 // Default
	}

	if multiAZ {
		baseCost *= 2
	}

	// Adicionar custo de storage (GB/mês)
	storageCost := float64(storage) * 0.115

	// Custo mensal
	monthlyCost := (baseCost * 24 * 30) + storageCost

	return monthlyCost
}

// GetDatabaseInsights obtém insights sobre um banco de dados
func (c *AWSClient) GetDatabaseInsights(identifier string) (string, error) {
	// Obter métricas do CloudWatch
	cmd := exec.Command("aws", "cloudwatch", "get-metric-statistics",
		"--namespace", "AWS/RDS",
		"--metric-name", "CPUUtilization",
		"--dimensions", fmt.Sprintf("Name=DBInstanceIdentifier,Value=%s", identifier),
		"--start-time", "2024-01-01T00:00:00Z",
		"--end-time", "2024-01-02T00:00:00Z",
		"--period", "3600",
		"--statistics", "Average",
		"--output", "json",
		"--region", c.config.Region)

	if c.config.Profile != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("AWS_PROFILE=%s", c.config.Profile))
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter métricas: %w", err)
	}

	// Processar métricas e gerar insights
	insights := fmt.Sprintf("Insights para %s:\n- Métricas obtidas do CloudWatch\n- Dados: %s", identifier, string(output))

	return insights, nil
}

