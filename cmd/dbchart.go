package cmd

import (
	"fmt"
	"os"

	"github.com/snip/internal/dbcharts"
	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var (
	dbChartAnalysisID  int
	dbChartChartType   string
	dbChartOutputFile  string
)

func init() {
	dbChartCmd.Flags().IntVarP(&dbChartAnalysisID, "analysis-id", "a", 0, "ID da análise para gerar gráfico")
	dbChartCmd.Flags().StringVarP(&dbChartChartType, "type", "t", "", "Tipo de gráfico (line, bar, pie, area, table, ascii, html)")
	dbChartCmd.Flags().StringVarP(&dbChartOutputFile, "output", "o", "", "Arquivo de saída (para HTML)")

	rootCmd.AddCommand(dbChartCmd)
}

var dbChartCmd = &cobra.Command{
	Use:   "db-chart",
	Short: "Gerar gráficos de análises de banco de dados",
	Long: `Gera gráficos visuais a partir de resultados de análises de banco de dados.

A IA sugere o melhor tipo de gráfico e extrai dados automaticamente dos resultados.

Exemplos:
  snip db-chart --analysis-id 1
  snip db-chart --analysis-id 1 --type bar
  snip db-chart --analysis-id 1 --type html --output chart.html`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbChartAnalysisID == 0 {
			fmt.Println("Erro: --analysis-id é obrigatório")
			return
		}

		if err := executeWithDBAnalysisHandler(func(h handler.DBAnalysisHandler) error {
			// Obter análise
			analysis, err := h.GetAnalysisByID(dbChartAnalysisID)
			if err != nil {
				return fmt.Errorf("erro ao obter análise: %w", err)
			}

			if analysis.Result == "" {
				return fmt.Errorf("análise ainda não foi executada")
			}

			// Criar gerador de gráficos
			chartGenerator, err := dbcharts.NewChartGenerator()
			if err != nil {
				return fmt.Errorf("erro ao criar gerador de gráficos: %w", err)
			}

			// Determinar tipo de gráfico
			chartType := dbcharts.ChartTypeASCII
			if dbChartChartType != "" {
				chartType = dbcharts.ChartType(dbChartChartType)
			} else {
				// Sugerir tipo com IA
				suggestedType, reason, err := chartGenerator.SuggestChartType(analysis.Result)
				if err == nil {
					chartType = suggestedType
					fmt.Printf("Tipo de gráfico sugerido: %s\nRazão: %s\n\n", suggestedType, reason)
				}
			}

			// Gerar gráfico
			chart, err := chartGenerator.GenerateChartFromAnalysis(analysis.Result, chartType)
			if err != nil {
				return fmt.Errorf("erro ao gerar gráfico: %w", err)
			}

			// Salvar ou exibir
			if dbChartOutputFile != "" && chartType == dbcharts.ChartTypeHTML {
				err := os.WriteFile(dbChartOutputFile, []byte(chart), 0644)
				if err != nil {
					return fmt.Errorf("erro ao salvar arquivo: %w", err)
				}
				fmt.Printf("Gráfico salvo em: %s\n", dbChartOutputFile)
			} else {
				fmt.Println(chart)
			}

			return nil
		}); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	},
}

