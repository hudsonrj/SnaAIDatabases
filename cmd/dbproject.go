package cmd

import (
	"fmt"

	"github.com/snip/internal/dbproject"
	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var (
	dbProjectAnalysisID int
	dbProjectIncident   string
)

func init() {
	dbProjectCmd.Flags().IntVarP(&dbProjectAnalysisID, "analysis-id", "a", 0, "ID da análise para transformar em projeto")
	dbProjectCmd.Flags().StringVarP(&dbProjectIncident, "incident", "i", "", "Descrição do incidente (opcional)")

	rootCmd.AddCommand(dbProjectCmd)
}

var dbProjectCmd = &cobra.Command{
	Use:   "db-project",
	Short: "Transformar análises ou incidentes em projetos",
	Long: `Transforma resultados de análises ou incidentes em projetos estruturados com tarefas.

A IA cria:
- Projeto com nome e descrição
- Tarefas priorizadas
- Passo a passo detalhado
- Tempo estimado e datas sugeridas

Exemplos:
  snip db-project --analysis-id 1
  snip db-project --analysis-id 1 --incident "Banco de dados lento durante picos"
  snip db-project --incident "Erro de conexão" --analysis-id 2`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbProjectAnalysisID == 0 && dbProjectIncident == "" {
			fmt.Println("Erro: --analysis-id ou --incident é obrigatório")
			return
		}

		// Obter análise se fornecida
		var analysisResult string
		var analysisTitle string
		var analysisType string
		var dbType string

		if dbProjectAnalysisID > 0 {
			if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
				analysis, err := h.GetAnalysisByID(dbProjectAnalysisID)
				if err != nil {
					return fmt.Errorf("erro ao obter análise: %w", err)
				}

				analysisResult = analysis.Result
				analysisTitle = analysis.Title
				analysisType = string(analysis.AnalysisType)
				dbType = string(analysis.DatabaseType)

				return nil
			}); err != nil {
				fmt.Printf("Erro: %v\n", err)
				return
			}
		}

		// Criar gerador de projetos
		generator, err := dbproject.NewProjectGenerator()
		if err != nil {
			fmt.Printf("Erro ao criar gerador: %v\n", err)
			return
		}

		var project *dbproject.ProjectFromAnalysis

		// Gerar projeto
		if dbProjectIncident != "" {
			// Projeto de incidente
			if analysisResult == "" {
				fmt.Println("Erro: --analysis-id é necessário quando usar --incident")
				return
			}
			project, err = generator.GenerateIncidentProject(dbProjectIncident, analysisResult, dbType)
		} else {
			// Projeto de análise
			if analysisResult == "" {
				fmt.Println("Erro: análise ainda não foi executada")
				return
			}
			project, err = generator.GenerateProjectFromAnalysis(analysisTitle, analysisResult, analysisType, dbType)
		}

		if err != nil {
			fmt.Printf("Erro ao gerar projeto: %v\n", err)
			return
		}

		// Exibir projeto
		fmt.Println(project.FormatProject())

		// Perguntar se deseja criar no sistema
		fmt.Println("\n💡 Dica: Use 'snip project create' para criar este projeto no sistema")
	},
}

