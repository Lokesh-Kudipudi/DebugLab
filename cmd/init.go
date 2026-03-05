package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/config"
	"github.com/lokeshkudipudi/dblab/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a DebugLab workspace",
	Long:  "Creates the workspace directory and initializes the local SQLite database.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}

		// Initialize the database and workspace path
		dbPath, err := store.DBPath(absDir)
		if err != nil {
			return fmt.Errorf("failed to determine workspace path: %w", err)
		}

		s, err := store.NewStore(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize local database: %w", err)
		}
		s.Close()

		// Save to global config
		cfg := config.Config{WorkspaceDir: absDir}
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		color.Green("\n✅ DebugLab workspace initialized successfully!\n")
		fmt.Printf("   Workspace: %s\n", absDir)
		fmt.Printf("   Database:  %s\n\n", dbPath)
		fmt.Println("You can now list problems with:  dblab list")
		fmt.Println("And start a problem with:        dblab start <problem-name>")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
