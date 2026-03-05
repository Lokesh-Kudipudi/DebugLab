package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/runner"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <problem-name>",
	Short: "Start a debugging session for a problem",
	Long:  "Copies the problem to your workspace, starts a timer, and displays the bug ticket.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		problemName := args[0]

		// 1. Fetch remote catalog to find path
		catalog, err := runner.FetchCatalog()
		if err != nil {
			return fmt.Errorf("failed to fetch problem catalog: %w", err)
		}

		var targetPath string
		for _, entry := range catalog {
			if entry.Name == problemName {
				targetPath = entry.Path
				break
			}
		}

		if targetPath == "" {
			return fmt.Errorf("problem %q not found in the remote catalog", problemName)
		}

		// Check if there's already an active session
		existing, err := db.GetActiveSessionByProblem(problemName)
		if err == nil && existing != nil {
			color.Yellow("⚠ You already have an active session for %q.", problemName)
			color.Yellow("  Started at: %s", existing.StartedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
			color.White("  Run " + color.CyanString("dblab verify %s", problemName) + " when you're done.")
			return nil
		}

		// Set up workspace in the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current directory: %w", err)
		}

		workspaceDir := filepath.Join(cwd, "dblab-workspace", problemName)

		// Remove existing workspace if any
		if err := os.RemoveAll(workspaceDir); err != nil {
			return fmt.Errorf("failed to clean workspace: %w", err)
		}

		// Download problem from GitHub into workspace
		color.Cyan("\nDownloading %q from GitHub...\n", problemName)
		if err := runner.DownloadProblem(targetPath, workspaceDir); err != nil {
			os.RemoveAll(workspaceDir) // Clean up partial download
			return fmt.Errorf("failed to download problem to workspace: %w", err)
		}

		// Load problem definition from newly created workspace
		problem, err := runner.LoadProblem(workspaceDir)
		if err != nil {
			return fmt.Errorf("could not load problem definition after download: %w", err)
		}

		// Initialize git repo for tracking changes
		if err := initGitRepo(workspaceDir); err != nil {
			// Non-fatal: git tracking just won't work
			color.Yellow("⚠ Could not init git repo for change tracking: %v", err)
		}

		// Create session in database
		session, err := db.CreateSession(problemName)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		// Print the ticket
		fmt.Println()
		color.New(color.FgCyan, color.Bold).Printf("🔬 %s\n", problem.Title)
		fmt.Println()

		// Difficulty badge
		diff := difficultyColor(problem.Difficulty)
		diff.Printf("   Difficulty: %s", problem.Difficulty)
		fmt.Printf("  |  Language: %s  |  Par time: %d min\n", problem.Language, problem.ParTimeMinutes)
		fmt.Println()

		color.New(color.FgWhite, color.Bold).Println("📋 Ticket:")
		fmt.Println()
		color.White("   %s", problem.Ticket)
		fmt.Println()

		if len(problem.Hints) > 0 {
			color.New(color.FgYellow).Printf("💡 %d hint(s) available (using hints reduces your score)\n", len(problem.Hints))
		}

		fmt.Println()
		fmt.Println("  " + color.New(color.FgGreen, color.Bold).Sprint("─────────────────────────────────────────────"))
		fmt.Println()
		color.New(color.FgWhite).Printf("  📂 Workspace: %s\n", workspaceDir)
		color.New(color.FgWhite).Printf("  ⏱  Timer started at: %s\n", session.StartedAt.Format("15:04:05"))
		fmt.Println()
		color.New(color.FgGreen, color.Bold).Print("  Your editor is ready. ")
		color.New(color.FgWhite).Printf("Run ")
		color.New(color.FgCyan, color.Bold).Printf("dblab verify %s", problemName)
		color.New(color.FgWhite).Println(" when done.")
		fmt.Println()

		return nil
	},
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// initGitRepo initializes a git repo and commits all files.
func initGitRepo(dir string) error {
	commands := [][]string{
		{"git", "init"},
		{"git", "add", "."},
		{"git", "commit", "-m", "Initial problem state"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=dblab",
			"GIT_AUTHOR_EMAIL=dblab@local",
			"GIT_COMMITTER_NAME=dblab",
			"GIT_COMMITTER_EMAIL=dblab@local",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s failed: %w\n%s", args[0], err, string(out))
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
