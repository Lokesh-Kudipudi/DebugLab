package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a DebugLab workspace in the current directory",
	Long:  "Creates the ./dblab-workspace directory and initializes the local SQLite database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize the database and workspace path
		dbPath, err := store.DefaultDBPath()
		if err != nil {
			return fmt.Errorf("failed to determine workspace path: %w", err)
		}

		s, err := store.NewStore(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize local database: %w", err)
		}
		s.Close()

		color.Green("\n✅ DebugLab workspace initialized successfully!\n")
		fmt.Printf("   Workspace: ./dblab-workspace/\n")
		fmt.Printf("   Database:  %s\n\n", dbPath)
		fmt.Println("You can now list problems with:  dblab list")
		fmt.Println("And start a problem with:        dblab start <problem-name>")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
