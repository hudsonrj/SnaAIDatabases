package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aiSearchCmd)
}

var aiSearchCmd = &cobra.Command{
	Use:   "ai-search [query]",
	Short: "Improve search query with AI and search notes",
	Long: `Use AI to improve your search query and then search through your notes.

The AI will enhance your search query based on the context of your existing notes
to help you find more relevant results.

Examples:
  snip ai-search "meeting notes"
  snip ai-search "python tutorial"
  snip ai-search "project ideas"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			query := strings.Join(args, " ")
			return h.ImproveSearchWithAI(query)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

