package cmd

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/store"
	"github.com/spf13/cobra"
)

var (
	db          *store.Store
	problemsDir string
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
	// Initialize database
	dbPath, err := store.DefaultDBPath()
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	db, err = store.NewStore(dbPath)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	// Find problems directory relative to executable
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		problemsDir = "problems"
		return
	}

	// Check relative to executable
	execDir := filepath.Dir(execPath)
	candidate := filepath.Join(execDir, "problems")
	if _, err := os.Stat(candidate); err == nil {
		problemsDir = candidate
		return
	}

	// Fallback to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		problemsDir = "problems"
		return
	}
	problemsDir = filepath.Join(cwd, "problems")
}
