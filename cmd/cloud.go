package cmd

import (
	"fmt"
	"os"

	"github.com/snip/internal/cloud"
	"github.com/spf13/cobra"
)

var (
	cloudProvider string
	cloudRegion   string
	cloudProject  string
)

func init() {
	cloudListCmd.Flags().StringVarP(&cloudProvider, "provider", "p", "", "Provedor de cloud (aws, azure, gcp, oci)")
	cloudListCmd.Flags().StringVarP(&cloudRegion, "region", "r", "", "Região (opcional)")
	cloudListCmd.Flags().StringVarP(&cloudProject, "project", "j", "", "Project ID (GCP) ou Subscription ID (Azure)")

	rootCmd.AddCommand(cloudCmd)
	cloudCmd.AddCommand(cloudListCmd)
}

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Integração com clouds (AWS, Azure, GCP, OCI)",
	Long: `Comandos para integração com provedores de cloud.

Suporta:
  - AWS: RDS, DocumentDB, Redshift
  - Azure: SQL Database, PostgreSQL, MySQL, Cosmos DB
  - GCP: Cloud SQL, Firestore, BigQuery
  - OCI: Autonomous Database, MySQL, PostgreSQL`,
}

var cloudListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar bancos de dados na cloud",
	Long: `Lista todos os bancos de dados em um provedor de cloud.

Exemplos:
  snip cloud list --provider aws --region us-east-1
  snip cloud list --provider azure --project "subscription-id"
  snip cloud list --provider gcp --project "my-project"
  snip cloud list --provider oci --region us-ashburn-1`,
	Run: func(cmd *cobra.Command, args []string) {
		if cloudProvider == "" {
			fmt.Println("❌ Provedor é obrigatório (use --provider)")
			return
		}

		switch cloudProvider {
		case "aws":
			listAWSDatabases()
		case "azure":
			listAzureDatabases()
		case "gcp":
			listGCPDatabases()
		case "oci":
			listOCIDatabases()
		default:
			fmt.Printf("❌ Provedor não suportado: %s\n", cloudProvider)
			fmt.Println("Provedores suportados: aws, azure, gcp, oci")
		}
	},
}

func listAWSDatabases() {
	config := cloud.AWSConfig{
		Region:  cloudRegion,
		Profile: os.Getenv("AWS_PROFILE"),
	}

	if config.Region == "" {
		config.Region = "us-east-1"
	}

	client := cloud.NewAWSClient(config)
	databases, err := client.ListDatabases()
	if err != nil {
		fmt.Printf("❌ Erro ao listar bancos AWS: %v\n", err)
		return
	}

	fmt.Printf("📊 Bancos de Dados AWS (%d encontrados):\n\n", len(databases))
	for _, db := range databases {
		fmt.Printf("🔹 %s\n", db.Identifier)
		fmt.Printf("   Engine: %s %s\n", db.Engine, db.EngineVersion)
		fmt.Printf("   Status: %s\n", db.Status)
		fmt.Printf("   Classe: %s\n", db.InstanceClass)
		fmt.Printf("   Storage: %d GB\n", db.Storage)
		fmt.Printf("   Multi-AZ: %v\n", db.MultiAZ)
		fmt.Printf("   Endpoint: %s:%d\n", db.Endpoint, db.Port)
		fmt.Printf("   Custo estimado: $%.2f/mês\n", db.Cost)
		fmt.Printf("   Backup retention: %d dias\n", db.BackupRetention)
		fmt.Println()
	}
}

func listAzureDatabases() {
	config := cloud.AzureConfig{
		SubscriptionID: cloudProject,
		ResourceGroup:  "",
	}

	if config.SubscriptionID == "" {
		config.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}

	client := cloud.NewAzureClient(config)
	databases, err := client.ListDatabases()
	if err != nil {
		fmt.Printf("❌ Erro ao listar bancos Azure: %v\n", err)
		return
	}

	fmt.Printf("📊 Bancos de Dados Azure (%d encontrados):\n\n", len(databases))
	for _, db := range databases {
		fmt.Printf("🔹 %s\n", db.Name)
		fmt.Printf("   Tipo: %s\n", db.Type)
		fmt.Printf("   Status: %s\n", db.Status)
		fmt.Printf("   Tier: %s\n", db.Tier)
		fmt.Printf("   Size: %s\n", db.Size)
		fmt.Printf("   Location: %s\n", db.Location)
		fmt.Printf("   Endpoint: %s\n", db.Endpoint)
		fmt.Printf("   Custo estimado: $%.2f/mês\n", db.Cost)
		fmt.Printf("   Backup: %v\n", db.BackupEnabled)
		fmt.Println()
	}
}

func listGCPDatabases() {
	config := cloud.GCPConfig{
		ProjectID: cloudProject,
		Region:    cloudRegion,
	}

	if config.ProjectID == "" {
		config.ProjectID = os.Getenv("GCP_PROJECT")
	}

	if config.ProjectID == "" {
		fmt.Println("❌ Project ID é obrigatório (use --project ou GCP_PROJECT)")
		return
	}

	client := cloud.NewGCPClient(config)
	databases, err := client.ListDatabases()
	if err != nil {
		fmt.Printf("❌ Erro ao listar bancos GCP: %v\n", err)
		return
	}

	fmt.Printf("📊 Bancos de Dados GCP (%d encontrados):\n\n", len(databases))
	for _, db := range databases {
		fmt.Printf("🔹 %s\n", db.Name)
		fmt.Printf("   Tipo: %s\n", db.Type)
		fmt.Printf("   Status: %s\n", db.Status)
		fmt.Printf("   Tier: %s\n", db.Tier)
		fmt.Printf("   Region: %s\n", db.Region)
		fmt.Printf("   Endpoint: %s\n", db.Endpoint)
		fmt.Printf("   Custo estimado: $%.2f/mês\n", db.Cost)
		fmt.Printf("   Backup: %v\n", db.BackupEnabled)
		fmt.Println()
	}
}

func listOCIDatabases() {
	config := cloud.OCIConfig{
		Region:        cloudRegion,
		CompartmentID: os.Getenv("OCI_COMPARTMENT_ID"),
		TenancyOCID:   os.Getenv("OCI_TENANCY_OCID"),
		UserOCID:      os.Getenv("OCI_USER_OCID"),
		Fingerprint:   os.Getenv("OCI_FINGERPRINT"),
	}

	if config.Region == "" {
		config.Region = "us-ashburn-1"
	}

	if config.CompartmentID == "" {
		fmt.Println("❌ Compartment ID é obrigatório (configure OCI_COMPARTMENT_ID)")
		return
	}

	client := cloud.NewOCIClient(config)
	databases, err := client.ListDatabases()
	if err != nil {
		fmt.Printf("❌ Erro ao listar bancos OCI: %v\n", err)
		return
	}

	fmt.Printf("📊 Bancos de Dados OCI (%d encontrados):\n\n", len(databases))
	for _, db := range databases {
		fmt.Printf("🔹 %s\n", db.Name)
		fmt.Printf("   Tipo: %s\n", db.Type)
		fmt.Printf("   Status: %s\n", db.Status)
		fmt.Printf("   Shape: %s\n", db.Shape)
		fmt.Printf("   OCPUs: %d\n", db.OCPUs)
		fmt.Printf("   Storage: %d TB\n", db.Storage)
		fmt.Printf("   Region: %s\n", db.Region)
		fmt.Printf("   Custo estimado: $%.2f/mês\n", db.Cost)
		fmt.Printf("   Backup: %v\n", db.BackupEnabled)
		fmt.Println()
	}
}

