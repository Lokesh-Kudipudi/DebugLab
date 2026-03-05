package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/config"
	"github.com/lokeshkudipudi/dblab/internal/store"
	"github.com/spf13/cobra"
)

var (
	db *store.Store
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "dblab",
	Short: "DebugLab — Practice debugging real-world bugs",
	Long: color.New(color.FgCyan, color.Bold).Sprint("🔬 DebugLab") + ` — A CLI tool for practicing debugging skills.

Pick a problem, find and fix the bugs, then verify your solution.
You'll be scored on time, accuracy, and efficiency.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Try loading configure, but don't fail immediately because `dblab init` doesn't need it
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		// If running `dblab init` or `dblab help`, don't error out entirely
		// We'll let the individual subcommands error out if they require a workspace
		return
	}

	// Initialize database path from config
	dbPath, err := store.DBPath(cfg.WorkspaceDir)
	if err != nil {
		color.Red("Error resolving DB path: %v", err)
		os.Exit(1)
	}

	db, err = store.NewStore(dbPath)
	if err != nil {
		color.Red("Error opening DB: %v", err)
		os.Exit(1)
	}
}
