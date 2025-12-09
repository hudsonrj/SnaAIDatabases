package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var aiCodeLang string
var aiCodeContext string

func init() {
	aiCodeCmd.Flags().StringVarP(&aiCodeLang, "lang", "l", "go", "Programming language")
	aiCodeCmd.Flags().StringVarP(&aiCodeContext, "context", "c", "", "Additional context for code generation")
	rootCmd.AddCommand(aiCodeCmd)
}

var aiCodeCmd = &cobra.Command{
	Use:   "ai-code [description]",
	Short: "Generate code with AI",
	Long: `Generate code in a specific programming language using AI.

The AI will generate clean, well-documented code following best practices.

Examples:
  snip ai-code "function to reverse a string"
  snip ai-code "REST API endpoint" --lang "python" --context "Use FastAPI"
  snip ai-code "binary search algorithm" --lang "javascript"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			description := strings.Join(args, " ")
			context := ""
			if aiCodeContext != "" {
				context = aiCodeContext
			}
			lang := "go"
			if aiCodeLang != "" {
				lang = aiCodeLang
			}
			return h.GenerateCodeWithAI(lang, description, context)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

