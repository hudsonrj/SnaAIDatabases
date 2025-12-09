package cmd

import (
	"fmt"
	"os"

	"github.com/snip/internal/dbmaintenance"
	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var (
	dbMaintenanceAnalysisID int
	dbMaintenanceOutputFile string
)

func init() {
	dbMaintenanceCmd.Flags().IntVarP(&dbMaintenanceAnalysisID, "analysis-id", "a", 0, "ID da análise para gerar plano")
	dbMaintenanceCmd.Flags().StringVarP(&dbMaintenanceOutputFile, "output", "o", "", "Arquivo de saída (opcional)")

	rootCmd.AddCommand(dbMaintenanceCmd)
}

var dbMaintenanceCmd = &cobra.Command{
	Use:   "db-maintenance",
	Short: "Gerar planos de manutenção baseados em análises",
	Long: `Gera planos de manutenção detalhados usando IA baseados em resultados de análises.

O plano inclui:
- Tarefas priorizadas
- Passo a passo detalhado
- Tempo estimado
- Dependências entre tarefas

Exemplos:
  snip db-maintenance --analysis-id 1
  snip db-maintenance --analysis-id 1 --output maintenance-plan.md`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbMaintenanceAnalysisID == 0 {
			fmt.Println("Erro: --analysis-id é obrigatório")
			return
		}

		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			// Obter análise
			analysis, err := h.GetAnalysisByID(dbMaintenanceAnalysisID)
			if err != nil {
				return fmt.Errorf("erro ao obter análise: %w", err)
			}

			if analysis.Result == "" {
				return fmt.Errorf("análise ainda não foi executada")
			}

			// Criar planejador
			planner, err := dbmaintenance.NewMaintenancePlanner()
			if err != nil {
				return fmt.Errorf("erro ao criar planejador: %w", err)
			}

			// Gerar plano
			plan, err := planner.GenerateMaintenancePlan(
				analysis.Result,
				string(analysis.AnalysisType),
				string(analysis.DatabaseType),
			)
			if err != nil {
				return fmt.Errorf("erro ao gerar plano: %w", err)
			}

			// Formatar e exibir
			formatted := plan.FormatPlan()
			fmt.Println(formatted)

			// Salvar se solicitado
			if dbMaintenanceOutputFile != "" {
				err := os.WriteFile(dbMaintenanceOutputFile, []byte(formatted), 0644)
				if err != nil {
					return fmt.Errorf("erro ao salvar arquivo: %w", err)
				}
				fmt.Printf("\nPlano salvo em: %s\n", dbMaintenanceOutputFile)
			}

			return nil
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

