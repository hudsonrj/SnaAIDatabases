package cmd

import (
	"fmt"
	"time"

	"github.com/snip/internal/checklist"
	"github.com/spf13/cobra"
)

var (
	bulkChecklistCSV    string
	bulkChecklistOutput string
	bulkChecklistType   string
	bulkChecklistTemplateOutput string
)

func init() {
	bulkChecklistCmd.Flags().StringVarP(&bulkChecklistCSV, "csv", "c", "", "Caminho do arquivo CSV com os itens do checklist")
	bulkChecklistCmd.Flags().StringVarP(&bulkChecklistOutput, "output", "o", "", "Nome do arquivo markdown de saída (opcional)")

	bulkChecklistTemplateCmd.Flags().StringVarP(&bulkChecklistType, "type", "t", "generic", "Tipo de checklist (generic, daily, weekly, deep, backup, security, performance, maintenance)")
	bulkChecklistTemplateCmd.Flags().StringVarP(&bulkChecklistTemplateOutput, "output", "o", "", "Caminho do arquivo CSV de saída (opcional)")

	// Adicionar ao checklistCmd (definido em checklist.go)
	checklistCmd.AddCommand(bulkChecklistCmd)
	bulkChecklistCmd.AddCommand(bulkChecklistTemplateCmd)
}

var bulkChecklistCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Processar checklist em massa a partir de CSV",
	Long: `Processa um checklist em massa a partir de um arquivo CSV.

O CSV deve ter os seguintes campos (obrigatórios marcados com *):
  - title*: Título do item
  - description*: Descrição do item
  - category*: Categoria do item
  - priority: Prioridade (high/medium/low) - padrão: medium
  - status: Status (pending/completed/failed) - padrão: pending
  - notes: Notas adicionais

Exemplo de CSV:
  title,description,category,priority,status,notes
  "Verificar backups","Verificar se backups estão sendo executados","Backup",high,pending,"Verificar logs"
  "Testar restore","Testar procedimento de restore","Backup",medium,completed,"Teste realizado com sucesso"

Exemplos:
  snip checklist bulk --csv checklist.csv
  snip checklist bulk --csv items.csv --output resultado.md`,
	Run: func(cmd *cobra.Command, args []string) {
		if bulkChecklistCSV == "" {
			fmt.Println("❌ Caminho do arquivo CSV é obrigatório (use --csv)")
			return
		}

		fmt.Printf("📋 Processando checklist em massa de: %s\n", bulkChecklistCSV)
		fmt.Println("Aguarde...\n")

		result, err := checklist.ProcessBulkChecklistFromCSV(bulkChecklistCSV)
		if err != nil {
			fmt.Printf("❌ Erro ao processar CSV: %v\n", err)
			return
		}

		fmt.Printf("✅ Processamento concluído!\n")
		fmt.Printf("  Total de itens: %d\n", result.TotalItems)
		fmt.Printf("  ✅ Concluídos: %d\n", result.Completed)
		fmt.Printf("  ⏳ Pendentes: %d\n", result.Pending)
		fmt.Printf("  ❌ Falhas: %d\n", result.Failed)
		fmt.Printf("  Tempo de execução: %s\n\n", result.ExecutionTime.Round(time.Second))

		// Exportar para markdown
		filePath, err := checklist.ExportBulkChecklistToMarkdown(result, bulkChecklistOutput)
		if err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao exportar: %v\n", err)
			return
		}

		fmt.Printf("📄 Relatório exportado para: %s\n", filePath)
	},
}

var bulkChecklistTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Gerar template CSV para checklist",
	Long: `Gera um template CSV para um tipo específico de checklist.

Tipos disponíveis:
  - generic: Checklist genérico com campos básicos
  - daily: Checklist diário para tarefas rotineiras
  - weekly: Checklist semanal para revisões periódicas
  - deep: Checklist profundo para análises detalhadas
  - backup: Checklist específico para backups e restores
  - security: Checklist de segurança e compliance
  - performance: Checklist de performance e monitoramento
  - maintenance: Checklist de manutenção preventiva

Exemplos:
  snip checklist bulk template --type daily --output daily_checklist.csv
  snip checklist bulk template --type backup
  snip checklist bulk template --type security -o security_audit.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		checklistType := checklist.ChecklistType(bulkChecklistType)
		
		// Validar tipo
		availableTypes := checklist.GetAvailableChecklistTypes()
		valid := false
		for _, t := range availableTypes {
			if t == checklistType {
				valid = true
				break
			}
		}
		
		if !valid {
			fmt.Printf("❌ Tipo de checklist inválido: %s\n", bulkChecklistType)
			fmt.Println("\nTipos disponíveis:")
			for _, t := range availableTypes {
				fmt.Printf("  - %s: %s\n", t, checklist.GetChecklistTypeDescription(t))
			}
			return
		}

		outputPath := bulkChecklistTemplateOutput
		if outputPath == "" {
			outputPath = fmt.Sprintf("checklist_template_%s.csv", bulkChecklistType)
		}

		fmt.Printf("📝 Gerando template CSV para checklist tipo: %s\n", bulkChecklistType)
		fmt.Printf("   Descrição: %s\n", checklist.GetChecklistTypeDescription(checklistType))
		fmt.Println("Aguarde...\n")

		err := checklist.GenerateCSVTemplate(checklistType, outputPath)
		if err != nil {
			fmt.Printf("❌ Erro ao gerar template: %v\n", err)
			return
		}

		fmt.Printf("✅ Template gerado com sucesso!\n")
		fmt.Printf("   Arquivo: %s\n", outputPath)
		fmt.Printf("\n💡 Você pode editar este arquivo e usar com:\n")
		fmt.Printf("   snip checklist bulk --csv %s\n", outputPath)
	},
}

