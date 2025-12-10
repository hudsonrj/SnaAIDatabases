package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// OCIConfig armazena configuração do OCI
type OCIConfig struct {
	TenancyOCID  string `json:"tenancy_ocid"`
	UserOCID     string `json:"user_ocid"`
	Fingerprint  string `json:"fingerprint"`
	PrivateKey   string `json:"private_key"`
	Region       string `json:"region"`
	CompartmentID string `json:"compartment_id"`
}

// OCIDatabase representa um banco de dados no OCI
type OCIDatabase struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"` // Autonomous Database, MySQL, PostgreSQL
	Status       string            `json:"status"`
	Shape        string            `json:"shape"`
	OCPUs        int               `json:"ocpus"`
	Storage      int               `json:"storage"`
	Region       string            `json:"region"`
	Endpoint     string            `json:"endpoint"`
	Cost         float64           `json:"cost"`
	Tags         map[string]string `json:"tags"`
	BackupEnabled bool             `json:"backup_enabled"`
}

// OCIClient é o cliente para OCI
type OCIClient struct {
	config OCIConfig
}

// NewOCIClient cria um novo cliente OCI
func NewOCIClient(config OCIConfig) *OCIClient {
	return &OCIClient{config: config}
}

// ListDatabases lista todos os bancos de dados Autonomous Database
func (c *OCIClient) ListDatabases() ([]OCIDatabase, error) {
	var databases []OCIDatabase

	// Usar OCI CLI para listar Autonomous Databases
	cmd := exec.Command("oci", "db", "autonomous-database", "list",
		"--compartment-id", c.config.CompartmentID,
		"--output", "json",
		"--region", c.config.Region)

	// Configurar autenticação OCI
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("OCI_TENANCY_OCID=%s", c.config.TenancyOCID))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OCI_USER_OCID=%s", c.config.UserOCID))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OCI_FINGERPRINT=%s", c.config.Fingerprint))
	cmd.Env = append(cmd.Env, fmt.Sprintf("OCI_REGION=%s", c.config.Region))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar Autonomous Databases: %w", err)
	}

	var result struct {
		Data []struct {
			ID              string `json:"id"`
			DisplayName     string `json:"display-name"`
			LifecycleState  string `json:"lifecycle-state"`
			DbWorkload      string `json:"db-workload"`
			IsFreeTier      bool   `json:"is-free-tier"`
			CpuCoreCount    int    `json:"cpu-core-count"`
			DataStorageSizeInTBs int `json:"data-storage-size-in-tbs"`
			ServiceConsoleURL string `json:"service-console-url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta OCI: %w", err)
	}

	for _, db := range result.Data {
		cost := c.estimateCost(db.CpuCoreCount, db.DataStorageSizeInTBs, db.IsFreeTier)

		databases = append(databases, OCIDatabase{
			ID:            db.ID,
			Name:          db.DisplayName,
			Type:          fmt.Sprintf("Autonomous %s", db.DbWorkload),
			Status:        db.LifecycleState,
			Shape:         fmt.Sprintf("%d OCPUs", db.CpuCoreCount),
			OCPUs:         db.CpuCoreCount,
			Storage:       db.DataStorageSizeInTBs,
			Region:        c.config.Region,
			Endpoint:      db.ServiceConsoleURL,
			Cost:          cost,
			BackupEnabled: true, // Autonomous Database tem backup automático
		})
	}

	return databases, nil
}

func (c *OCIClient) estimateCost(ocpus, storage int, isFreeTier bool) float64 {
	if isFreeTier {
		return 0.0
	}

	// Valores aproximados mensais (pode ser melhorado com OCI Pricing API)
	// Autonomous Database: ~$0.30 por OCPU/hora + storage
	cpuCost := float64(ocpus) * 0.30 * 24 * 30
	storageCost := float64(storage) * 200.0 // ~$200 por TB/mês

	return cpuCost + storageCost
}

// GetDatabaseInsights obtém insights sobre um banco de dados
func (c *OCIClient) GetDatabaseInsights(name string) (string, error) {
	insights := fmt.Sprintf("Insights para %s:\n- Métricas disponíveis no OCI Monitoring\n- Performance insights disponíveis\n- Recomendações de otimização", name)
	return insights, nil
}

