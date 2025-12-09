package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aiAskCmd)
}

var aiAskCmd = &cobra.Command{
	Use:   "ai-ask [question]",
	Short: "Ask a question to AI based on your notes",
	Long: `Ask a question to AI that will use your notes as context to provide answers.

The AI will search through your notes and provide answers based on the information
found. If the answer is not in your notes, it will use general knowledge.

Examples:
  snip ai-ask "What did I write about Python?"
  snip ai-ask "Summarize my meeting notes"
  snip ai-ask "What are the main topics in my notes?"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			question := strings.Join(args, " ")
			return h.AskAI(question)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

