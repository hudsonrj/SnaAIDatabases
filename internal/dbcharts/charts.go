package dbcharts

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/snip/internal/ai"
)

// ChartType representa o tipo de gráfico
type ChartType string

const (
	ChartTypeLine      ChartType = "line"
	ChartTypeBar       ChartType = "bar"
	ChartTypePie       ChartType = "pie"
	ChartTypeArea      ChartType = "area"
	ChartTypeTable     ChartType = "table"
	ChartTypeASCII     ChartType = "ascii"
	ChartTypeHTML      ChartType = "html"
)

// ChartData representa dados para um gráfico
type ChartData struct {
	Labels   []string    `json:"labels"`
	Series   []Series    `json:"series"`
	Title    string      `json:"title"`
	XAxis    string      `json:"x_axis"`
	YAxis    string      `json:"y_axis"`
	ChartType ChartType  `json:"chart_type"`
}

// Series representa uma série de dados
type Series struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
	Color  string    `json:"color,omitempty"`
}

// ChartGenerator gera gráficos de análises
type ChartGenerator struct {
	aiClient *ai.GroqClient
}

// NewChartGenerator cria um novo gerador de gráficos
func NewChartGenerator() (*ChartGenerator, error) {
	aiClient, err := ai.NewGroqClient()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente IA: %w", err)
	}
	return &ChartGenerator{aiClient: aiClient}, nil
}

// GenerateChartFromData gera um gráfico a partir de dados estruturados
func (c *ChartGenerator) GenerateChartFromData(data ChartData) (string, error) {
	switch data.ChartType {
	case ChartTypeASCII:
		return c.generateASCIIChart(data)
	case ChartTypeHTML:
		return c.generateHTMLChart(data)
	case ChartTypeTable:
		return c.generateTable(data)
	default:
		return c.generateASCIIChart(data)
	}
}

// GenerateChartFromAnalysis gera gráfico a partir de resultado de análise usando IA
func (c *ChartGenerator) GenerateChartFromAnalysis(analysisResult string, chartType ChartType) (string, error) {
	// Usar IA para extrair dados estruturados do resultado
	extractedData, err := c.extractDataWithAI(analysisResult, chartType)
	if err != nil {
		return "", fmt.Errorf("erro ao extrair dados: %w", err)
	}

	// Gerar gráfico
	return c.GenerateChartFromData(extractedData)
}

// extractDataWithAI usa IA para extrair dados estruturados do resultado da análise
func (c *ChartGenerator) extractDataWithAI(analysisResult string, chartType ChartType) (ChartData, error) {
	prompt := fmt.Sprintf(`Analise o seguinte resultado de análise de banco de dados e extraia dados numéricos que possam ser visualizados em um gráfico do tipo %s.

Resultado da análise:
%s

Retorne APENAS um JSON válido com a seguinte estrutura:
{
  "labels": ["label1", "label2", ...],
  "series": [
    {
      "name": "Nome da série",
      "values": [valor1, valor2, ...]
    }
  ],
  "title": "Título do gráfico",
  "x_axis": "Rótulo do eixo X",
  "y_axis": "Rótulo do eixo Y",
  "chart_type": "%s"
}

Se não houver dados numéricos suficientes, retorne um JSON com arrays vazios.`, chartType, analysisResult, chartType)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em análise de dados que extrai informações estruturadas de textos para visualização em gráficos.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 2000, 0.3)
	if err != nil {
		return ChartData{}, err
	}

	// Limpar resposta (remover markdown code blocks se houver)
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var data ChartData
	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		return ChartData{}, fmt.Errorf("erro ao parsear JSON: %w", err)
	}

	data.ChartType = chartType
	return data, nil
}

// generateASCIIChart gera gráfico ASCII
func (c *ChartGenerator) generateASCIIChart(data ChartData) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("\n%s\n", data.Title))
	result.WriteString(strings.Repeat("=", len(data.Title)))
	result.WriteString("\n\n")

	if len(data.Series) == 0 || len(data.Labels) == 0 {
		return "Dados insuficientes para gerar gráfico", nil
	}

	// Encontrar valor máximo para escala
	maxValue := 0.0
	for _, series := range data.Series {
		for _, val := range series.Values {
			if val > maxValue {
				maxValue = val
			}
		}
	}

	if maxValue == 0 {
		return "Todos os valores são zero", nil
	}

	// Altura do gráfico
	height := 20
	width := len(data.Labels)

	// Para gráfico de barras
	if data.ChartType == ChartTypeBar || data.ChartType == "" {
		result.WriteString(fmt.Sprintf("%s\n", data.YAxis))
		result.WriteString("│\n")

		for i := height; i >= 0; i-- {
			value := maxValue * float64(i) / float64(height)
			result.WriteString(fmt.Sprintf("%6.1f│", value))

			for j := 0; j < width && j < len(data.Labels); j++ {
				hasBar := false
				for _, series := range data.Series {
					if j < len(series.Values) && series.Values[j] >= value {
						hasBar = true
						break
					}
				}
				if hasBar {
					result.WriteString("█")
				} else {
					result.WriteString(" ")
				}
			}
			result.WriteString("\n")
		}

		result.WriteString("      └")
		result.WriteString(strings.Repeat("─", width))
		result.WriteString("\n      ")

		// Labels
		for i, label := range data.Labels {
			if i < width {
				if len(label) > 8 {
					label = label[:8]
				}
				result.WriteString(fmt.Sprintf("%-8s", label))
			}
		}
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("      %s\n\n", data.XAxis))

		// Legenda
		for i, series := range data.Series {
			if i < len(data.Series) {
				result.WriteString(fmt.Sprintf("  %s: █\n", series.Name))
			}
		}
	}

	// Para gráfico de linha
	if data.ChartType == ChartTypeLine {
		result.WriteString(fmt.Sprintf("%s\n", data.YAxis))
		result.WriteString("│\n")

		for i := height; i >= 0; i-- {
			value := maxValue * float64(i) / float64(height)
			result.WriteString(fmt.Sprintf("%6.1f│", value))

			for j := 0; j < width && j < len(data.Labels); j++ {
				hasPoint := false
				for _, series := range data.Series {
					if j < len(series.Values) {
						diff := math.Abs(series.Values[j] - value)
						if diff < maxValue/float64(height) {
							hasPoint = true
							break
						}
					}
				}
				if hasPoint {
					result.WriteString("●")
				} else {
					result.WriteString(" ")
				}
			}
			result.WriteString("\n")
		}

		result.WriteString("      └")
		result.WriteString(strings.Repeat("─", width))
		result.WriteString("\n      ")

		for i, label := range data.Labels {
			if i < width {
				if len(label) > 8 {
					label = label[:8]
				}
				result.WriteString(fmt.Sprintf("%-8s", label))
			}
		}
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("      %s\n\n", data.XAxis))
	}

	return result.String(), nil
}

// generateHTMLChart gera gráfico HTML usando Chart.js
func (c *ChartGenerator) generateHTMLChart(data ChartData) (string, error) {
	var result strings.Builder

	result.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>` + data.Title + `</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .chart-container { width: 800px; height: 400px; margin: 20px 0; }
    </style>
</head>
<body>
    <h1>` + data.Title + `</h1>
    <div class="chart-container">
        <canvas id="chart"></canvas>
    </div>
    <script>
        const ctx = document.getElementById('chart').getContext('2d');
        const chart = new Chart(ctx, {
            type: '` + string(data.ChartType) + `',
            data: {
                labels: `)

	labelsJSON, _ := json.Marshal(data.Labels)
	result.WriteString(string(labelsJSON))
	result.WriteString(`,
                datasets: [`)

	for i, series := range data.Series {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(`{
                    label: '` + series.Name + `',
                    data: `)
		valuesJSON, _ := json.Marshal(series.Values)
		result.WriteString(string(valuesJSON))
		if series.Color != "" {
			result.WriteString(`,
                    backgroundColor: '` + series.Color + `',
                    borderColor: '` + series.Color + `'`)
		}
		result.WriteString(`
                }`)
	}

	result.WriteString(`
            ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: '` + data.YAxis + `'
                        }
                    },
                    x: {
                        title: {
                            display: true,
                            text: '` + data.XAxis + `'
                        }
                    }
                }
            }
        });
    </script>
</body>
</html>`)

	return result.String(), nil
}

// generateTable gera tabela formatada
func (c *ChartGenerator) generateTable(data ChartData) (string, error) {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("\n%s\n", data.Title))
	result.WriteString(strings.Repeat("=", len(data.Title)))
	result.WriteString("\n\n")

	if len(data.Labels) == 0 {
		return "Sem dados para exibir", nil
	}

	// Cabeçalho
	result.WriteString("| ")
	result.WriteString(data.XAxis)
	for _, series := range data.Series {
		result.WriteString(" | ")
		result.WriteString(series.Name)
	}
	result.WriteString(" |\n")

	// Separador
	result.WriteString("|")
	for i := 0; i <= len(data.Series); i++ {
		result.WriteString("---|")
	}
	result.WriteString("\n")

	// Dados
	for i, label := range data.Labels {
		result.WriteString("| ")
		result.WriteString(label)
		for _, series := range data.Series {
			result.WriteString(" | ")
			if i < len(series.Values) {
				result.WriteString(fmt.Sprintf("%.2f", series.Values[i]))
			} else {
				result.WriteString("N/A")
			}
		}
		result.WriteString(" |\n")
	}

	return result.String(), nil
}

// SuggestChartType sugere o tipo de gráfico mais apropriado usando IA
func (c *ChartGenerator) SuggestChartType(analysisResult string) (ChartType, string, error) {
	prompt := fmt.Sprintf(`Analise o seguinte resultado de análise de banco de dados e sugira o tipo de gráfico mais apropriado para visualizar os dados.

Resultado:
%s

Retorne APENAS um JSON com:
{
  "chart_type": "line" ou "bar" ou "pie" ou "area" ou "table",
  "reason": "explicação breve do porquê"
}`, analysisResult)

	messages := []ai.Message{
		{
			Role:    "system",
			Content: "Você é um especialista em visualização de dados que sugere o melhor tipo de gráfico para diferentes tipos de dados.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.aiClient.Chat(messages, 500, 0.3)
	if err != nil {
		return ChartTypeTable, "", err
	}

	// Limpar resposta
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var suggestion struct {
		ChartType string `json:"chart_type"`
		Reason    string `json:"reason"`
	}

	err = json.Unmarshal([]byte(response), &suggestion)
	if err != nil {
		return ChartTypeTable, "", err
	}

	return ChartType(suggestion.ChartType), suggestion.Reason, nil
}
